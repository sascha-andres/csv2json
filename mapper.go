package csv2json

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// OptionFunc defines a function signature for configuring a Mapper instance with specific options or parameters.
type OptionFunc func(*Mapper) error

// WithIn sets the "In" field of the Mapper instance to the provided string if it's not empty, otherwise returns an error.
func WithIn(in string) OptionFunc {
	return func(mapper *Mapper) error {
		if strings.TrimSpace(in) == "" {
			return errors.New("-in may not be empty")
		}
		mapper.In = in
		return nil
	}
}

// WithOut sets the "Out" field of the Mapper instance to the provided string if it's not empty, otherwise returns an error.'
func WithOut(out string) OptionFunc {
	return func(mapper *Mapper) error {
		if strings.TrimSpace(out) == "" {
			return errors.New("-out may not be empty")
		}
		mapper.Out = out
		return nil
	}
}

// WithArray sets the "Array" field of the Mapper instance to the provided boolean value.
func WithArray(array bool) OptionFunc {
	return func(mapper *Mapper) error {
		mapper.Array = array
		return nil
	}
}

// WithNamed sets the "Named" field of the Mapper instance to the provided boolean value.
func WithNamed(named bool) OptionFunc {
	return func(mapper *Mapper) error {
		mapper.Named = named
		return nil
	}
}

// WithMappingFile sets the "MappingFile" field of the Mapper instance to the provided string.
func WithMappingFile(mappingFile string) OptionFunc {
	return func(mapper *Mapper) error {
		mapper.MappingFile = mappingFile
		return nil
	}
}

// WithOutputType sets the specified output type for marshaling data in a Mapper instance.
func WithOutputType(outputType string) OptionFunc {
	return func(mapper *Mapper) error {
		mapper.MarshalWith = outputType
		switch outputType {
		case "json":
		case "yaml":
		case "toml":
			break
		case "":
			mapper.MarshalWith = "json"
			break
		default:
			return errors.New(fmt.Sprintf("unknown marshaling type %q", outputType))
		}
		return nil
	}
}

// WithTomlPropertyName sets the property name for TOML array output.
func WithTomlPropertyName(propertyName string) OptionFunc {
	return func(mapper *Mapper) error {
		if propertyName == "" {
			mapper.TomlPropertyName = "data"
		} else {
			mapper.TomlPropertyName = propertyName
		}
		return nil
	}
}

// NewMapper creates and initializes a new Mapper instance using the provided OptionFunc configurations.
func NewMapper(options ...OptionFunc) (*Mapper, error) {
	mapper := &Mapper{}
	for _, option := range options {
		if err := option(mapper); err != nil {
			return nil, err
		}
	}
	switch mapper.MarshalWith {
	case "json":
		mapper.marshaler = json.Marshal
		break
	case "yaml":
		mapper.Array = true
		mapper.marshaler = yaml.Marshal
		break
	case "toml":
		mapper.Array = true
		mapper.marshaler = toml.Marshal
		break
	}
	return mapper, nil
}

