package controllers

import (
	"net/http"

	db "github.com/Fidel-wole/Transaction_Management_Service/db"
	sqlc "github.com/Fidel-wole/Transaction_Management_Service/db/sqlc"
	"github.com/Fidel-wole/Transaction_Management_Service/utils"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
    var user sqlc.CreateUserParams
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data", "error": err.Error()})
        return
    }

    queries := db.GetQueries()

    // Check if the user already exists
    existingUser, err := queries.GetUserByEmail(c, user.Email)
    if err == nil && existingUser.Email == user.Email {
        c.JSON(http.StatusConflict, gin.H{"message": "User already exists"})
        return
    }

    // Hash the user's password
    hashedPassword, err := utils.HashPassword(user.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password", "error": err.Error()})
        return
    }
    user.Password = hashedPassword

    // Create the user in the database
    createdUser, err := queries.CreateUser(c, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err.Error()})
        return
    }

    accountParams := sqlc.CreateAccountParams{
        UserID: int64(createdUser.ID),  
        AccountNumber: utils.GenerateAccountNumber(),
		Currency: "USD",
		Balance: "0.00",
    }

    createdAccount, err := queries.CreateAccount(c, accountParams)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create account", "error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "User and account created successfully",
        "user":    createdUser,
        "account": createdAccount,
    })
}


func Login(c *gin.Context) {
	var user sqlc.GetUserByEmailRow
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data", "error": err.Error()})
	}
	queries := db.GetQueries()

	existingUser, err := queries.GetUserByEmail(c, user.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	err = utils.ComparePassword(existingUser.Password, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
		return
	}
    token, err := utils.GenerateToken(user.Email, int64(user.ID))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "data": token})
}
