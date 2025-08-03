package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/invopop/jsonschema"
	jsonStream "github.com/joan-cb/schema-stream/jsonStream"
	schemabuilder "github.com/joan-cb/schema-stream/schemaBuilder"
)

var schema *jsonschema.Schema
var structFields []reflect.StructField

func readJSONFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Validate it's proper JSON
	if !json.Valid(data) {
		return nil, fmt.Errorf("file %s contains invalid JSON", filename)
	}

	return data, nil
}

func processJSONToSchema(inputFile string) error {
	// 1. Read input
	jsonData, err := readJSONFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	// 2. Parse JSON to struct fields
	structFields, err = jsonStream.ReturnStructDefinition(jsonData)
	if err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	// 4. Generate schema
	schema = schemabuilder.ReturnSchemaFromStructFields(structFields)

	return nil
}

func findJSONFilesInDirectory(directory string) ([]string, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", directory, err)
	}
	jsonFiles := []string{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			jsonFiles = append(jsonFiles, file.Name())
		}
	}
	return jsonFiles, nil
}
