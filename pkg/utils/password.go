package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

func PasswordValidate(userPassword, incomingPassword string) error {
	parts := strings.Split(userPassword, ".")
	if len(parts) != 2 {
		return HandleError(errors.New("invalid password"), "Err: Invalid password hash")
	}

	saltBase64 := parts[0]
	hashedPasswordBase64 := parts[1]
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return HandleError(err, "Password salt decode failed")
	}

	hashedPassword, err := base64.StdEncoding.DecodeString(hashedPasswordBase64)
	if err != nil {
		return HandleError(err, "Hashed password salt decode failed")
	}

	hash := argon2.IDKey([]byte(incomingPassword), salt, 1, 60*1024, 4, 32)
	if len(hash) != len(hashedPassword) {
		return HandleError(err, "Err: Password incorrect!")
	}

	// Compares the hashes now after validating the length;
	if subtle.ConstantTimeCompare(hash, hashedPassword) == 1 {
		// Do nothing (i.e, this is the positive case - Password Correct)
		return nil
	} else {
		return HandleError(err, "Err: Password incorrect!")
	}
}

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", HandleError(err, "Err: Cannot generate salt!")
	}

	// Password encoding;
	hash := argon2.IDKey([]byte(password), salt, 1, 60*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	return encodedHash, nil
}
