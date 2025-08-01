package schemaStream

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

func ParseAndGetTypes(jsonData []byte) (map[string]string, error) {
	reader := strings.NewReader(string(jsonData))
	decoder := json.NewDecoder(reader)

	types := make(map[string]string)
	err := analyzeJSON(decoder, "", types)
	if err != nil {
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}

	return types, nil
}

func analyzeJSON(decoder *json.Decoder, path string, types map[string]string) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	// Handle delimiters
	if delim, ok := token.(json.Delim); ok {
		switch delim {
		case '{':
			return analyzeObject(decoder, path, types)
		case '[':
			return analyzeArray(decoder, path, types)
		}
	}

	// Handle primitive types using reflection
	types[path] = getTypeName(token)
	return nil
}

func getTypeName(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "bool"
	case reflect.Float64, reflect.Int64:
		return "number"
	case reflect.Map:
		return "object"
	case reflect.Slice:
		return "array"
	default:
		return "unknown"
	}
}

func analyzeObject(decoder *json.Decoder, path string, types map[string]string) error {

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		if delim, ok := token.(json.Delim); ok && delim == '}' {
			break
		}

		key := token.(string)
		currentPath := key
		if path != "" {
			currentPath = path + "." + key
		}

		if err := analyzeJSON(decoder, currentPath, types); err != nil {
			return err
		}
	}
	return nil
}

func analyzeArray(decoder *json.Decoder, path string, types map[string]string) error {
	// Analyze first element to determine type with proper path
	types[path] = "array"
	firstElementPath := path + "[0]"
	if err := analyzeJSON(decoder, firstElementPath, types); err != nil {
		return err
	}

	// Skip to end of array
	skipArray(decoder)

	// Get the type from the first element and format as array
	if firstType, exists := types[firstElementPath]; exists {
		types[path] = "[" + firstType + "]"
		delete(types, firstElementPath)
	} else {
		log.Println("Unknown type for array:", path)
		types[path] = "[unknown]"
	}

	return nil
}

func skipArray(decoder *json.Decoder) {
	for {
		token, _ := decoder.Token()
		if delim, ok := token.(json.Delim); ok && delim == ']' {
			break
		}
	}
}

func skipObject(decoder *json.Decoder) {
	for {
		token, _ := decoder.Token()
		if delim, ok := token.(json.Delim); ok && delim == '}' {
			break
		}
	}
}
