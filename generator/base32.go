package generator

import (
	"encoding/base32"
	"encoding/binary"
	"time"
)

// RandomBase32String ...
func RandomBase32String() string {
	now := time.Now().Unix()
	var x [16]byte
	value := binary.PutVarint(x[:], now)

	codex := base32.StdEncoding.EncodeToString(x[:value])
	return codex
}
