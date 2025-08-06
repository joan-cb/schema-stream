package jsonStream

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
)

type tokenType int

const (
	tokenOpenCurlyBrace tokenType = iota
	tokenCloseCurlyBrace
	tokenOpenSquareBracket
	tokenCloseSquareBracket
	tokenString
	tokenNumber
	tokenBool
	tokenNull
)

// can remove for loop as this will go in the  handle object or handle array
func traverseOuterObject() ([]reflect.StructField, error) {
	outerStruct := []reflect.StructField{}
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return outerStruct, nil
			}
			return nil, fmt.Errorf("JSON parsing error: %v", err)
		}
		
		if delim, ok := token.(json.Delim); ok {
			switch delim {
			case '{':
				return handleObject(&outerStruct)
			case '[':
				return handleArray( &outerStruct)
			default:
				return fmt.Errorf("unexpected delimiter: %v", delim)
			}
		}
	}	
}

  

func handleObject(*[]reflect.StructField) error {
	innerStruct := reflect.StructField{}
	nextToken, err := decoder.Token()
	if err != nil {
		return err
	}


	return nil
}

func handleArray(outerStruct *[]reflect.StructField) error {
	innerStruct := reflect.StructField{}

	nextToken, err := decoder.Token()
	if err != nil {
		return err
	}
	switch t := nextToken.(type) {
	case json.Delim:
		switch t {
		case '[':
			err := handleArray(tokenName, tagName, outerStruct)
			if err != nil {
				return err
			}
		case '{':
			// Process nested object
			err := handleObject(tokenName, tagName, outerStruct)
			if err != nil {
				return err
			}
		}
	case string, float64, bool, int64:
		processPrimitive(tokenName, tagName, t)

	}

	// Consume remaining tokens in the array until we reach the closing bracket
	for {
		token, err := decoder.Token()
		if err != nil {
			return reflect.StructField{}, err
		}
		if delim, ok := token.(json.Delim); ok && delim == ']' {
			break
		}

	*outerStruct = append(*outerStruct, innerStruct)
	return nil
}

func returnFieldNameAndTag(token any) (string, string, error) {
	if str, ok := token.(string); ok {
		tokenName := sanitizeTokenName(str)
		tagName := strings.ToLower(tokenName[:1]) + tokenName[1:]
		return tokenName, tagName, nil
	}
	return "", "", fmt.Errorf("unexpected token type: %T", token)
}

 





func processPrimitive(tokenName, tagName string, valueToken any) (reflect.StructField, error) {
	return reflect.StructField{
		Name: sanitizeTokenName(tokenName),
		Type: reflect.TypeOf(valueToken),
		Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s" jsonschema:"required"`, tagName)),
	}, nil
}

func sanitizeTokenName(key string) string {
	if key == "" {
		log.Println("Warning: Empty key provided, returning default name.")
		return "unamedField"
	}

	// Capitalize first letter
	sanitizedToken := strings.ToUpper(key[:1]) + key[1:]

	// Replace invalid characters with underscores
	invalidChars := []string{"-", " ", ".", ",", ":", ";", "!", "?", "@", "#", "$", "%", "^", "&", "*", "(", ")", "[", "]", "{", "}", "|", "\\", "/", "<", ">", "=", "+", "~", "`", "'", "\""}
	for _, char := range invalidChars {
		sanitizedToken = strings.ReplaceAll(sanitizedToken, char, "_")
	}

	// Remove multiple consecutive underscores
	for strings.Contains(sanitizedToken, "__") {
		sanitizedToken = strings.ReplaceAll(sanitizedToken, "__", "_")
	}

	// Remove leading/trailing underscores
	sanitizedToken = strings.Trim(sanitizedToken, "_")

	// Ensure it starts with a letter
	if sanitizedToken == "" || !isLetter(rune(sanitizedToken[0])) {
		sanitizedToken = "Field" + sanitizedToken
	}

	return sanitizedToken
}

// isLetter checks if a rune is a letter
func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
