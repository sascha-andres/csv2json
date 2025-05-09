# csv2json

Reads a CSV file and converts it to a JSON file. To do it, it takes a mapping from a mapping file that looks like this:

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

Instead of using a 0 based index index to adress csv columns you might use names when you pass the `-named` flag.

To provide the mapping the tool will look for a `mapping.json` file in the current directory unless `-mapping` is used to provide a path to one.

By default it uses stdin to read the json unless `-in` is used to provide a file. Accordingly it prints to stdout unless `-out` is provided.

Given the following CSV:

```csv
1,"hello",2.3
```

The output will look like this (pretty printed):

```json
{
  "property1": 1,
  "property2": {
    "property3": "Hello"
  },
  "property4": 2.3
}
```

For multiple rows in the csv it will print one JSON document. Except when given the `-array` flag, then the result will look like this:

```json
[
  {
    "property1": 1,
    "property2": {
      "property3": "Hello"
    },
    "property4": 2.3
  }
]
```