package defs

type ParameterRef struct {
	Name  string `yaml:"name"`
	Index int32  `yaml:"index"`
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
		ConstraintName string `yaml:"constraint_name"`
		Attributes     []int  `yaml:"attributes"`
	} `yaml:"unique_constraints"`
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
	DriverName string `yaml:"driver_name"`
	DBConfigId string `yaml:"db_config_id"`
}

type DataConfig struct {
	Models           []ModelConfig     `yaml:"models"`
	ConnectionConfig *ConnectionConfig `yaml:"connection_config,omitempty"`
}
