package controllers

import (
	"log"
	"net/http"

	db "github.com/Fidel-wole/Transaction_Management_Service/db"
	sqlc "github.com/Fidel-wole/Transaction_Management_Service/db/sqlc"
	"github.com/gin-gonic/gin"
)

func CreateAccount(c *gin.Context) {
	var account sqlc.CreateAccountParams
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	queries := db.GetQueries()
	createdaccount, err := queries.CreateAccount(c, account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdaccount)
}

func GetAccountByUserId(c *gin.Context) {
    // Retrieve userId from the context
    userId, exists := c.Get("userId")
    if !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
        return
    }

    // Perform type assertion to convert userId to int64
    userIdInt, ok := userId.(int64) 
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID type"})
        return
    }

    queries := db.GetQueries()
    account, err := queries.GetAccountByUserId(c, userIdInt)
	if err != nil {
		// Log the error for debugging
		log.Printf("Error retrieving account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve account"})
		return
	}

    c.JSON(http.StatusOK, account)
}
