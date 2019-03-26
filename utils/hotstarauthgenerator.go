package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

//Generate generates Hotstar authorization code for given time.
func Generate(start int64) string {

	expiry := start + 6000

	message := fmt.Sprintf("st=%d~exp=%d~acl=/*", start, expiry)

	akamaiEncryptionKey := "05fc1a01cac94bc412fc53120775f9ee"

	hexedAkamaiEncryptionKey, err := hex.DecodeString(akamaiEncryptionKey)
	if err != nil {
		panic(err)
	}

	hmacInstance := hmac.New(sha256.New, hexedAkamaiEncryptionKey)
	hmacInstance.Write([]byte(message))
	hmacedMessage := hex.EncodeToString(hmacInstance.Sum(nil))

	return fmt.Sprintf("%s~hmac=%s", message, hmacedMessage)
}

//GenerateHotstarAuth generates Hotstar authorization code.
func GenerateHotstarAuth() string {
	start := time.Now().UnixNano() / 1e9
	return Generate(start)
}
