package utils

import (
	"github.com/golang-jwt/jwt"
)

// JWT Claims constants
const (
	Name         = "name"
	EmailAddress = "emailaddress"
	Role         = "role"
)

func CreateJWT(username string, userEmail string, userRole string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		Name:         username,
		EmailAddress: userEmail,
		Role:         userRole,
	})

	privateKey := GetJWTPrivateKey()

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
