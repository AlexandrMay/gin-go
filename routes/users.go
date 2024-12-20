package routes

import (
	"gin-go/models"
	"gin-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func signup(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = user.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save a user"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func login(ctx *gin.Context) {
	var user models.User

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = user.ValidateCreds()

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization failed"})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
