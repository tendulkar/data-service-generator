package defs

import (
	"fmt"

	"stellarsky.ai/platform/codegen/data-service-generator/config"
	"stellarsky.ai/platform/codegen/data-service-generator/db/models"
)

type ParameterRef struct {
	Name     string        `yaml:"name"`
	Index    int32         `yaml:"index"`
	FuncName string        `yaml:"func_name"`
	FuncArgs []interface{} `yaml:"func_args"`
}

type Filter struct {
	Attribute      string   `yaml:"attribute,omitempty"`
	Transformation string   `yaml:"transformation,omitempty"`
	Operator       string   `yaml:"operator"`
	ParamName      string   `yaml:"param_name,omitempty"`
	Conditions     []Filter `yaml:"conditions,omitempty"`
}

type Model struct {
	ID                int     `yaml:"id"`
	Namespace         string  `yaml:"namespace"`
	Family            string  `yaml:"family"`
	Name              string  `yaml:"name"`
	Attributes        []int64 `yaml:"attributes"`
	UniqueConstraints []struct {
		ConstraintName string  `yaml:"constraint_name"`
		Attributes     []int64 `yaml:"attributes"`
	} `yaml:"unique_constraints"`
	Indexes []struct {
		IndexName  string  `yaml:"index_name"`
		Attributes []int64 `yaml:"attributes"`
	} `yaml:"indexes"`
}

type Access struct {
	Find         []AccessConfig `yaml:"find"`
	Update       []AccessConfig `yaml:"update"`
	Add          []AccessConfig `yaml:"add"`
	AddOrReplace []AccessConfig `yaml:"add_or_replace"`
	Delete       []AccessConfig `yaml:"delete"`
}

type ModelConfig struct {
	Model  `yaml:"model"`
	Access `yaml:"access"`
}

func (m *ModelConfig) GetAllAccessConfig() []AccessConfig {
	accessConfig := []AccessConfig{}
	accessConfig = append(accessConfig, m.Access.Find...)
	accessConfig = append(accessConfig, m.Access.Update...)
	accessConfig = append(accessConfig, m.Access.Add...)
	accessConfig = append(accessConfig, m.Access.AddOrReplace...)
	accessConfig = append(accessConfig, m.Access.Delete...)
	return accessConfig
}

func (m *ModelConfig) GetAllFilters() []Filter {
	filters := []Filter{}
	for _, accessConfig := range m.GetAllAccessConfig() {
		filters = append(filters, accessConfig.Filter...)
	}
	return filters
}

func (m *ModelConfig) GetAttributes() ([]*models.AttributeRow, error) {
	attributes := []*models.AttributeRow{}
	for _, attributeID := range m.Attributes {
		attribute, ok := config.Attributes[attributeID]
		if !ok {
			return nil, fmt.Errorf("ModelConfig attribute %d not found", attributeID)
		}
		attributes = append(attributes, &attribute)
	}
	return attributes, nil
}

type Index struct {
	IndexName  string  `yaml:"index_name"`
	Attributes []int64 `yaml:"attributes"`
	IsUnique   bool    `yaml:"is_unique"`
}

func (m *Model) GetIndexes() []Index {
	indexes := []Index{}
	for _, index := range m.Indexes {
		indexes = append(indexes, Index{
			IndexName:  index.IndexName,
			Attributes: index.Attributes,
			IsUnique:   false,
		})
	}

	for _, constraint := range m.UniqueConstraints {
		indexes = append(indexes, Index{
			IndexName:  constraint.ConstraintName,
			Attributes: constraint.Attributes,
			IsUnique:   true,
		})
	}
	return indexes
}

type Parameter struct {
	Param string `yaml:"param"`
}

type Request struct {
	Parameters []Parameter `yaml:"parameters"`
}

type AccessConfig struct {
	Name             string   `yaml:"name"`
	Request          *Request `yaml:"request,omitempty"`
	Attributes       []string `yaml:"attributes,omitempty"`
	Filter           []Filter `yaml:"filter,omitempty"`
	Set              []Update `yaml:"set,omitempty"`
	Autoincrement    []string `yaml:"autoincrement,omitempty"`
	CaptureTimestamp []string `yaml:"capture_timestamp,omitempty"`
	Values           []Update `yaml:"values,omitempty"`
}

type Update struct {
	Attribute string `yaml:"attribute"`
	ParamName string `yaml:"param_name"`
}

type Attribute struct {
	Attribute string `yaml:"attribute"`
}

type ConnectionConfig struct {
	IdleTimeoutSecs int `yaml:"idle_timeout_secs"`
	MaxLifetimeMins int `yaml:"conn_max_lifetime_mins"`
}

type ConnectionPoolConfig struct {
	MaxIdleConns int `yaml:"max_idle_conns"`
	MaxOpenConns int `yaml:"max_open_conns"`
}

type DatabaseConfig struct {
	DriverName           string                `yaml:"driver_name"`
	DBConfigId           string                `yaml:"db_config_id"`
	UserName             string                `yaml:"user_name"`
	Password             string                `yaml:"password"`
	Host                 string                `yaml:"host"`
	Port                 int                   `yaml:"port"`
	DBName               string                `yaml:"db_name"`
	ConnectionConfig     *ConnectionConfig     `yaml:"conn_config,omitempty"`
	ConnectionPoolConfig *ConnectionPoolConfig `yaml:"conn_pool_config,omitempty"`
}

type DataConfig struct {
	FamilyName     string          `yaml:"family_name"`
	Models         []ModelConfig   `yaml:"models"`
	DatabaseConfig *DatabaseConfig `yaml:"connection_config,omitempty"`
}
