package uuid_test

import (
	"database/sql/driver"
	"testing"

	"github.com/imega/stock-miner/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	r := uuid.NewID()
	r1 := uuid.NewID()

	assert.NotEqual(t, r, r1)
}

func TestUID_Scan(t *testing.T) {
	r := []byte{0x43, 0xa1, 0xf3, 0x83, 0x8c, 0x37, 0x71, 0x84, 0xa3, 0xf5, 0x56, 0x2a, 0x99, 0x9c, 0x99, 0x7b}
	exp, _ := uuid.DecodeString("43a1f3838c377184a3f5562a999c997b")

	var n uuid.ID
	err := n.Scan(r)

	assert.Nil(t, err)
	assert.Equal(t, n, exp)
}

func TestUID_Value(t *testing.T) {
	u, _ := uuid.DecodeString("43a1f3838c377184a3f5562a999c997b")
	exp := []byte{0x43, 0xa1, 0xf3, 0x83, 0x8c, 0x37, 0x71, 0x84, 0xa3, 0xf5, 0x56, 0x2a, 0x99, 0x9c, 0x99, 0x7b}
	dr, err := u.Value()

	assert.Nil(t, err)
	assert.Equal(t, dr, driver.Value(exp))
}

func TestUID_ScanError(t *testing.T) {
	r := 123
	var u uuid.ID
	err := u.Scan(r)

	assert.Error(t, err)
	assert.Equal(t, uuid.ID([16]byte{}), u)
}

func TestUID_ValueError(t *testing.T) {
	_, err := uuid.DecodeString("123")
	assert.Error(t, err)
}

func BenchmarkNewID(b *testing.B) {
	b.SetBytes(2)
	for i := 0; i < b.N; i++ {
		uuid.NewID()
	}
}

func TestNewFromByte(t *testing.T) {
	expected := uuid.ID(uuid.ID{0x50, 0x0, 0x0, 0xf, 0x25, 0x85, 0xcb, 0x29, 0xa0, 0x0, 0x1d, 0x86, 0xcd, 0xa1, 0x84, 0x7c})
	actual, err := uuid.NewFromString("2585cb29000f5000a0001d86cda1847c")
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
