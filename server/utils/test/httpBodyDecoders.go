package test

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

// EncodeJSON encodes json data in bytes
func EncodeJSON(data interface{}) ([]byte, error) {
	if data == nil {
		return []byte{}, nil
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(data)
	if err != nil {
		return []byte{}, errors.Wrap(err, "Error while encoding json")
	}

	return buf.Bytes(), nil

}
