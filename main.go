package main

import (
	"fmt"
	"os"

	"github.com/joan-cb/schema-stream/schemaStream"
)

func main() {
	jsonData, err := os.ReadFile("test.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	types, err := schemaStream.ReturnStructDefinition(jsonData)
	if err != nil {
		fmt.Println("error parsing json input:", err)
		return
	}
	for key, value := range types {
		fmt.Println(key, value)
	}
	// Generate the struct definition
	fmt.Println("Generated Go struct:")

	for k, v := range types {
		fmt.Println(k, v)
	}
}
