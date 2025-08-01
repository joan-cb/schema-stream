package schemaStream

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
)

var reader *strings.Reader
var decoder *json.Decoder
var types []reflect.StructField

func ParseAndGetTypes(jsonData []byte) ([]reflect.StructField, error) {
	reader = strings.NewReader(string(jsonData))
	decoder = json.NewDecoder(reader)
	types = []reflect.StructField{}

	err := analyzeJSON(&types)
	if err != nil {
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}

	log.Println("types", types)
	return types, nil
}

func analyzeJSON(types *[]reflect.StructField) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	// Handle delimiters
	if delim, ok := token.(json.Delim); ok {
		switch delim {
		case '{':
			return analyzeObject()
		case '[':
			return analyzeArray()
		}
	}
	return nil
}

func getStructField(token any) (reflect.StructField, error) {
	if token == nil {
		return reflect.StructField{}, nil
	}

	tokenName, tagName, err := returnFieldNameAndTag(token)
	if err != nil {
		return reflect.StructField{}, err
	}
	log.Println("Field name:", tokenName)

	valueToken, err := decoder.Token()
	if err != nil {
		return reflect.StructField{}, err
	}
	log.Println("Value token:", valueToken)

	switch reflect.TypeOf(valueToken).Kind() {
	case reflect.String:
		return reflect.StructField{
			Name: tokenName,
			Type: reflect.TypeOf(""),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, tagName)),
		}, nil
	case reflect.Bool:
		return reflect.StructField{
			Name: tokenName,
			Type: reflect.TypeOf(true),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, tagName)),
		}, nil
	case reflect.Float64, reflect.Int64:
		return reflect.StructField{
			Name: tokenName,
			Type: reflect.TypeOf(0.0),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, tagName)),
		}, nil
	//TO DO: CLEAN THIS UP BELOW
	case reflect.Map:
		log.Fatal("map")
		return reflect.StructField{}, nil
	case reflect.Slice:
		log.Fatal("slice")
		return reflect.StructField{}, nil
	default:
		log.Fatal("default")
		return reflect.StructField{}, fmt.Errorf("unknown type: %v", token)
	}
}

func analyzeObject() error {
	token, err := decoder.Token()
	if err != nil {
		log.Printf("analyzeObject returned error when fetching token: %s", err)
	}
	for {
		if fieldName, ok := token.(string); ok {
			// Get the value token
			valueToken, err := decoder.Token()
			if err != nil {
				return err
			}

			// Handle nested objects
			if delim, ok := valueToken.(json.Delim); ok && delim == '{' {
				// Skip nested object for now
				for {
					t, _ := decoder.Token()
					if d, ok := t.(json.Delim); ok && d == '}' {
						break
					}
				}
				continue
			}

			// Handle primitive values
			structInstance, err := getStructField(fieldName)
			if err != nil {
				return err
			}
			types = append(types, structInstance)
		}
	}
}

func analyzeArray() error {
	return nil
}

func returnFieldNameAndTag(token any) (string, string, error) {
	if str, ok := token.(string); ok {
		tokenName := strings.ToUpper(str[:1]) + str[1:]
		tagName := strings.ToLower(str[:1]) + str[1:]
		return tokenName, tagName, nil
	}
	return "", "", fmt.Errorf("unexpected error: %v", token)
}

// 	// Skip to end of array
// 	skipArray(decoder)

// 	// Get the type from the first element and format as array
// 	if firstType, exists := types[firstElementPath]; exists {
// 		types[path] = "[" + firstType + "]"
// 		delete(types, firstElementPath)
// 	} else {
// 		log.Println("Unknown type for array:", path)
// 		types[path] = "[unknown]"
// 	}

// 	return nil
// }

// func skipArray(decoder *json.Decoder) {
// 	for {
// 		token, _ := decoder.Token()
// 		if delim, ok := token.(json.Delim); ok && delim == ']' {
// 			break
// 		}
// 	}
// }

// func skipObject(decoder *json.Decoder) {
// 	for {
// 		token, _ := decoder.Token()
// 		if delim, ok := token.(json.Delim); ok && delim == '}' {
// 			break
// 		}
// 	}
// }
