package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/alperdrsnn/clime"
	"github.com/invopop/jsonschema"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func getSchemaTable(schema *jsonschema.Schema) *clime.Table {
	table := clime.NewTable().
		WithStyle(clime.TableStyleDefault).
		WithHeaderColor(clime.MagentaColor).
		WithBorderColor(clime.BlueColor).
		ShowBorders(true).
		ShowHeader(true).
		AutoResize(true).
		WithMaxWidth(200).
		AddColumn("name").
		AddColumn("type").
		AddColumn("description").
		AddColumn("required").
		AddColumn("examples").
		AddColumn("enum").
		AddColumn("items")
	for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
		key := pair.Key
		description := pair.Value.Description
		propertyType := pair.Value.Type
		required := pair.Value.Required
		examples := pair.Value.Examples
		enum := pair.Value.Enum
		items := pair.Value.Items
		table.AddRow(key, propertyType, description, handleSliceOfStringProperty(required), handleSliceOfAnyProperty(examples), handleSliceOfAnyProperty(enum), hasItems(items))
	}
	return table

}

func editKey(key, property, value string) error {
	for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
		if pair.Key == key {
			err := handlePropertySwitch(pair, key, property, value)
			if err != nil {
				return fmt.Errorf("error handling property switch: %w", err)
			}
			return nil
		}
	}
	return errors.New("key not found in schema properties")
}

func handleSliceOfStringProperty(input []string) string {
	var stringValue string
	for _, v := range input {
		stringValue = stringValue + ", " + v
	}
	// Remove the leading ", " if needed
	if len(stringValue) > 0 {
		stringValue = stringValue[2:]
	}
	return stringValue
}

func handleSliceOfAnyProperty(input []any) string {
	var stringValue string
	for _, v := range input {
		stringValue = stringValue + ", " + fmt.Sprintf("%v", v)
	}
	// Remove the leading ", " if needed
	if len(stringValue) > 0 {
		stringValue = stringValue[2:]
	}
	return stringValue
}

func handlePropertySwitch(pair *orderedmap.Pair[string, *jsonschema.Schema], key, property, value string) error {
	switch property {
	case "description":
		pair.Value.Description = value
		return nil
	case "required":
		_, err := strconv.ParseBool(value)
		if err != nil {
			return errors.New("invalid boolean value for 'required'")
		}
		pair.Value.Required = append(pair.Value.Required, value)
		return nil
	case "examples":
		pair.Value.Examples = append(pair.Value.Examples, value)
		return nil
	case "enum":
		pair.Value.Enum = append(pair.Value.Enum, value)
	default:
		return fmt.Errorf("unsupported property: %s", property)
	}
	return nil
}

func hasItems(schema *jsonschema.Schema) string {
	if schema == nil {
		return "false"
	}
	if schema.Items == nil {
		return "false"
	}
	return "true"
}
