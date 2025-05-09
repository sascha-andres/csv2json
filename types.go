package csv2json

type (
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
