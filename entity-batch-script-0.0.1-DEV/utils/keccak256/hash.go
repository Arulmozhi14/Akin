package keccak256

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

// Hash asks for a string and return the hash
func Hash(data string) string {
	var bytes []byte
	keccak256 := sha3.NewLegacyKeccak256()

	keccak256.Write([]byte(data))
	bytes = keccak256.Sum(nil)

	return hex.EncodeToString(bytes)
}
