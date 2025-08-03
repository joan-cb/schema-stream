package cli

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/alperdrsnn/clime"
)

func Run() {
	getBox("Welcome to schema-stream CLI app", "This CLI tool takes a JSON document as its input and returns the corresponding JSON schema.", "Note that only schemas whose top-level property is an object are supported.", "The app will start with looking for JSON files in the current directory.", "You can also edit the descriptions of the properties in the schema after it is generated.")

	// Interactive input
	jsonFilesInCurrentDirectory, err := findJSONFilesInDirectory(".")
	if err != nil {
		log.Println(err)
	}

	getSpinner("Searching for JSON files in the current directory...")
	printRainbow("=============================================================")
	printBold("Available JSON files in the current directory:")

	for _, file := range jsonFilesInCurrentDirectory {
		printInfo(file)
	}
	printRainbow("=============================================================")
	var inputFile string
	for {
		inputFile, _ = clime.Ask("Enter the name of the JSON file to process")
		if inputFile == "" || !slices.Contains(jsonFilesInCurrentDirectory, inputFile) {
			fmt.Println(clime.Error.Sprint("Incorrect file name. Please try again."))
		} else {
			break
		}
	}

	// Validate input file is provided
	if inputFile == "" {
		printError("No input file provided. Please provide a valid JSON file.")
		os.Exit(1)
	}

	// Run the main processing logic
	if err := processJSONToSchema(inputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	clime.SuccessBanner("The schema returned with the following properties.")
	schemaTable := getSchemaTable(schema)
	schemaTable.Println()
	getBox("Schema Properties", "You can edit the properties in the schema if you want.", "Editable properties are all columns in the table except name and type")

outerLoop:
	for {
		name, err := clime.Ask("Do you want to add descriptions to any properties? (yes/no)")
		if err != nil {
			clime.ErrorBanner("Error reading title!")
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}

		switch name {
		case "yes":
			key, err := clime.Ask("Enter the name of the property to edit")
			if err != nil {
				log.Println(err)
			}

			property, err := clime.Ask("Enter the property type")
			if err != nil {
				log.Println(err)
			}
			description, err := clime.Ask("Enter the description for the property")
			if err != nil {
				log.Println(err)
			}
			err = editKey(key, property, description)
			if err != nil {
				clime.ErrorBanner(fmt.Sprintf("Error editing key %s: %v", key, err))
				log.Println(err)
			} else {
				clime.SuccessBanner(fmt.Sprintf("Key %s edited successfully!", key))
			}
			schemaTable := getSchemaTable(schema)
			schemaTable.Println()
		case "no":
			clime.SuccessBanner("Moving to save schama!")
			break outerLoop

		}
	}
	getProgressBar("Saving schema file...")
	err = writeSchemaFile(inputFile, schema)
	if err != nil {
		log.Println(err)
		clime.ErrorBanner(fmt.Sprintf("Error saving schema file: %v", err))
	}
	getSpinner("Saving schema file...")
	clime.SuccessBanner("Schema file saved successfully!")
}
