package storer

type (
	// OutputType represents the output type of the engine (JSON, YAML or TOML)
	OutputType string

	// Project represents a project for the daemon mapped to
	// a relational database or serialized
	Project struct {

		// Id is the unique identifier for the project
		// xid is used by the daemon for create requests requiring an identifier
		Id string `json:"id"`

		// Name is a canonical name
		Name string `json:"name"`

		// Description can should be used to aid in understanding the
		// use casef for the project
		Description *string `json:"description"`

		// OutputType may be set to one of the supported output types (JSON, YAML or TOML)
		OutputType *OutputType `json:"output_type"`

		// Named should be true if the first line in the CSV file containes header
		Named *bool `json:"named"`

		// Array should be true if you need to return valid arrays, implicitly true
		// for YAML and TOML
		Array *bool `json:"array"`

		// NestedPropertyName can be used to name the property to which the array is assigned
		NestedPropertyName *string `json:"nested_property_name"`

		// Separator defaults to , and is passed to CSV parsing
		Separator *string `json:"separator"`
	}
)
