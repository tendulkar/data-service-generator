package second

type Condition struct {
	Attribute string
}

type EqualCondition struct {
	Attribute string
	ParamName interface{}
}

type NotEqualCondition struct {
	Attribute string
	ParamName interface{}
}

type LessThanCondition struct {
	Attribute string
	ParamName interface{}
}

type LessThanOrEqualCondition struct {
	Attribute string
	ParamName interface{}
}

type GreaterThanCondition struct {
	Attribute string
	ParamName interface{}
}

type GreaterThanOrEqualCondition struct {
	Attribute string
	ParamName interface{}
}

type AmongCondition struct {
	Attribute string        `yaml:"attribute"`
	ParamName []interface{} `yaml:"param"`
}

type RangeCondition struct {
	Attribute string
	LowParam  interface{}
	HighParam interface{}
}

type CommonFilter struct {
	EqualConditions              []EqualCondition              `yaml:"equal"`
	NotEqualConditions           []NotEqualCondition           `yaml:"not_equal"`
	LessThanConditions           []LessThanCondition           `yaml:"less_than"`
	LessThanOrEqualConditions    []LessThanOrEqualCondition    `yaml:"less_than_or_equal"`
	GreaterThanConditions        []GreaterThanCondition        `yaml:"greater_than"`
	GreaterThanOrEqualConditions []GreaterThanOrEqualCondition `yaml:"greater_than_or_equal"`
	AmongConditions              []AmongCondition              `yaml:"among"`
	RangeConditions              []RangeCondition              `yaml:"range"`
}

type AndFilter struct {
	CommonFilter
}

type OrFilter struct {
	CommonFilter
}

type NotFilter struct {
	CommonFilter
}

type Filter struct {
	*AndFilter                   `yaml:"and"`
	*OrFilter                    `yaml:"or"`
	*NotFilter                   `yaml:"not"`
	*EqualCondition              `yaml:"equal"`
	*NotEqualCondition           `yaml:"not_equal"`
	*LessThanCondition           `yaml:"less_than"`
	*LessThanOrEqualCondition    `yaml:"less_than_or_equal"`
	*GreaterThanCondition        `yaml:"greater_than"`
	*GreaterThanOrEqualCondition `yaml:"greater_than_or_equal"`
	*AmongCondition              `yaml:"among"`
	*RangeCondition              `yaml:"range"`
}

var exampleFilter string = `
filter:
  and:
	equal: 
		- attribute: name
		  param: name
	not_equal:
		- attribute: name
		  param: name
	less_than:
		- attribute: name
			param: name
	less_than_or_equal:
		- attribute: name
			param: name
	greater_than:
		- attribute: name
			param: name
	greater_than_or_equal:
		- attribute: name
			param: name
	among:
		- attribute: name
			param: [name, name]
	range:
		- attribute: name
			low_param: name
			high_param: name
		- attribute: name
			low_param: name
			high_param: name
`

func prepareStatBuilder(filter *Filter) string {

	return ""
}
