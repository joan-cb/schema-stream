package jsonStream

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
)

// var decoder *json.Decoder

// to do add global to access current token
//pass decoder in each func

func ReturnStructDefinition(jsonData []byte) ([]reflect.StructField, error) {
	reader := strings.NewReader(string(jsonData))
	decoder := json.NewDecoder(reader)
	outerStruct := []reflect.StructField{}

	token, err := decoder.Token()
	udpateCurrentToken(token)

	if err != nil {
		if err == io.EOF {
			return outerStruct, errors.New("the provided file is empty.")
		}
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}

	tokenType, err := returnTokenType(token)
	if err != nil {
		return []reflect.StructField{}, err
	}

	switch tokenType {
	case tokenOpenSquareBracket:
		handleArray(&outerStruct, decoder)
	case tokenOpenCurlyBrace:
		handleObject(&outerStruct, decoder)
	default:
		return []reflect.StructField{}, errors.New(fmt.Sprintf("unexpected token type for first token, expect [ or {. Got %v", token))
	}

	return outerStruct, nil
}

func returnTokenType(token json.Token) (tokenType, error) {
	if delim, ok := token.(json.Delim); ok {
		switch delim {
		case '[':
			return tokenOpenSquareBracket, nil
		case ']':
			return tokenCloseSquareBracket, nil
		case '{':
			return tokenOpenCurlyBrace, nil
		case '}':
			return tokenCloseCurlyBrace, nil

		}
	}

	switch reflect.TypeOf(token).Kind() {
	case reflect.String:
		return tokenString, nil
	case reflect.Int64, reflect.Float64:
		return tokenNumber, nil
	case reflect.Bool:
		return tokenBool, nil
		// TO DO HANDLE NIL SPECIFICALLY
	default:
		return unexpectedToken, errors.New("nil or invalid token")

	}

}

// func traverseJSON() ([]reflect.StructField, error) {
// 	token, err := decoder.T
// oken()
// 	if err != nil {
// 		return nil, fmt.Errorf("JSON parsing error: %v", err)
// 	}

// 	switch token.(json.Delim) {
// 	case '[':
// 		err := traverseArray(token)
// 		if err != nil {
// 			return nil, fmt.Errorf("JSON parsing error: %v", err)
// 		}
// 		return []reflect.StructField{}, errors.New("The JSON document contains an array as a root level item, which is not supported.")
// 	case '{':
// 		structFields, err := traverseObject()
// 		if err != nil {
// 		}
// 		return structFields, nil
// 	default:
// 		return nil, fmt.Errorf("unexpected token: %v", token)
// 	}
// }

// func traverseObject() ([]reflect.StructField, error) {
// 	innerStruct := []reflect.StructField{}

// 	for {
// 		//fetch first key in object
// 		token, err := decoder.Token()
// 		if err != nil {
// 			if err == io.EOF {
// 				return innerStruct, nil
// 			}
// 			return nil, fmt.Errorf("JSON parsing error: %v", err)
// 		}

// 		// Check if we've reached the end of an object
// 		if delim, ok := token.(json.Delim); ok && delim == '}' {
// 			return innerStruct, nil
// 		}

// 		// Fetch token name (will always be a string name of the first key)
// 		tokenName, tagName, err := returnFieldNameAndTag(token)
// 		if err != nil {
// 			return nil, fmt.Errorf("expected field name, got: %v", token)
// 		}

// 		// Fetch token value (either primitive or delimiter)
// 		valueToken, err := decoder.Token()
// 		if err != nil {
// 			return nil, fmt.Errorf("JSON parsing error: %v", err)
// 		}

// 		// check if the value is a delimiter (nested object/array) or primitive, and dispatch accordingly
// 		if delim, ok := valueToken.(json.Delim); ok {
// 			switch delim {
// 			case '{':
// 				objectField, err := processObject(tokenName, tagName)
// 				if err != nil {
// 					return nil, err
// 				}
// 				innerStruct = append(innerStruct, objectField)
// 			case '[':
// 				arrayField, err := processArray(tokenName, tagName)
// 				if err != nil {
// 					return nil, err
// 				}
// 				innerStruct = append(innerStruct, arrayField)
// 			}
// 		} else {
// 			// Primitive value
// 			primitiveField := processPrimitive(tokenName, tagName, valueToken)
// 			innerStruct = append(innerStruct, primitiveField)
// 		}
// 	}
// }

func traverseArray(token any) error {
	log.Fatal("The JSON document contains an array as a root level item, which is not supported.")
	return errors.New("The JSON document contains an array as a root level item, which is not supported.")
}
