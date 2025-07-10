package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func SignToken(userId, username, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiry := os.Getenv("JWT_EXPIRY")

	claims := jwt.MapClaims{
		"uid":      userId,
		"username": username,
		"role":     role,
	}

	if jwtExpiry != "" {
		duration, err := time.ParseDuration(jwtExpiry)
		if err != nil {
			return "", HandleError(err, "Err: JWT expiry parsing failed!")
		}

		claims["exp"] = jwt.NewNumericDate(time.Now().Add(duration))
	} else {
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", HandleError(err, "Err JWT signing failed")
	}
	return signedToken, nil
}
