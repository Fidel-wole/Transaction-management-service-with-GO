package controllers

import (
	"net/http"

	db "github.com/Fidel-wole/Transaction_Management_Service/db"
	sqlc "github.com/Fidel-wole/Transaction_Management_Service/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DepositMoney(c *gin.Context) {
	var depositData sqlc.DepositParams
	if err := c.ShouldBindJSON(&depositData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data", "error": err.Error()})
		return
	}
	referenceID := uuid.New().String()
	depositData.ReferenceID = referenceID
	queries := db.GetQueries()

	// Start a transaction
	tx, err := db.GetDB().Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to start transaction", "error": err.Error()})
		return
	}

	// Create the deposit transaction
	_, err = queries.Deposit(c, depositData) 
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to deposit money", "error": err.Error()})
		return
	}
	updateBalanceParams := sqlc.UpdateAccountBalanceParams{
		ID: depositData.AccountID,
		Balance: depositData.Amount,
	}
	// Update the account balance
	err = queries.UpdateAccountBalance(c, updateBalanceParams) 
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update account balance", "error": err.Error()})
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to commit transaction", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit successful"})
}
