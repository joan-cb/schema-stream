package jsonStream

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReturnStructDefinition(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	path := "./test"
	err := filepath.WalkDir(path, func(filePath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(filePath), ".json") {
			jsonData, readErr := os.ReadFile(filePath)
			if readErr != nil {
				t.Errorf("Failed to read file %s: %v", filePath, readErr)
				return nil // Continue with other files
			}

			structFields, parseErr := ReturnStructDefinition(jsonData)
			if parseErr != nil {
				t.Errorf("Failed to return struct definition for %s: %v", filePath, parseErr)
				return nil // Continue with other files
			}

			// Basic validation that we got some fields
			if len(structFields) == 0 {
				t.Errorf("No struct fields returned for %s", filePath)
			}

			t.Logf("Successfully parsed %s with %d fields", filePath, len(structFields))
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk directory: %v", err)
	}
}
