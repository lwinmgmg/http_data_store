package helper

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/google/uuid"
)

func HexString(input string) string {
	shaData := sha512.Sum512([]byte(input))
	return hex.EncodeToString(shaData[:])
}

func GetUniqueKey() string {
	myUUID := uuid.New()
	return myUUID.String()
}
