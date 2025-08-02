package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	jsonStream "github.com/joan-cb/schema-stream/jsonStream"
	schemabuilder "github.com/joan-cb/schema-stream/schemaBuilder"
)

func main() {
	jsonData, err := os.ReadFile("test.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	types, err := jsonStream.ReturnStructDefinition(jsonData)
	if err != nil {
		fmt.Println("error parsing json input:", err)
		return
	}

	// Add this debug line
	fmt.Printf("DEBUG: Found %d struct fields\n", len(types))
	for i, field := range types {
		fmt.Printf("  Field %d: %s (type: %v)\n", i, field.Name, field.Type)
	}
	// Generate the struct definition
	fmt.Println("Generated Go struct definition:")
	log.Println(types)

	jsonSchema := schemabuilder.ReturnSchemaFromStructFields(types)
	if err != nil {
		fmt.Println("error generating schema:", err)
		return
	}
	schemaFile, err := os.Create("schema.json")
	if err != nil {
		fmt.Println("error creating schema file:", err)
		return
	}
	json.NewEncoder(schemaFile).Encode(jsonSchema)
	fmt.Println(jsonSchema)
}
