package main

import (
	"net/http"
	"time"

	"github.com/alireza-akbarzadeh/restful-app/internal/database"
	"github.com/alireza-akbarzadeh/restful-app/internal/messages"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type registeredRoutes struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=3"`
}
type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
type loginResponse struct {
	Token string `json:"token"`
}

func (app *application) login(c *gin.Context) {
	var auth loginRequest
	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existingUser, err := app.models.Users.GetByEmail(auth.Email)
	if existingUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid email or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Something wend wrong"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid email or password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": existingUser.Id})
	tokenString, err := token.SignedString([]byte(app.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, loginResponse{Token: tokenString})
}

func (app *application) registerUser(c *gin.Context) {

	var input registeredRoutes
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"expr":  time.Now().Add(time.Hour * 72).Unix(),
		})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrFailedToRegisterUser})
		return
	}
	input.Password = string(hashedPassword)

	user := database.User{
		Email:    input.Email,
		Password: input.Password,
		Name:     input.Name,
	}
	_, err = app.models.Users.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrFailedToRegisterUser, "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}
