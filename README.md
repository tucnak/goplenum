# Enum [![GoDoc](https://godoc.org/github.com/tucnak/enum?status.svg)](https://godoc.org/github.com/tucnak/enum) [![Go Report Card](https://goreportcard.com/badge/github.com/tucnak/enum)](https://goreportcard.com/report/github.com/tucnak/enum) [![cover.run go](https://cover.run/go/github.com/tucnak/enum.svg?tag=golang-1.10)](https://cover.run/go/github.com/tucnak/enum?tag=golang-1.10)
Enum is a tool to generate Go code that adds useful methods to Go enums (constants with a specific type).
It's a fork of [alvaroloes/enumer](https://github.com/alvaroloes/enumer), which originally started as a fork of Rob Pike’s [Stringer](https://godoc.org/golang.org/x/tools/cmd/stringer).

## Generated functions and methods
When `enum` is applied to a type, it will generate:

* The following basic methods/functions: 

  * Method `String()`: returns the string representation of the enum value. This makes the enum conform
the `Stringer` interface, so whenever you print an enum value, you'll get the string name instead of a number.
  * `MarshalJSON()` and `UnmarshalJSON()`. These make the enum conform to the `json.Marshaler` and `json.Unmarshaler` interfaces.
    Very useful to use it in JSON APIs.
  * When the flag `text` is provided, two additional methods will be generated, `MarshalText()` and `UnmarshalText()`. These make
    the enum conform to the `encoding.TextMarshaler` and `encoding.TextUnmarshaler` interfaces. 
    **Note:** If you use your enum values as keys in a map and you encode the map as _JSON_, you need this flag set to true to properly
    convert the map keys to json (strings). If not, the numeric values will be used instead.
  * When the flag `yaml` is provided, two additional methods will be generated, `MarshalYAML()` and `UnmarshalYAML()`. These make
    the enum conform to the `gopkg.in/yaml.v2.Marshaler` and `gopkg.in/yaml.v2.Unmarshaler` interfaces.
  * Method `Scan()` and `Value()` interfaces required for storing the enum in a database.
  * When the flag `nojson` is provided, two additional JSON-associated methods will **not** be generated. 
  * When the flag `nosql` is provided, the methods for implementing the Scanner and Valuer interfaces will **not** be generated.
  * Function `enum<Type>Of(s string)`: returns the enum value from its string representation. This is useful
    when you need to read enum values from command line arguments, from a configuration file, or
    from a REST API request... In short, from those places where using the real enum value (an integer) would
    be almost meaningless or hard to trace or use by a human.

For example, if we have an enum type called `Pill`,
```go
type Pill int

const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen
	Paracetamol
	Acetaminophen = Paracetamol
)
```
executing `enum Pill` will generate a new file with the following methods:
```go
func (i Pill) String() string { 
	//...
}

func (i Pill) MarshalJSON() ([]byte, error) {
	//...
}

func (i *Pill) UnmarshalJSON(data []byte) error {
	//...
}

func (i Pill) Value() (driver.Value, error) {
	//...
}

func (i *Pill) Scan(value interface{}) error {
	//...
}
```

The generated code is exactly the same as the Stringer tool plus the mentioned additions, so you can use
**Enum** where you are already using **Stringer** without any code change.

## Transforming the string representation of the enum value

By default, Enum uses the same name of the enum value for generating the string representation (usually CamelCase in Go).

```go
type MyType int

 ...

name := MyTypeValue.String() // name => "MyTypeValue"
```

Sometimes you need to use some other string representation format than CamelCase (i.e. in JSON).

To transform it from CamelCase to snake_case or kebab-case, you can use the `transform` flag.

For example, the command `enum -type=MyType -json -transform=snake` would generate the following string representation:

```go
name := MyTypeValue.String() // name => "my_type_value"
```
**Note**: The transformation only works form CamelCase to snake_case or kebab-case, not the other way around.

## How to use
The usage of Enum is the same as Stringer, so you can refer to the [Stringer docs](https://godoc.org/golang.org/x/tools/cmd/stringer)
for more information.

There are four boolean flags: `text`, `yaml`, `nojson` and `nosql` (as JSON and SQL generators are implicit).

You can use any combination of them (i.e. `enum Pill -nojson -text`).

For enum string representation transformation the `transform` and `trimprefix` flags
were added (i.e. `enum -transform=snake`).
Possible transform values are `snake` and `kebab` for transformation to snake_case and kebab-case accordingly.
The default value for `transform` flag is `noop` which means no transformation will be performed.

If a prefix is provided via the `trimprefix` flag, it will be trimmed from the start of each name (before
it is transformed). If a name doesn't have the prefix it will be passed unchanged.

## Inspiring projects
* [Stringer](https://godoc.org/golang.org/x/tools/cmd/stringer)
* [jsonenums](https://github.com/campoy/jsonenums)
* [alvaroloes/enumer](https://github.com/alvaroloes/enumer)
