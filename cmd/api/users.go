package main

import (
	"example/bluebean-go/internal/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) regiserUserHandler(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BindJSON(&input); err != nil {
		panic("Wrong json format.")
	}

	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
		Role:  input.Role,
	}

	err := user.Password.Set(input.Password)
	if err != nil {
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, user)
}
