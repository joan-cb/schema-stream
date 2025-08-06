package jsonStream

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/Microsoft/go-winio/pkg/process"
)

var previousToken json.Token

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
	unexpectedToken
)

func dispatchFirstToken(outerStruct *[]reflect.StructField)  error {
	
	
	if err != nil {
		if err == io.EOF {
				return outerStruct, nil
			}
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}
	if delim, ok := token.(json.Delim); ok && delim == '{' {
		err = handleArray(outerStruct)
		if err != nil{
			return err
		}
	}

	return nil
}	


func udpatepreviousToken(token json.Token){
	previousToken = token
}




func handleObject(outerStruct *[]reflect.StructField, decoder *json.Decoder) error {
	nextToken, err := decoder.Token()
	udpatepreviousToken(token)
	
	tokenType, err := returnTokenType(tokenType)
	if err != nil{
		return  err
	}
	 switch tokenType{
	 case tokenString:
		handleString(outerStruct,decoder)
	 case  tokenOpenSquareBracket:
		handleArray(outerStruct, decoder)
	 default:
		return errors.New(fmt.Sprintf("unexpected token type for first token, expect [ or {. Got %v", token))
	}
	return  nil
}


func handleString(outerStruct *[]reflect.StructField, decoder *json.Decoder) error {	
	nextToken, err := decoder.Token()
	udpatepreviousToken(token)
		
	tokenType, err := returnTokenType(tokenType)
	if err != nil{
		return  err
	}
	switch tokenType{
	case tokenString:
		handleString(outerStruct,decoder)
	case  tokenOpenSquareBracket:
		handleArray(outerStruct, decoder)
	default:
		return errors.New(fmt.Sprintf("unexpected token type for first token, expect [ or {. Got %v", token))
	}
	return  nil
}



func handleArray(outerStruct *[]reflect.StructField, decode *json.Decoder) error {
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


func processObject (token *json.Token, outerStruct *[]reflect.StructField, decoder *json.Decoder) error {
	innerStruct := []reflect.StructField{}
	for{
		token, err :=decode.Token()
		if err!=nil {
			return err
		}
		udpateCurrentToken(token)
		err = processPrimitive(token, outerStruct,decoder)
		if err !=nil {
			return err
		}
		innerStruct = append(innerStruct, )
		}

	}
}


func processPrimitive(token *json.Token, outerStruct *[]reflect.StructField, decoder *json.Decoder) error  {
	tokenName, tagName, err := returnFieldNameAndTag(previousToken)
	if err !=nil {
		return err 
	}

	return reflect.StructField{
		Name: sanitizeTokenName(tokenName),
		Type: reflect.TypeOf(token),
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
