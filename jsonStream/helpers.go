package jsonStream

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func processArray(tokenName, tagName string) (reflect.StructField, error) {
	// Get the first element to determine array type
	firstElement, err := decoder.Token()
	if err != nil {
		return reflect.StructField{}, fmt.Errorf("JSON parsing error: %v", err)
	}

	var arrayType reflect.Type

	// Check if it's an empty array
	if delim, ok := firstElement.(json.Delim); ok && delim == ']' {
		// Empty array
		arrayType = reflect.TypeOf([]interface{}{})
	} else if delim == '{' {
		// Array of objects - need to process the object structure
		nestedFields, err := traverseObject(firstElement)
		if err != nil {
			return reflect.StructField{}, err
		}

		// Create the struct type for array elements
		elementStructType := reflect.StructOf(nestedFields)
		arrayType = reflect.SliceOf(elementStructType)

		// Skip remaining array elements
		for {
			token, err := decoder.Token()
			if err != nil {
				return reflect.StructField{}, err
			}
			if delim, ok := token.(json.Delim); ok && delim == ']' {
				break
			}
		}
	} else {
		// Create array type based on first element type
		switch firstElement.(type) {
		case string:
			arrayType = reflect.TypeOf([]string{})
		case float64:
			arrayType = reflect.TypeOf([]float64{})
		case bool:
			arrayType = reflect.TypeOf([]bool{})
		case nil:
			arrayType = reflect.TypeOf([]*interface{}{})
		default:
			arrayType = reflect.TypeOf([]interface{}{})
		}

		// Skip the rest of the array elements until we find ']'
		for {
			token, err := decoder.Token()
			if err != nil {
				return reflect.StructField{}, err
			}

			if delim, ok := token.(json.Delim); ok && delim == ']' {
				break // Found closing bracket, we're done
			}
			// Continue consuming tokens until we reach the end
		}
	}

	// Create struct field with array type
	return reflect.StructField{
		Name: tokenName,
		Type: arrayType,
		Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s" jsonschema:"required"`, tagName)),
	}, nil
}

func returnFieldNameAndTag(token any) (string, string, error) {
	if str, ok := token.(string); ok {
		tokenName := strings.ToUpper(str[:1]) + str[1:]
		tagName := strings.ToLower(str[:1]) + str[1:]
		return tokenName, tagName, nil
	}
	return "", "", fmt.Errorf("unexpected error: %v", token)
}

func processObject(tokenName, tagName string, valueToken any) (reflect.StructField, error) {
	// Recursively call traverseObject to get nested fields
	nestedFields, err := traverseObject(valueToken)
	if err != nil {
		return reflect.StructField{}, err
	}

	// Create struct type from nested fields
	nestedStructType := reflect.StructOf(nestedFields)

	// Return the struct field for the nested object
	return reflect.StructField{
		Name: tokenName,
		Type: nestedStructType,
		Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s" jsonschema:"required"`, tagName)),
	}, nil
}

func processPrimitive(tokenName, tagName string, valueToken any) reflect.StructField {
	return reflect.StructField{
		Name: tokenName,
		Type: reflect.TypeOf(valueToken),
		Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s" jsonschema:"required"`, tagName)),
	}
}
