# csv2json

A command-line tool that converts CSV data to JSON, YAML, or TOML format using a configurable mapping.

## Overview

csv2json reads CSV data from a file or standard input, transforms it according to a mapping configuration, and outputs the result in JSON, YAML, or TOML format to a file or standard output.

## Command-Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-in` | `-` (stdin) | Input file path. Use `-` for standard input. |
| `-out` | `-` (stdout) | Output file path. Use `-` for standard output. |
| `-array` | `false` | Output all records as a single array instead of separate documents. |
| `-named` | `false` | Use CSV header row for column names instead of numeric indices. |
| `-mapping` | `mapping.json` | Path to the mapping configuration file. |
| `-output-type` | `json` | Output format type. One of: `json`, `yaml`, or `toml`. |
| `-toml-property` | `data` | Property name for TOML array output. |

**Note:** When using `yaml` or `toml` as the output type, the `-array` flag is automatically set to `true`.

## Environment Variables

All flags can also be set using environment variables with the prefix `CSV2JSON_`. For example:
- `CSV2JSON_IN=input.csv` is equivalent to `-in input.csv`
- `CSV2JSON_ARRAY=true` is equivalent to `-array`

## Mapping Configuration

The mapping configuration is a JSON file that defines how CSV columns are mapped to properties in the output. The configuration has the following structure:

```json
{
  "mapping": {
    "key1": {
      "property": "propertyName",
      "type": "dataType"
    },
    "key2": {
      "property": "nested.property",
      "type": "dataType"
    }
  }
}
```

Where:
- `key1`, `key2`, etc. are either:
  - Column indices (0, 1, 2, ...) when not using the `-named` flag
  - Column names from the CSV header when using the `-named` flag
- `propertyName` is the name of the property in the output
- `nested.property` demonstrates how to create nested objects using dot notation
- `dataType` is one of:
  - `int` - converts the value to an integer
  - `float` - converts the value to a floating-point number
  - `bool` - converts the value to a boolean
  - `string` (default) - keeps the value as a string

## Output Behavior

### Without `-array` Flag (Default)

When processing multiple CSV rows without the `-array` flag, each row is converted to a separate JSON document and written to the output with newlines between them. This produces a newline-delimited JSON format (NDJSON/JSON Lines), where each line is a valid JSON object, but the file as a whole is not a standard JSON array.

### With `-array` Flag

When using the `-array` flag, all rows are collected into a single array and output as one document.

### TOML Output Format

When using TOML as the output format, the array data is wrapped in a field specified by the `-toml-property` flag (defaults to "data"):

```toml
[[data]]
property1 = 1
# ...

[[data]]
property1 = 2
# ...
```

You can customize this property name using the `-toml-property` flag:

```
csv2json -in products.csv -output-type toml -toml-property items
```

Will produce:

```toml
[[items]]
property1 = 1
# ...

[[items]]
property1 = 2
# ...
```

## Examples

### Basic Usage

Given the following CSV:

```csv
1,"hello",2.3
```

And this mapping.json:

```json
{
  "mapping": {
    "0": {
      "property": "property1",
      "type": "int"
    },
    "1": {
      "property": "property2.property3",
      "type": "string"
    },
    "2": {
      "property": "property4",
      "type": "float"
    }
  }
}
```

The default output will be:

```json
{
  "property1": 1,
  "property2": {
    "property3": "hello"
  },
  "property4": 2.3
}
```

### Using Named Columns

Given the following CSV:

```csv
id,name,price
1,"Product A",19.99
2,"Product B",29.99
```

And this mapping.json:

```json
{
  "mapping": {
    "id": {
      "property": "productId",
      "type": "int"
    },
    "name": {
      "property": "productName",
      "type": "string"
    },
    "price": {
      "property": "pricing.retail",
      "type": "float"
    }
  }
}
```

Running with the `-named` flag:

```
csv2json -in products.csv -named
```

Will produce (in NDJSON format, with each line being a separate JSON document):

```
{"productId":1,"productName":"Product A","pricing":{"retail":19.99}}
{"productId":2,"productName":"Product B","pricing":{"retail":29.99}}
```

Note: The actual output will not be pretty-printed but shown as compact JSON objects, one per line.

### Output as Array in YAML Format

```
csv2json -in products.csv -named -output-type yaml
```

Will produce:

```yaml
- productId: 1
  productName: Product A
  pricing:
    retail: 19.99
- productId: 2
  productName: Product B
  pricing:
    retail: 29.99
```
