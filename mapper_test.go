package csv2json

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

// TestCreate tests the creation of Mapper instances using various configurations and validates expected errors or outcomes.
func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		options []OptionFunc
		want    *Mapper
		wantErr bool
	}{
		{
			name: "basic mapper",
			options: []OptionFunc{
				WithIn("input.csv"),
				WithOut("output.json"),
			},
			want: &Mapper{
				In:  "input.csv",
				Out: "output.json",
			},
			wantErr: false,
		},
		{
			name: "mapper with array and named",
			options: []OptionFunc{
				WithIn("input.csv"),
				WithOut("output.json"),
				WithArray(true),
				WithNamed(true),
			},
			want: &Mapper{
				In:    "input.csv",
				Out:   "output.json",
				Array: true,
				Named: true,
			},
			wantErr: false,
		},
		{
			name: "empty input error",
			options: []OptionFunc{
				WithIn(""),
				WithOut("output.json"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty output error",
			options: []OptionFunc{
				WithIn("input.csv"),
				WithOut(""),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMapper(tt.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMapper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !compareMappers(got, tt.want) {
				t.Errorf("NewMapper() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMap tests the Map method with various configurations.
func TestMap(t *testing.T) {
	// Create a temporary mapping file
	mappingJSON := `{
		"mapping": {
			"id": {
				"property": "property1",
				"type": "int"
			},
			"text": {
				"property": "property2.property3",
				"type": "string"
			},
			"value": {
				"property": "property4",
				"type": "float"
			},
			"b": {
				"property": "property2.property5",
				"type": "bool"
			}
		}
	}`

	tempMappingFile, err := os.CreateTemp("", "mapping*.json")
	if err != nil {
		t.Fatalf("Failed to create temp mapping file: %v", err)
	}
	defer os.Remove(tempMappingFile.Name())

	if _, err := tempMappingFile.Write([]byte(mappingJSON)); err != nil {
		t.Fatalf("Failed to write to temp mapping file: %v", err)
	}
	tempMappingFile.Close()

	// Create a temporary CSV file
	csvData := `id,text,value,b
1,hello,2.3,true
2,world,3.4,false`

	tempCSVFile, err := os.CreateTemp("", "input*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(tempCSVFile.Name())

	if _, err := tempCSVFile.Write([]byte(csvData)); err != nil {
		t.Fatalf("Failed to write to temp CSV file: %v", err)
	}
	tempCSVFile.Close()

	// Create a CSV file with invalid data for error testing
	invalidCSVData := `id,text,value,b
1,hello,not_a_float,true`

	tempInvalidCSVFile, err := os.CreateTemp("", "invalid_input*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp invalid CSV file: %v", err)
	}
	defer os.Remove(tempInvalidCSVFile.Name())

	if _, err := tempInvalidCSVFile.Write([]byte(invalidCSVData)); err != nil {
		t.Fatalf("Failed to write to temp invalid CSV file: %v", err)
	}
	tempInvalidCSVFile.Close()

	// Create an invalid mapping file for error testing
	invalidMappingJSON := `{
		"mapping": {
			"missing_column": {
				"property": "property1",
				"type": "int"
			}
		}
	}`

	tempInvalidMappingFile, err := os.CreateTemp("", "invalid_mapping*.json")
	if err != nil {
		t.Fatalf("Failed to create temp invalid mapping file: %v", err)
	}
	defer os.Remove(tempInvalidMappingFile.Name())

	if _, err := tempInvalidMappingFile.Write([]byte(invalidMappingJSON)); err != nil {
		t.Fatalf("Failed to write to temp invalid mapping file: %v", err)
	}
	tempInvalidMappingFile.Close()

	tests := []struct {
		name           string
		options        []OptionFunc
		expectedOutput string
		wantErr        bool
	}{
		{
			name: "json output with named columns",
			options: []OptionFunc{
				WithIn(tempCSVFile.Name()),
				WithOut("-"),
				WithNamed(true),
				WithArray(true),
				WithMappingFile(tempMappingFile.Name()),
				WithOutputType("json"),
			},
			expectedOutput: `[{"property1":1,"property2":{"property3":"hello","property5":true},"property4":2.3},{"property1":2,"property2":{"property3":"world","property5":false},"property4":3.4}]`,
			wantErr:        false,
		},
		{
			name: "json output without array",
			options: []OptionFunc{
				WithIn(tempCSVFile.Name()),
				WithOut("-"),
				WithNamed(true),
				WithArray(false),
				WithMappingFile(tempMappingFile.Name()),
				WithOutputType("json"),
			},
			expectedOutput: `{"property1":1,"property2":{"property3":"hello","property5":true},"property4":2.3}
{"property1":2,"property2":{"property3":"world","property5":false},"property4":3.4}`,
			wantErr: false,
		},
		{
			name: "yaml output",
			options: []OptionFunc{
				WithIn(tempCSVFile.Name()),
				WithOut("-"),
				WithNamed(true),
				WithArray(true),
				WithMappingFile(tempMappingFile.Name()),
				WithOutputType("yaml"),
			},
			expectedOutput: "", // We won't check the exact YAML output format
			wantErr:        false,
		},
		{
			name: "toml output",
			options: []OptionFunc{
				WithIn(tempCSVFile.Name()),
				WithOut("-"),
				WithNamed(true),
				WithArray(true),
				WithMappingFile(tempMappingFile.Name()),
				WithOutputType("toml"),
				WithTomlPropertyName("records"),
			},
			expectedOutput: "", // We won't check the exact TOML output format
			wantErr:        false,
		},
		{
			name: "error - invalid float value",
			options: []OptionFunc{
				WithIn(tempInvalidCSVFile.Name()),
				WithOut("-"),
				WithNamed(true),
				WithArray(true),
				WithMappingFile(tempMappingFile.Name()),
				WithOutputType("json"),
			},
			expectedOutput: "",
			wantErr:        true,
		},
		{
			name: "error - missing column in mapping",
			options: []OptionFunc{
				WithIn(tempCSVFile.Name()),
				WithOut("-"),
				WithNamed(true),
				WithArray(true),
				WithMappingFile(tempInvalidMappingFile.Name()),
				WithOutputType("json"),
			},
			expectedOutput: "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect stdout to capture output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create and run mapper
			mapper, err := NewMapper(tt.options...)
			if err != nil {
				t.Fatalf("Failed to create mapper: %v", err)
			}

			err = mapper.Map()

			// Restore stdout and get output
			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := strings.TrimSpace(buf.String())

			// Check for errors
			if (err != nil) != tt.wantErr {
				t.Errorf("Map() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// For JSON output, normalize both expected and actual JSON for comparison
			if !tt.wantErr && strings.Contains(tt.name, "json") {
				// Check if the test case is for array output
				isArrayOutput := false
				for _, opt := range tt.options {
					// Create a temporary mapper to check if this option sets Array to true
					tempMapper := &Mapper{}
					opt(tempMapper)
					if tempMapper.Array {
						isArrayOutput = true
						break
					}
				}

				// If it's array output, we can directly compare JSON
				if isArrayOutput {
					var expected, actual interface{}
					json.Unmarshal([]byte(tt.expectedOutput), &expected)
					json.Unmarshal([]byte(output), &actual)

					if !reflect.DeepEqual(expected, actual) {
						t.Errorf("Map() output = %v, want %v", output, tt.expectedOutput)
					}
				} else {
					// For non-array output, split by newline and compare each JSON object
					expectedLines := strings.Split(tt.expectedOutput, "\n")
					actualLines := strings.Split(output, "\n")

					if len(expectedLines) != len(actualLines) {
						t.Errorf("Map() output has %d lines, want %d lines", len(actualLines), len(expectedLines))
						return
					}

					for i := range expectedLines {
						var expectedObj, actualObj interface{}
						json.Unmarshal([]byte(expectedLines[i]), &expectedObj)
						json.Unmarshal([]byte(actualLines[i]), &actualObj)

						if !reflect.DeepEqual(expectedObj, actualObj) {
							t.Errorf("Map() output line %d = %v, want %v", i+1, actualLines[i], expectedLines[i])
						}
					}
				}
			}
		})
	}
}

// TestInitialize tests the initialize method which reads the mapping file and opens input/output files.
func TestInitialize(t *testing.T) {
	// Create a temporary mapping file
	mappingJSON := `{
		"mapping": {
			"id": {
				"property": "property1",
				"type": "int"
			}
		}
	}`

	tempMappingFile, err := os.CreateTemp("", "mapping*.json")
	if err != nil {
		t.Fatalf("Failed to create temp mapping file: %v", err)
	}
	defer os.Remove(tempMappingFile.Name())

	if _, err := tempMappingFile.Write([]byte(mappingJSON)); err != nil {
		t.Fatalf("Failed to write to temp mapping file: %v", err)
	}
	tempMappingFile.Close()

	// Create a temporary CSV file
	csvData := `id
1`

	tempCSVFile, err := os.CreateTemp("", "input*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(tempCSVFile.Name())

	if _, err := tempCSVFile.Write([]byte(csvData)); err != nil {
		t.Fatalf("Failed to write to temp CSV file: %v", err)
	}
	tempCSVFile.Close()

	// Create a temporary output file
	tempOutFile, err := os.CreateTemp("", "output*.json")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer os.Remove(tempOutFile.Name())
	tempOutFile.Close()

	tests := []struct {
		name      string
		mapper    *Mapper
		wantErr   bool
		errString string
	}{
		{
			name: "successful initialization",
			mapper: &Mapper{
				In:          tempCSVFile.Name(),
				Out:         tempOutFile.Name(),
				MappingFile: tempMappingFile.Name(),
			},
			wantErr: false,
		},
		{
			name: "non-existent mapping file",
			mapper: &Mapper{
				In:          tempCSVFile.Name(),
				Out:         tempOutFile.Name(),
				MappingFile: "non_existent_file.json",
			},
			wantErr:   true,
			errString: "failed to read mapping file",
		},
		{
			name: "invalid mapping file",
			mapper: &Mapper{
				In:          tempCSVFile.Name(),
				Out:         tempOutFile.Name(),
				MappingFile: tempCSVFile.Name(), // Using CSV file as mapping file will cause JSON parsing error
			},
			wantErr:   true,
			errString: "failed to parse mapping file",
		},
		{
			name: "non-existent input file",
			mapper: &Mapper{
				In:          "non_existent_file.csv",
				Out:         tempOutFile.Name(),
				MappingFile: tempMappingFile.Name(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, writer, err := tt.mapper.initialize()

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("initialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check error message if expected
			if tt.wantErr && tt.errString != "" && (err == nil || !strings.Contains(err.Error(), tt.errString)) {
				t.Errorf("initialize() error = %v, should contain %v", err, tt.errString)
				return
			}

			// Check if reader and writer are non-nil when no error
			if !tt.wantErr {
				if reader == nil {
					t.Errorf("initialize() reader is nil")
				} else {
					reader.Close()
				}

				if writer == nil {
					t.Errorf("initialize() writer is nil")
				} else {
					writer.Close()
				}
			}
		})
	}
}

// TestGetApplicationValue tests the getApplicationValue function which computes values based on application data.
func TestGetApplicationValue(t *testing.T) {
	mapper := &Mapper{}

	tests := []struct {
		name      string
		field     CalculatedField
		index     int
		want      any
		wantErr   bool
		errString string
	}{
		{
			name: "record format returns index",
			field: CalculatedField{
				Format: "record",
				Type:   "int",
			},
			index:   42,
			want:    42,
			wantErr: false,
		},
		{
			name: "unknown format returns error",
			field: CalculatedField{
				Format: "unknown",
				Type:   "int",
			},
			index:     0,
			want:      nil,
			wantErr:   true,
			errString: "unknown format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.getApplicationValue(tt.field, tt.index)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("getApplicationValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check error message if expected
			if tt.wantErr && tt.errString != "" && (err == nil || !strings.Contains(err.Error(), tt.errString)) {
				t.Errorf("getApplicationValue() error = %v, should contain %v", err, tt.errString)
				return
			}

			// Check result
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getApplicationValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetDateTimeValue tests the getDateTimeValue function which generates formatted date and time values.
func TestGetDateTimeValue(t *testing.T) {
	mapper := &Mapper{}

	tests := []struct {
		name    string
		field   CalculatedField
		wantErr bool
	}{
		{
			name: "date format",
			field: CalculatedField{
				Format: "2006-01-02",
				Type:   "string",
			},
			wantErr: false,
		},
		{
			name: "time format",
			field: CalculatedField{
				Format: "15:04:05",
				Type:   "string",
			},
			wantErr: false,
		},
		{
			name: "datetime format",
			field: CalculatedField{
				Format: "2006-01-02 15:04:05",
				Type:   "string",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapper.getDateTimeValue(tt.field)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("getDateTimeValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check that result is a string
			if !tt.wantErr {
				if _, ok := got.(string); !ok {
					t.Errorf("getDateTimeValue() result is not a string: %v", got)
					return
				}

				// Verify that the result can be parsed using the same format
				_, err := time.Parse(tt.field.Format, got.(string))
				if err != nil {
					t.Errorf("getDateTimeValue() result cannot be parsed with the same format: %v", err)
				}
			}
		})
	}
}

// TestMapWithCalculatedFields tests the Map method with various calculated fields.
func TestMapWithCalculatedFields(t *testing.T) {
	// Create a temporary mapping file with calculated fields
	mappingJSON := `{
		"mapping": {
			"id": {
				"property": "property1",
				"type": "int"
			},
			"text": {
				"property": "property2",
				"type": "string"
			}
		},
		"calculated": [
			{
				"property": "calculated.record",
				"kind": "application",
				"format": "record",
				"type": "int"
			},
			{
				"property": "calculated.date",
				"kind": "datetime",
				"format": "2006-01-02",
				"type": "string"
			},
			{
				"property": "calculated.extra",
				"kind": "extra",
				"format": "test-var",
				"type": "string"
			}
		],
		"extra_variables": {
			"test-var": {
				"value": "test value"
			}
		}
	}`

	tempMappingFile, err := os.CreateTemp("", "mapping_calc*.json")
	if err != nil {
		t.Fatalf("Failed to create temp mapping file: %v", err)
	}
	defer os.Remove(tempMappingFile.Name())

	if _, err := tempMappingFile.Write([]byte(mappingJSON)); err != nil {
		t.Fatalf("Failed to write to temp mapping file: %v", err)
	}
	tempMappingFile.Close()

	// Create a temporary CSV file
	csvData := `id,text
1,hello
2,world`

	tempCSVFile, err := os.CreateTemp("", "input_calc*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(tempCSVFile.Name())

	if _, err := tempCSVFile.Write([]byte(csvData)); err != nil {
		t.Fatalf("Failed to write to temp CSV file: %v", err)
	}
	tempCSVFile.Close()

	// Create a temporary output file
	tempOutFile, err := os.CreateTemp("", "output_calc*.json")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer os.Remove(tempOutFile.Name())
	tempOutFile.Close()

	// Create and run mapper
	mapper, err := NewMapper(
		WithIn(tempCSVFile.Name()),
		WithOut(tempOutFile.Name()),
		WithNamed(true),
		WithArray(true),
		WithMappingFile(tempMappingFile.Name()),
		WithOutputType("json"),
	)
	if err != nil {
		t.Fatalf("Failed to create mapper: %v", err)
	}

	err = mapper.Map()
	if err != nil {
		t.Fatalf("Failed to map: %v", err)
	}

	// Read the output file
	outputData, err := os.ReadFile(tempOutFile.Name())
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Parse the JSON output
	var result []map[string]interface{}
	if err := json.Unmarshal(outputData, &result); err != nil {
		t.Fatalf("Failed to parse output JSON: %v", err)
	}

	// Verify the calculated fields
	if len(result) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(result))
	}

	// Check record numbers
	record0, ok := result[0]["calculated"].(map[string]interface{})["record"]
	if !ok || record0 != float64(0) { // JSON numbers are float64
		t.Errorf("Expected record 0, got %v", record0)
	}

	record1, ok := result[1]["calculated"].(map[string]interface{})["record"]
	if !ok || record1 != float64(1) { // JSON numbers are float64
		t.Errorf("Expected record 1, got %v", record1)
	}

	// Check date fields (just verify they exist, as the actual date will vary)
	date0, ok := result[0]["calculated"].(map[string]interface{})["date"]
	if !ok || date0 == "" {
		t.Errorf("Expected date field in record 0, got %v", date0)
	}

	// Check extra variable
	extra0, ok := result[0]["calculated"].(map[string]interface{})["extra"]
	if !ok || extra0 != "test value" {
		t.Errorf("Expected extra field 'test value', got %v", extra0)
	}
}

// TestMapWithEnvironmentVariables tests the Map method with environment variables.
func TestMapWithEnvironmentVariables(t *testing.T) {
	// Set a test environment variable
	testEnvVar := "TEST_ENV_VAR_FOR_CSV2JSON"
	testEnvValue := "test environment value"
	os.Setenv(testEnvVar, testEnvValue)
	defer os.Unsetenv(testEnvVar)

	// Create a temporary mapping file with environment variable
	mappingJSON := `{
		"mapping": {
			"id": {
				"property": "property1",
				"type": "int"
			}
		},
		"calculated": [
			{
				"property": "env.test",
				"kind": "environment",
				"format": "TEST_ENV_VAR_FOR_CSV2JSON",
				"type": "string"
			}
		]
	}`

	tempMappingFile, err := os.CreateTemp("", "mapping_env*.json")
	if err != nil {
		t.Fatalf("Failed to create temp mapping file: %v", err)
	}
	defer os.Remove(tempMappingFile.Name())

	if _, err := tempMappingFile.Write([]byte(mappingJSON)); err != nil {
		t.Fatalf("Failed to write to temp mapping file: %v", err)
	}
	tempMappingFile.Close()

	// Create a temporary CSV file
	csvData := `id
1`

	tempCSVFile, err := os.CreateTemp("", "input_env*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp CSV file: %v", err)
	}
	defer os.Remove(tempCSVFile.Name())

	if _, err := tempCSVFile.Write([]byte(csvData)); err != nil {
		t.Fatalf("Failed to write to temp CSV file: %v", err)
	}
	tempCSVFile.Close()

	// Create a temporary output file
	tempOutFile, err := os.CreateTemp("", "output_env*.json")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer os.Remove(tempOutFile.Name())
	tempOutFile.Close()

	// Create and run mapper
	mapper, err := NewMapper(
		WithIn(tempCSVFile.Name()),
		WithOut(tempOutFile.Name()),
		WithNamed(true),
		WithArray(true),
		WithMappingFile(tempMappingFile.Name()),
		WithOutputType("json"),
	)
	if err != nil {
		t.Fatalf("Failed to create mapper: %v", err)
	}

	err = mapper.Map()
	if err != nil {
		t.Fatalf("Failed to map: %v", err)
	}

	// Read the output file
	outputData, err := os.ReadFile(tempOutFile.Name())
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Parse the JSON output
	var result []map[string]interface{}
	if err := json.Unmarshal(outputData, &result); err != nil {
		t.Fatalf("Failed to parse output JSON: %v", err)
	}

	// Verify the environment variable
	if len(result) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(result))
	}

	envValue, ok := result[0]["env"].(map[string]interface{})["test"]
	if !ok || envValue != testEnvValue {
		t.Errorf("Expected environment value '%s', got %v", testEnvValue, envValue)
	}
}

// compareMappers compares two Mapper instances for equality, considering their public fields: In, Out, Array, and Named.
// TODO switch to cmp
func compareMappers(a, b *Mapper) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.In == b.In &&
		a.Out == b.Out &&
		a.Array == b.Array &&
		a.Named == b.Named
}
