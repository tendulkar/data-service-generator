// config/load.go
package config

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
	"stellarsky.ai/platform/codegen/data-service-generator/base"
	"stellarsky.ai/platform/codegen/data-service-generator/base/parser"
	"stellarsky.ai/platform/codegen/data-service-generator/db/models"
)

const ProjectDir = "/Users/yugandhar/code/stellarsky.ai/platform/codegen/data-service-generator"

var (
	Types            map[int64]models.TypeInfo
	Validations      map[int64]models.Validation
	PostgresTypeMaps map[int64]models.TypeMapping
	Attributes       map[int64]models.AttributeRow
)

func readBaseData(cfg *Config) {
	modelConfig := cfg.Model
	projectDir := ProjectDir
	baseDir := modelConfig.BaseDir
	typesPath := path.Join(projectDir, baseDir, modelConfig.TypesPath)
	validationsPath := path.Join(projectDir, baseDir, modelConfig.ValidationsPath)
	postgresTypeMapsPath := path.Join(projectDir, baseDir, modelConfig.PostgresTypeMapsPath)
	attributesPath := path.Join(projectDir, baseDir, modelConfig.AttributesPath)

	typesRead := parser.MustReadJsonToSlice[models.TypeInfo](typesPath, "types")
	validationsRead := parser.MustReadJsonToSlice[models.Validation](validationsPath, "validations")
	typeMapsRead := parser.MustReadJsonToSlice[models.TypeMapping](postgresTypeMapsPath, "postgres_type_mappings")
	attributesRead := parser.MustReadJsonToSlice[models.AttributeRow](attributesPath, "attributes")

	Types = make(map[int64]models.TypeInfo)
	for _, t := range typesRead {
		Types[t.ID] = t
	}
	Validations = make(map[int64]models.Validation)
	for _, v := range validationsRead {
		Validations[v.ID] = v
	}
	PostgresTypeMaps = make(map[int64]models.TypeMapping)
	for _, t := range typeMapsRead {
		PostgresTypeMaps[t.TypeID] = t
	}

	Attributes = make(map[int64]models.AttributeRow)
	for _, attribute := range attributesRead {
		Attributes[attribute.ID] = attribute
	}
	base.LOG.Info("Data loaded", "Types", Types, "Validations", Validations, "typeMaps", PostgresTypeMaps, "Attributes", Attributes)
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(path.Join(ProjectDir, "config"))
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		base.LOG.Error("Error reading config file", "error", err)
		os.Exit(1)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		base.LOG.Error("Error unmarshaling config", "error", err)
		os.Exit(1)
	}

	base.LOG.Info("Application Config loaded", "config", config)
	readBaseData(&config)
	return &config
}
