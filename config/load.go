// config/load.go
package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"stellarsky.ai/platform/codegen/data-service-generator/base"
)

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
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

	return &config
}
