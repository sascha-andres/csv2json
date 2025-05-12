# csv2json

A command-line tool that converts CSV data to JSON, YAML, or TOML format using a configurable mapping.

[![CI/CD](https://github.com/sascha-andres/csv2json/actions/workflows/ci.yml/badge.svg)](https://github.com/sascha-andres/csv2json/actions/workflows/ci.yml)

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
| `-nested-property` | `data` | Property name for nested array output. When specified, array output is nested under this property name. |

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
  },
  "calculated": [
    {
      "property": "calculatedProperty",
      "kind": "kindOfCalculation",
      "format": "formatString",
      "type": "dataType"
    }
  ],
  "extra_variables": {
    "variable-name": {
      "value": "variable-value"
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

### Calculated Fields

Calculated fields allow you to add dynamic values to your output that are not directly derived from the CSV input. These fields are defined in the `calculated` array of the mapping configuration.

Each calculated field has the following properties:
- `property`: The name of the property in the output (supports dot notation for nested objects)
- `kind`: The type of calculation to perform (see below)
- `format`: Additional information for the calculation, varies by kind
- `type`: The data type of the calculated value (`int`, `float`, `bool`, or `string`)
- `location`: Where the calculated field should be applied - either `record` (default) or `document`

#### Kinds of Calculated Fields

1. **datetime**: Adds the current date/time formatted according to the format string
   - `format`: A Go time format string (e.g., "2006-01-02" for date, "15:04:05" for time)

2. **application**: Adds application-specific values
   - `format`: Currently only supports "record", which adds the record index (0-based)

3. **environment**: Adds the value of an environment variable
   - `format`: The name of the environment variable to read

4. **extra**: Adds the value of an extra variable defined in the configuration
   - `format`: The name of the extra variable to use
   - Extra variables are defined in the `extra_variables` section of the configuration

5. **mapping**: Maps values from a source field to different output values
   - `format`: Specified as "field:mapping_list" where:
     - `field` is the source field name (when using `-named`) or index
     - `mapping_list` is a comma-separated list of "from=to" pairs
     - A special "default" mapping can be specified for values that don't match any explicit mapping
   - Example: "from-to:a=0,b=1,default=99" maps values from the "from-to" field:
     - "a" becomes 0
     - "b" becomes 1
     - Any other value becomes 99

#### Record-Level vs Document-Level Calculated Fields

Calculated fields can be applied at two different levels:

1. **Record-Level Fields** (`location: "record"`): 
   - Applied to each individual record in the output
   - This is the default if no location is specified
   - Always included regardless of output format

2. **Document-Level Fields** (`location: "document"`):
   - Applied to the entire document, not to individual records
   - Only applied when using array output (with `-array` flag) or when using TOML/YAML output formats
   - Typically used for metadata about the entire document
   - Often placed under a top-level property like `_meta`

Document-level calculated fields are useful for adding metadata about the entire dataset, such as:
- Total number of records processed
- Processing timestamp
- Global configuration values

**Note:** Document-level calculated fields are only applied when the output is a single document containing all records (array mode). They are not applied when outputting individual records as separate JSON objects.

#### Example

```json
{
  "mapping": {
    "id": {
      "property": "productId",
      "type": "int"
    }
  },
  "calculated": [
    {
      "property": "metadata.recordNumber",
      "kind": "application",
      "format": "record",
      "type": "int",
      "location": "record"
    },
    {
      "property": "metadata.processedDate",
      "kind": "datetime",
      "format": "2006-01-02",
      "type": "string",
      "location": "record"
    },
    {
      "property": "metadata.processedTime",
      "kind": "datetime",
      "format": "15:04:05",
      "type": "string",
      "location": "record"
    },
    {
      "property": "metadata.userHome",
      "kind": "environment",
      "format": "HOME",
      "type": "string",
      "location": "record"
    },
    {
      "property": "metadata.version",
      "kind": "extra",
      "format": "app-version",
      "type": "string",
      "location": "record"
    },
    {
      "property": "_meta.totalRecords",
      "kind": "application",
      "format": "records",
      "type": "int",
      "location": "document"
    },
    {
      "property": "_meta.processedAt",
      "kind": "datetime",
      "format": "2006-01-02 15:04:05",
      "type": "string",
      "location": "document"
    }
  ],
  "extra_variables": {
    "app-version": {
      "value": "1.0.0"
    }
  }
}
```

This configuration would add the following calculated fields:

Record-level fields (added to each record):
- `metadata.recordNumber`: The 0-based index of the record
- `metadata.processedDate`: The current date in YYYY-MM-DD format
- `metadata.processedTime`: The current time in HH:MM:SS format
- `metadata.userHome`: The value of the HOME environment variable
- `metadata.version`: The string "1.0.0" from the extra variable "app-version"

Document-level fields (added to the top-level document when using array output):
- `_meta.totalRecords`: The total number of records processed
- `_meta.processedAt`: The date and time when the document was processed

## Output Behavior

### Without `-array` Flag (Default)

When processing multiple CSV rows without the `-array` flag, each row is converted to a separate JSON document and written to the output with newlines between them. This produces a newline-delimited JSON format (NDJSON/JSON Lines), where each line is a valid JSON object, but the file as a whole is not a standard JSON array.

### With `-array` Flag

When using the `-array` flag, all rows are collected into a single array and output as one document.

### Nested Property Output

When using the `-nested-property` flag with the `-array` flag (or when using TOML output which implicitly enables array mode), the output data is nested under the specified property name:

#### JSON Output with Nested Property

```
csv2json -in products.csv -array -nested-property items
```

Will produce:

```json
{
  "items": [
    {
      "property1": 1,
      "property2": {
        "property3": "hello"
      }
    },
    {
      "property1": 2,
      "property2": {
        "property3": "world"
      }
    }
  ]
}
```

#### YAML Output with Nested Property

```
csv2json -in products.csv -output-type yaml -nested-property items
```

Will produce:

```yaml
items:
  - property1: 1
    property2:
      property3: hello
  - property1: 2
    property2:
      property3: world
```

#### TOML Output Format

When using TOML as the output format, the array data is always wrapped in a property. By default, this property is named "data", but you can customize it using the `-nested-property` flag:

```
csv2json -in products.csv -output-type toml -nested-property items
```

Will produce:

```toml
[_meta]
  processedAt = "2023-05-12 15:30:45"
  totalRecords = 2

[[items]]
property1 = 1
property2 = { property3 = "hello" }

[[items]]
property1 = 2
property2 = { property3 = "world" }
```

Note how the document-level calculated fields appear in the `_meta` section at the top of the document, while record-level calculated fields would appear within each record.

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

### Value Mapping Example

Given the following CSV:

```csv
id,status,value
1,"active",10.5
2,"inactive",20.3
3,"pending",15.7
```

And this mapping.json with value mapping:

```json
{
  "mapping": {
    "id": {
      "property": "id",
      "type": "int"
    },
    "status": {
      "property": "originalStatus",
      "type": "string"
    },
    "value": {
      "property": "amount",
      "type": "float"
    }
  },
  "calculated": [
    {
      "property": "statusCode",
      "kind": "mapping",
      "format": "status:active=1,inactive=0,pending=2,default=-1",
      "type": "int",
      "location": "record"
    }
  ]
}
```

Running with the `-named` flag:

```
csv2json -in statuses.csv -named
```

Will produce:

```
{"id":1,"originalStatus":"active","amount":10.5,"statusCode":1}
{"id":2,"originalStatus":"inactive","amount":20.3,"statusCode":0}
{"id":3,"originalStatus":"pending","amount":15.7,"statusCode":2}
```

This example demonstrates how to map string status values to numeric codes using the value mapping feature.

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

## Development

### CI/CD Pipeline

This project uses GitHub Actions for continuous integration and delivery:

- **Automated Testing**: All commits and pull requests trigger automated tests to ensure code quality.
- **Automated Releases**: When a tag with format `v*` (e.g., `v1.0.0`) is pushed, the CI pipeline automatically:
  1. Runs all tests
  2. Builds binaries for multiple platforms (Linux, macOS, Windows)
  3. Creates a GitHub release with these binaries

### Releases

Releases are available on the [GitHub Releases page](https://github.com/sascha-andres/csv2json/releases). Each release includes pre-built binaries for:

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

To create a new release:

1. Tag the commit with a version number: `git tag v1.0.0`
2. Push the tag: `git push origin v1.0.0`
3. The GitHub Actions workflow will automatically create the release with binaries

## License

This project is licensed under a Custom Commercial Use License.

- You are free to use this software for commercial purposes within your organization
- You must give appropriate credit and provide a link to the license
- You may not sell this software or any derivative works as part of a product outside your organization
- You may make modifications for internal use only, but may not distribute modified versions
- Distribution of this software as part of a commercial product to external customers is strictly prohibited unless explicitly permitted in writing by the owner

For more details, see the [LICENSE](LICENSE) file in the repository.

## Contributing

We welcome contributions to csv2json! Before contributing, please read our [Contributing Guidelines](CONTRIBUTING.md) which includes our Contributor License Agreement (CLA).

### Contributor License Agreement

All contributors to this project must agree to our Contributor License Agreement (CLA). The CLA ensures that the project maintainers have the necessary rights to use and distribute your contributions.

When you submit a pull request, our automated system will check if you have agreed to the CLA. If not, you'll be prompted to do so by commenting on the pull request with:

```
I agree to the CLA
```

For more details about the CLA and the contribution process, please see the [CONTRIBUTING.md](CONTRIBUTING.md) file.