// Map processes input CSV data, maps it to JSON according to the configuration, and writes the result to the output destination.
func (m *Mapper) Map() error {
	reader, writer, err := m.initialize()
	if err != nil {
		return err
	}
	defer reader.Close()
	defer writer.Close()

	csvIn := csv.NewReader(reader)
	csvIn.ReuseRecord = false

	var (
		arrResult []map[string]any
		header    []string
	)

	// Read header if needed
	if m.Named {
		header, err = csvIn.Read()
		if err != nil {
			return err
		}
	}
	// from now on we can reuse the record
	csvIn.ReuseRecord = true
	if m.Array {
		arrResult = make([]map[string]any, 0)
	}
	recordNumber := 0
	// Read all records
	for {
		record, err := csvIn.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		out := make(map[string]interface{})
		for i := range record {
			key := fmt.Sprintf("%d", i)
			if m.Named {
				key = header[i]
			}
			var (
				v  ColumnConfiguration
				ok bool
			)
			if v, ok = m.configuration.Mapping[key]; !ok {
				return errors.New("mapping configuration missing for key " + key)
			}

			out = setValue(strings.Split(v.Property, "."), convertToType(v.Type, record[i]), out)
		}
		// calculated fields
		for _, field := range m.configuration.Calculated {
			var val any
			switch field.Kind {
			case "application":
				val, err = m.getApplicationValue(field, recordNumber)
				if err != nil {
					return err
				}
				break
			case "datetime":
				val, err = m.getDateTimeValue(field)
				if err != nil {
					return err
				}
				break
			case "environment":
				e := os.Getenv(field.Format)
				val = convertToType(field.Type, e)
				break
			case "extra":
				e, ok := m.configuration.ExtraVariables[field.Format]
				if !ok {
					return errors.New("extra variable " + field.Format + " not found")
				}
				val = convertToType(field.Type, e.Value)
			default:
				return errors.New("unknown kind " + field.Kind)
			}
			out = setValue(strings.Split(field.Property, "."), val, out)
		}
		if m.Array {
			arrResult = append(arrResult, out)
		} else {
			d, err := m.marshaler(out)
			if err != nil {
				return err
			}
			if recordNumber > 0 {
				_, _ = writer.Write([]byte("\n"))
			}
			_, err = writer.Write(d)
			if err != nil {
				return err
			}
		}
		recordNumber++
	}
	if m.Array {
		var d []byte
		if m.MarshalWith == "toml" {
			// Set default property name if not specified
			propertyName := "data"
			if m.TomlPropertyName != "" {
				propertyName = m.TomlPropertyName
			}

			// Create a map with the custom property name as the key
			tomlData := map[string]any{
				propertyName: arrResult,
			}

			d, err = m.marshaler(tomlData)
			if err != nil {
				return err
			}
		} else {
			d, err = m.marshaler(arrResult)
			if err != nil {
				return err
			}
		}
		_, _ = writer.Write(d)
	}
	return nil
}

// getDateTimeValue generates a date and time value formatted based on the Format field of the CalculatedField structure.
func (m *Mapper) getDateTimeValue(field CalculatedField) (any, error) {
	return time.Now().Format(field.Format), nil
}

// getApplicationValue computes and returns a value based on the specified CalculatedField and index.
// Returns the computed value or an error if the field format is unknown.
func (m *Mapper) getApplicationValue(field CalculatedField, i int) (any, error) {
	switch field.Format {
	case "record":
		return convertToType("int", strconv.Itoa(i)), nil
	}
	return nil, errors.New("unknown format " + field.Format)
}

// convertToType converts the input string `val` to a specified type `t` such as "int", "float", or "bool".
// Returns the converted value as `any` or an error if the conversion fails.
func convertToType(t, val string) any {
	switch t {
	case "int":
		i, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		return i
	case "float":
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		return f
	case "bool":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		return b
	}
	return val
}

// setValue creates and maps nested dictionaries based on a hierarchy of keys, assigning a final value.
func setValue(hierarchy []string, value any, data map[string]interface{}) map[string]interface{} {
	v := setValueInternal(hierarchy, value, data)
	data[hierarchy[0]] = v
	return data
}

// setValueInternal recursively creates and maps nested dictionaries based on a hierarchy of keys, assigning a final value.
func setValueInternal(hierarchy []string, value any, inside map[string]any) any {
	if len(hierarchy) == 1 {
		return value
	}
	v := make(map[string]any)
	if val, ok := inside[hierarchy[0]]; ok {
		if reflected, ok := val.(map[string]any); ok {
			v = reflected
		}
	}
	v[hierarchy[1]] = setValueInternal(hierarchy[1:], value, v)
	return v
}

// initialize initializes the Mapper instance by reading the mapping file and opening the input and output files.
func (m *Mapper) initialize() (io.ReadCloser, io.WriteCloser, error) {
	mappingFile := "mapping.json"
	if m.MappingFile != "" {
		mappingFile = m.MappingFile
	}

	configData, err := os.ReadFile(mappingFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read mapping file: %w", err)
	}

	if err := json.Unmarshal(configData, &m.configuration); err != nil {
		return nil, nil, fmt.Errorf("failed to parse mapping file: %w", err)
	}
	var (
		fIn  io.ReadCloser
		fOut io.WriteCloser
	)
	if m.In == "-" {
		fIn = os.Stdin
	} else {
		fIn, err = os.OpenFile(m.In, os.O_RDONLY, 0600)
		if err != nil {
			return nil, nil, err
		}
	}
	if m.Out == "-" {
		fOut = os.Stdout
	} else {
		fOut, err = os.OpenFile(m.Out, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return nil, nil, err
		}
	}
	return fIn, fOut, nil
}
