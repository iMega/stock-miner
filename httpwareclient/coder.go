package httpwareclient

import (
	"io"

	"github.com/imega/stock-miner/httpwareclient/coders/json"
)

const (
	// JSON type coder.
	JSON CoderType = iota
	// XML type coder.
	XML
)

type (
	// CoderType type coder.
	CoderType int

	// Coder is a interface coding.
	Coder interface {
		Encode(data interface{}) (io.Reader, error)
		Decode(reader io.Reader, data interface{}) error
	}
)

// GetCoder return coder by type.
func GetCoder(coderType CoderType) Coder {
	switch coderType {
	case JSON:
		return &json.Coder{}
	case XML:
		return &nullCoder{}
	default:
		return &nullCoder{}
	}
}

type nullCoder struct{}

func (nullCoder) Encode(data interface{}) (io.Reader, error) {
	return nil, nil
}

func (nullCoder) Decode(reader io.Reader, data interface{}) error {
	return nil
}
