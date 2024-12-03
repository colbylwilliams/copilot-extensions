package jsonschema

import "encoding/json"

// DataType is a string that specifies the data type of a JSON Schema.
type DataType string

const (
	// Object is the data type for an object.
	Object DataType = "object"

	// Number is the data type for a number.
	Number DataType = "number"

	// Integer is the data type for an integer.
	Integer DataType = "integer"

	// String is the data type for a string.
	String DataType = "string"

	// Array is the data type for an array.
	Array DataType = "array"

	// Null is the data type for null.
	Null DataType = "null"

	// Boolean is the data type for a boolean.
	Boolean DataType = "boolean"
)

// Definition is a struct for describing a JSON Schema.
// It is fairly limited, and you may have better luck using a third-party library.
type Definition struct {
	// Type specifies the data type of the schema.
	Type DataType `json:"type,omitempty"`

	// Description is the description of the schema.
	Description string `json:"description,omitempty"`

	// Enum is used to restrict a value to a fixed set of values. It must be an array with at least
	// one element, where each element is unique. You will probably only use this with strings.
	Enum []string `json:"enum,omitempty"`

	// Properties describes the properties of an object, if the schema type is Object.
	Properties map[string]Definition `json:"properties"`

	// Required specifies which properties are required, if the schema type is Object.
	Required []string `json:"required,omitempty"`

	// Items specifies which data type an array contains, if the schema type is Array.
	Items *Definition `json:"items,omitempty"`
}

// MarshalJSON marshals a Definition to JSON.
func (d Definition) MarshalJSON() ([]byte, error) {
	if d.Properties == nil {
		d.Properties = make(map[string]Definition) //nolint:revive // not meant to be visible externally
	}
	type Alias Definition
	return json.Marshal(struct {
		Alias
	}{
		Alias: Alias(d),
	})
}
