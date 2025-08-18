package json

import (
	"encoding/json"
	"io"
)

// Decode loads source to destination using json deserialization.
func Decode(destination any, source io.Reader) error {
	decoder := json.NewDecoder(source)
	err := decoder.Decode(destination)
	if err != nil {
		return err
	}

	return nil
}
