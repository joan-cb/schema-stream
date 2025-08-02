package schemaStream

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
)

var reader *strings.Reader
var decoder *json.Decoder

func ReturnStructDefinition(jsonData []byte) ([]reflect.StructField, error) {
	reader = strings.NewReader(string(jsonData))
	decoder = json.NewDecoder(reader)

	outerStruct, err := traverseJSON()
	if err != nil {
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}

	return outerStruct, nil
}

func traverseJSON() ([]reflect.StructField, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}

	switch token.(json.Delim) {
	case '[':
		err := traverseArray(token)
		if err != nil {
			return nil, fmt.Errorf("JSON parsing error: %v", err)
		}
		return []reflect.StructField{}, errors.New("The JSON document contains an array as a root level item, which is not supported.")
	case '{':
		structFields, err := traverseObject(token)
		if err != nil {
			return nil, fmt.Errorf("JSON parsing error: %v", err)
		}
		return structFields, nil
	default:
		return nil, fmt.Errorf("unexpected token: %v", token)
	}
}

func traverseObject(token any) ([]reflect.StructField, error) {
	innerStruct := []reflect.StructField{}

	for {
		//fetch first key in object
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return innerStruct, nil
			}
			return nil, fmt.Errorf("JSON parsing error: %v", err)
		}

		// Check if we've reached the end of an object
		if delim, ok := token.(json.Delim); ok && delim == '}' {
			return innerStruct, nil
		}

		// Fetch token name (will always be a string name of the first key)
		tokenName, tagName, err := returnFieldNameAndTag(token)
		if err != nil {
			return nil, fmt.Errorf("expected field name, got: %v", token)
		}

		// Fetch token value (either primitive or delimiter)
		valueToken, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("JSON parsing error: %v", err)
		}

		// check if the value is a delimiter (nested object/array) or primitive, and dispatch accordingly
		if delim, ok := valueToken.(json.Delim); ok {
			switch delim {
			case '{':
				objectField, err := processObject(tokenName, tagName, valueToken)
				if err != nil {
					return nil, err
				}
				innerStruct = append(innerStruct, objectField)
			case '[':
				arrayField, err := processArray(tokenName, tagName)
				if err != nil {
					return nil, err
				}
				innerStruct = append(innerStruct, arrayField)
			}
		} else {
			// Primitive value
			primitiveField := processPrimitive(tokenName, tagName, valueToken)
			innerStruct = append(innerStruct, primitiveField)
		}
	}
}

func traverseArray(token any) error {
	log.Fatal("The JSON document contains an array as a root level item, which is not supported.")
	return errors.New("The JSON document contains an array as a root level item, which is not supported.")
}
