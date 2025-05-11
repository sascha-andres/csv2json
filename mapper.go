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
				v   ColumnConfiguration
				val any
				ok  bool
			)
			if v, ok = m.configuration.Mapping[key]; !ok {
				return errors.New("mapping configuration missing for key " + key)
			}

			switch v.Type {
			case "int":
				i, err := strconv.Atoi(record[i])
				if err != nil {
					return err
				}
				val = i
				break
			case "float":
				f, err := strconv.ParseFloat(record[i], 64)
				if err != nil {
					return err
				}
				val = f
				break
			case "bool":
				b, err := strconv.ParseBool(record[i])
				if err != nil {
					return err
				}
				val = b
				break
			default:
				val = record[i]
				break
			}
			out = setValue(strings.Split(v.Property, "."), val, out)
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
