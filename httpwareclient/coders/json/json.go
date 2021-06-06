package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

// Coder is a json coder.
type Coder struct{}

// Encode encode to json.
func (Coder) Encode(data interface{}) (io.Reader, error) {
	var bodyReader io.Reader
	if data != nil {
		buffer := new(bytes.Buffer)
		if err := json.NewEncoder(buffer).Encode(data); err != nil {
			return nil, err
		}
		bodyReader = buffer
	}
	return bodyReader, nil
}

// Decode from json.
func (Coder) Decode(reader io.Reader, data interface{}) error {
	respData, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read body from response, %s", err)
	}

	if err := json.Unmarshal(respData, data); err != nil {
		return fmt.Errorf("failed to unmarshal body, %s", err)
	}

	return nil
}
