# Schema Stream

A command-line tool that automatically generates JSON schemas from JSON documents. This interactive CLI application takes a JSON file as input and produces a corresponding JSON schema with the ability to add custom descriptions to properties.

## Features

- 🚀 **Interactive CLI**: User-friendly command-line interface with colorful output
- 📁 **Auto-discovery**: Automatically finds JSON files in the current directory
- 🔍 **Parsing**: Converts JSON documents to structured schemas
- ✏️ **Editable properties**: Add custom descriptions to schema properties
- 💾 **Output generation**: Saves the generated schema to `outputSchema.json`

## Installation

### Prerequisites

- Go 1.24.4 or higher

### Build from source

1. Clone the repository:
```bash
git clone <repository-url>
cd schema-stream
```

2. Build the application:
```bash
go build -o schema-stream
```

3. (Optional) Add to your PATH:
```bash
# On macOS/Linux
export PATH=$PATH:$(pwd)
```

## Usage

### Basic Usage

1. **Navigate to a directory containing JSON files**:
```bash
cd /path/to/your/json/files
```

2. **Run the application**:
```bash
./schema-stream
```

3. **Follow the interactive prompts**:
   - The app will show you available JSON files in the current directory
   - Select the JSON file you want to process
   - Optionally add descriptions to schema properties
   - The schema will be saved as `outputSchema.json`

### Example Workflow

1. **Prepare your JSON file** (e.g., `data.json`):
```json
{
    "name": "John",
    "age": 30,
    "email": "john@example.com",
    "address": {
        "street": "123 Main St",
        "city": "Anytown"
    }
}
```

2. **Run the CLI**:
```bash
./schema-stream
```

3. **Select your file** when prompted

4. **Add descriptions** (optional):
   - Choose "yes" when asked about adding descriptions
   - Enter property names and descriptions
   - Choose "no" when finished

5. **Find your schema** in `outputSchema.json`

## Supported Input

- **JSON objects**: The tool works best with JSON documents whose top-level property is an object
- **Nested structures**: Supports complex nested objects and arrays
- **Data types**: Automatically detects and maps JSON types to schema types

## Output

The application generates a JSON schema file (`outputSchema.json`) that includes:
- Property definitions with correct types
- Nested object structures
- Array definitions
- Custom descriptions (if added during the interactive process)
