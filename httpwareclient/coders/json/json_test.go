package json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestJsonEncoder(t *testing.T) {
	var coder Coder

	type Data struct {
		Field string `json:"filed"`
	}

	data := &Data{
		Field: "test",
	}

	expected, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}

	reader, err := coder.Encode(data)
	if err != nil {
		t.Error(err)
	}

	actual, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Error(err)
	}

	if reflect.DeepEqual(expected, actual) {
		t.Error("data not equal")
	}
}

func TestJsonDecoder(t *testing.T) {
	var coder Coder

	type Data struct {
		Field string `json:"filed"`
	}

	expected := &Data{
		Field: "test",
	}

	b, err := json.Marshal(expected)
	if err != nil {
		t.Error(err)
	}

	actual := Data{}

	if err := coder.Decode(bytes.NewBuffer(b), &actual); err != nil {
		t.Error(err)
	}

	if reflect.DeepEqual(expected, actual) {
		t.Error("data not equal")
	}
}
