package main

import (
	"fmt"
)

// Arguments to format are:
//	[1]: type name
const stringNameToValueMethod = `// enum%[1]s retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func enum%[1]sOf(s string) (%[1]s, error) {
	if val, ok := enum%[1]sMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%%s does not belong to %[1]s values", s)
}
`

func (g *Generator) buildBasicExtras(runs [][]Value, typeName string, runsThreshold int) {
	// At this moment, either "g.declareIndexAndNameVars()" or "g.declareNameVars()" has been called

	// Print the slice of values
	g.Printf("\nvar enum%sValues = []%s{", typeName, typeName)
	for _, values := range runs {
		for _, value := range values {
			g.Printf("\t%s, ", value.str)
		}
	}
	g.Printf("}\n\n")

	// Print the map between name and value
	g.Printf("\nvar enum%sMap = map[string]%s{\n", typeName, typeName)
	thereAreRuns := len(runs) > 1 && len(runs) <= runsThreshold
	var n int
	var runID string
	for i, values := range runs {
		if thereAreRuns {
			runID = "_" + fmt.Sprintf("%d", i)
			n = 0
		} else {
			runID = ""
		}

		for _, value := range values {
			g.Printf("\tenum%sRepr%s[%d:%d]: %s,\n", typeName, runID, n, n+len(value.name), &value)
			n += len(value.name)
		}
	}
	g.Printf("}\n\n")

	// Print the basic extra methods
	g.Printf(stringNameToValueMethod, typeName)
}

// Arguments to format are:
//	[1]: type name
const jsonMethods = `
// MarshalJSON implements the json.Marshaler interface for %[1]s
func (i %[1]s) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for %[1]s
func (i *%[1]s) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("%[1]s should be a string, got %%s", data)
	}

	var err error
	*i, err = enum%[1]sOf(s)
	return err
}
`

func (g *Generator) buildJSONMethods(runs [][]Value, typeName string, runsThreshold int) {
	g.Printf(jsonMethods, typeName)
}

// Arguments to format are:
//	[1]: type name
const textMethods = `
// MarshalText implements the encoding.TextMarshaler interface for %[1]s
func (i %[1]s) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for %[1]s
func (i *%[1]s) UnmarshalText(text []byte) error {
	var err error
	*i, err = enum%[1]sOf(string(text))
	return err
}
`

func (g *Generator) buildTextMethods(runs [][]Value, typeName string, runsThreshold int) {
	g.Printf(textMethods, typeName)
}

// Arguments to format are:
//	[1]: type name
const yamlMethods = `
// MarshalYAML implements a YAML Marshaler for %[1]s
func (i %[1]s) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for %[1]s
func (i *%[1]s) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = enum%[1]sOf(s)
	return err
}
`

func (g *Generator) buildYAMLMethods(runs [][]Value, typeName string, runsThreshold int) {
	g.Printf(yamlMethods, typeName)
}
