package helpers

import (
	"math/rand"
	"strconv"
	"time"
)

func RandomIntToString(min, max int) string {
	rand.Seed(time.Now().UnixNano())

	id := rand.Intn(max-min+1) + min

	return strconv.Itoa(int(id))
}

func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(max-min+1) + min
}
