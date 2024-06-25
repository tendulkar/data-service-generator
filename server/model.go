package server

// Define structures to match your JSON input for attributes and models
type Attribute struct {
	ID          int    `json:"id"`
	Namespace   string `json:"namespace"`
	Family      string `json:"family"`
	Name        string `json:"name"`
	Label       string `json:"label"`
	TypeID      int    `json:"type_id"`
	Validations []int  `json:"validations,omitempty"`
}

type Model struct {
	ID                int    `json:"id"`
	Namespace         string `json:"namespace"`
	Family            string `json:"family"`
	Name              string `json:"name"`
	Attributes        []int  `json:"attributes"`
	UniqueConstraints []struct {
		ConstraintName string `json:"constraint_name"`
		Attributes     []int  `json:"attributes"`
	} `json:"unique_constraints,omitempty"`
	Relationships []struct {
		Type          string `json:"type"`
		TargetModelID int    `json:"target_model_id"`
	} `json:"relationships,omitempty"`
}

type RequestData struct {
	Attributes []Attribute `json:"attributes"`
	Models     []Model     `json:"models"`
}
