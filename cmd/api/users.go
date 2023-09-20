package main

import (
	"net/http"

	"bitbucket.org/nemetschek-systems/bluebean-service/internal/data"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/errorconstants"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/utils"
	"bitbucket.org/nemetschek-systems/bluebean-service/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (app *application) registerUserHandler(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"json": errorconstants.InvalidJSONFormatError.Error()})
		return
	}

	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
		Role:  input.Role,
	}

	err := user.Password.Set(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	v := validator.New()
	if data.ValidateRegisterInput(v, user); !v.Valid() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors})
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		if err == errorconstants.DuplicateEmailError {
			c.JSON(http.StatusConflict, gin.H{"email": errorconstants.DuplicateEmailError.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	jwt, err := utils.CreateJWT(user.Name, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"jwt": jwt})
}

func (app *application) loginUserHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"json": errorconstants.InvalidJSONFormatError.Error()})
		return
	}

	v := validator.New()
	if data.ValidateLoginInput(v, input.Email, input.Password); !v.Valid() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": v.Errors})
		return
	}

	user, err := app.models.Users.Get(input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"user": errorconstants.FailedLoginError.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"user": errorconstants.FailedLoginError.Error()})
		return
	}

	userLoggedIn, err := app.models.Users.CanLoginUser(input.Password, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	if !userLoggedIn {
		c.JSON(http.StatusNotFound, gin.H{"error": errorconstants.FailedLoginError.Error()})
		return
	}

	jwt, err := utils.CreateJWT(user.Name, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": jwt})
}

func (app *application) getAllFacilitiesForUserHandler(c *gin.Context) {
	email := c.Param("email")

	claims, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	userEmail, exists := claims.(jwt.MapClaims)[utils.EmailAddress].(string)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	if email != userEmail {
		c.JSON(http.StatusForbidden, gin.H{"error": errorconstants.UserIsNotAuthorizedError.Error()})
		return
	}

	user, err := app.models.Users.Get(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"user": errorconstants.UserNotFoundError.Error()})
		return
	}

	facilities, err := app.models.Users.GetAllFacilitiesForUser(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorconstants.InternalServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, facilities)
}
