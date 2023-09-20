package main

import (
	"encoding/base64"
	"fmt"
)

func (app *application) generateRegisterLink(email, role string) string {
	emailBase64 := base64.StdEncoding.EncodeToString([]byte(email))
	roleBase64 := base64.StdEncoding.EncodeToString([]byte(role))

	registerLink := fmt.Sprintf("/register?email=%s&role=%s", emailBase64, roleBase64)

	return registerLink
}
