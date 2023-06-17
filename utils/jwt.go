package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//TODO: remove when User DTO is ready
type User struct {
	ID    int
	Email string
	Name  string
	Role  string
}

func GenerateToken(user User) (string, error) {
	secretKey := GetJWTKey()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nameidentifier": user.ID,
		"emailaddress":   user.Email,
		"name":           user.Name,
		"role":           user.Role,
		"exp":            time.Now().Add(time.Hour * 1).Unix(),
		"iss": 			  GetJWTIssuer(),
		"aud":			  GetJWTAudience(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//run the ExampleToken in main after LoadEnv() if you want to preview the Token
func ExampleToken() {
	user := User{
		ID:    10,
		Email: "dimiv@abv.bg",
		Name:  "Dimi Vlassis",
		Role:  "Owner",
	}

	tokenString, err := GenerateToken(user)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return
	}

	fmt.Println("JWT Token:", tokenString)
}