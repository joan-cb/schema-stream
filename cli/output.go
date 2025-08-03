package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
)

func writeSchemaFile(filename string, schema *jsonschema.Schema) error {
	file, err := os.Create("outputSchema.json")
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Pretty print

	if err := encoder.Encode(schema); err != nil {
		return fmt.Errorf("failed to encode schema: %w", err)
	}

	fmt.Printf("Schema written to outputSchema.json")
	return nil
}
