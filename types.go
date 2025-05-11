package csv2json

type (

	// Mapper defines a structure for mapping input data to output data, applying configuration and marshaling as needed.
	Mapper struct {

		// In specifies the input file path or '-' for standard input in the mapping process.
		In string

		// Out specifies the output file path or '-' for standard output in the mapping process.
		Out string

		// Array indicates whether the JSON output should be wrapped in an array format during the mapping process.
		Array bool

		// Named indicates whether the CSV header should be used for mapping column names to JSON fields.
		Named bool

		// MappingFile specifies the path to a JSON file containing the mapping configuration for the data transformation process.
		MappingFile string

		// MarshalWith specifies a custom Marshaler to be used for serializing data during the mapping process. json, yaml or toml
		MarshalWith string

		// TomlPropertyName specifies the property name to use for TOML array output (defaults to "data")
		TomlPropertyName string

		// marshaler defines a custom function for serializing a value of any type into a byte slice with error handling.
		marshaler func(v any) ([]byte, error)

		// configuration holds the mapping configuration used during the data transformation process.
		configuration Configuration
	}

	// ColumnConfiguration defines the structure for configuring a column's property and type in a mapping.
	ColumnConfiguration struct {

		// Property specifies the name of the column property in the mapping configuration.
		Property string `json:"property"`

		// Type specifies the data type of the column in the mapping configuration.
		Type string `json:"type"`
	}

	// Configuration represents a mapping configuration where keys map to ColumnConfiguration structures.
	Configuration struct {

		// Mapping represents a map of keys to their corresponding column configurations in the mapping structure.
		Mapping map[string]ColumnConfiguration `json:"mapping"`
	}
)
