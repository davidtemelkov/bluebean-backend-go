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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wrong json format"})
		return
	}

	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
		Role:  input.Role,
	}

	err := user.Password.Set(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}

	c.IndentedJSON(http.StatusCreated, user)
}

func (app *application) getUserByEmailHandler(c *gin.Context) {
	email := c.Param("email")

	user, err := app.models.Users.GetByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
