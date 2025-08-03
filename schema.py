import json
from genson import SchemaBuilder

def generate_schema(json_data):
    builder = SchemaBuilder()
    builder.add_object(json_data)
    return builder.to_schema()

def main():
    input_file = "test.json"
    output_file = "output_schema.json"

    # Load the JSON data
    with open(input_file, "r") as f:
        data = json.load(f)

    # Generate schema
    schema = generate_schema(data)

    # Save the schema to a file
    with open(output_file, "w") as f:
        json.dump(schema, f, indent=2)

    print(f"Schema saved to {output_file}")

if __name__ == "__main__":
    main()
