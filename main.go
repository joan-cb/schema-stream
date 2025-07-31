package main

import (
	"fmt"
	"os"

	"github.com/joan-cb/schema-stream/schemaStream"
)

func main() {
	jsonData, err := os.ReadFile("2083.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	types, err := schemaStream.ParseAndGetTypes(jsonData)
	if err != nil {
		fmt.Println("Error parsing json input:", err)
		return
	}
	for key, value := range types {
		fmt.Println(key, value)
	}
	// Generate the struct definition
	fmt.Println("Generated Go struct:")
	fmt.Println(types)
}
