package uuid

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

type ID [16]byte

func NewID() ID {
	var n ID

	uuidv4 := uuid.New()
	j := bytes.Join(
		[][]byte{
			uuidv4[6:8],
			uuidv4[4:6],
			uuidv4[0:4],
			uuidv4[8:10],
			uuidv4[10:],
		},
		[]byte{},
	)
	copy(n[:], j)

	return n
}

func NewFromString(str string) (ID, error) {
	var n ID

	b, err := hex.DecodeString(str)
	if err != nil {
		return ID([16]byte{}), err
	}

	j := bytes.Join([][]byte{b[6:8], b[4:6], b[0:4], b[8:10], b[10:]}, []byte{})
	copy(n[:], j)

	return n, nil
}

// String encode NID to string.
func (n ID) String() string {
	return hex.EncodeToString(n[:])
}

// DecodeString decode string to NID.
func DecodeString(str string) (ID, error) {
	var n ID

	b, err := hex.DecodeString(str)
	if err != nil {
		return ID([16]byte{}), err
	}

	copy(n[:], b)

	return n, nil
}

// Value encode unique identification for store to db.
func (n ID) Value() (driver.Value, error) {
	return n[:], nil
}

// Scan decode unique identification to get from db.
func (n *ID) Scan(src interface{}) error {
	var nid ID

	s, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("error scan uid %w", src)
	}

	copy(nid[:], s)
	*n = nid

	return nil
}

func DecodeStringWithoutErr(str string) ID {
	var n ID

	b, err := hex.DecodeString(str)
	if err != nil {
		return ID([16]byte{})
	}

	copy(n[:], b)

	return n
}
