package schemabuilder

import (
	"reflect"

	"github.com/invopop/jsonschema"
)

func ReturnSchemaFromStructFields(structFields []reflect.StructField) *jsonschema.Schema {
	concreteStructDefinition := reflect.StructOf(structFields)
	concreteStructInstance := reflect.New(concreteStructDefinition).Elem().Interface()
	// Create reflector with some options
	reflector := &jsonschema.Reflector{
		AllowAdditionalProperties:  false,
		RequiredFromJSONSchemaTags: true,
	}

	schema := reflector.Reflect(concreteStructInstance)
	return schema
}
