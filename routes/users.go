package routes

import (
	"net/http"

	"example.com/rest-api/models"
	"example.com/rest-api/utils"
	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data"})
		return
	}

	err = user.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not register a user. Try again later.",
			"error":   err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "New user registered"})
}

func login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse request data"})
		return
	}
	err = user.UserValidation()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error()})
		return
	}

	tokenStr, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not authorise user"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successfull!", "token": tokenStr})
}

func getUsers(context *gin.Context) {
	users, err := models.GetUsers()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch users.",
			"error":   err.Error()})
		return
	}
	context.JSON(http.StatusOK, users)
}
