package serializer

import (
	"encoding/json"
	"io"
)

// ToJSON serializes the given interface into a string based JSON format.
func ToJSON(i interface{}, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(i)
}

// FromJSON deserializes the object from JSON string in an io.Reader to the given interface.
func FromJSON(i interface{}, reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(i)
}
