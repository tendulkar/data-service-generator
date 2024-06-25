package models

type UniqueID struct {
	ID        int64  `yaml:"id"`
	Namespace string `yaml:"namespace"`
	Family    string `yaml:"family"`
	Name      string `yaml:"name"`
}

type TypeInfo struct {
	UniqueID
	ElementType string `yaml:"element_type"`
	WidgetType  string `yaml:"widget_type"`
}

type Validation struct {
	UniqueID
	RuleName string   `yaml:"rule_name"`
	Params   []string `yaml:"params"`
}

type Attribute struct {
	UniqueID
	Type        TypeInfo     `yaml:"type"`
	Validations []Validation `yaml:"validations"`
}

type ConstraintType int

const (
	Unique ConstraintType = iota
	ForeignKey
	NotNull
)

type RelationType int

const (
	Child RelationType = iota
	Children
	BelongsTo
	Referes
)

type Constraint struct {
	ConstraintName string      `yaml:"constraint_name"`
	Attributes     []Attribute `yaml:"attributes"`
}

type Model struct {
	UniqueID
	Attributes        []Attribute  `yaml:"attributes"`
	UniqueConstraints []Constraint `yaml:"unique_constraints"`
}

type DataAccessRequest struct {
	InputFieldMap map[string]interface{} `yaml:"input_field_map"`
	OutputFileds  []string               `yaml:"output_fields"`
}

type DataAccessDef struct {
	UniqueID
	Request DataAccessRequest `yaml:"request"`
}

type TypeMapping struct {
	TypeID     int64  `yaml:"type_id"`
	MappedType string `yaml:"mapped_type"`
}
