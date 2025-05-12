package main

import (
	"github.com/sascha-andres/reuse/flag"

	"github.com/sascha-andres/csv2json"
)

var (
	in                 string
	out                string
	array              bool
	named              bool
	mappingFile        string
	outputType         string
	nestedPropertyName string
)

// init initializes the command-line flags and environment variables.
func init() {
	flag.SetEnvPrefix("CSV2JSON")
	flag.StringVar(&in, "in", "-", "input file, defaults to stdin")
	flag.StringVar(&out, "out", "-", "output file, defaults to stdout")
	flag.BoolVar(&array, "array", false, "output as array (implicit for yaml and toml)")
	flag.BoolVar(&named, "named", false, "output as named")
	flag.StringVar(&mappingFile, "mapping", "mapping.json", "mapping file")
	flag.StringVar(&outputType, "output-type", "json", "output type, one of json, yaml or toml")
	flag.StringVar(&nestedPropertyName, "nested-property", "data", "property name for nested array output")
}

// main parses flags, executes the application logic via the run function, and handles any errors by panicking.
func main() {
	flag.Parse()

	if err := run(); err != nil {
		panic(err)
	}
}

// run initializes a Mapper instance with provided configurations and processes CSV input into JSON format.
// It returns an error if any step in the process fails.
func run() error {
	m, err := csv2json.NewMapper(
		csv2json.WithOutputType(outputType),
		csv2json.WithOut(out),
		csv2json.WithArray(array),
		csv2json.WithIn(in),
		csv2json.WithMappingFile(mappingFile),
		csv2json.WithNamed(named),
		csv2json.WithNestedPropertyName(nestedPropertyName))
	if err != nil {
		return err
	}
	return m.Map()
}
