package generator

import (
	"stellarsky.ai/platform/codegen/data-service-generator/db/models"
)

const (
	typesJsonPath            = "db/data/types.json"
	validationsJsonPath      = "db/data/validations.json"
	postgresTypeMapsJsonPath = "db/postgres/data/type_maps.json"
)

var (
	types            map[string]models.TypeInfo
	validations      map[string]models.Validation
	postgresTypeMaps map[string]models.TypeMapping
)

func readJsonData() {
}

func ReadConfig() {

}
