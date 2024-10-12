package controllers

import (
	"context"
	"fmt"
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
		ID:      depositData.AccountID,
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

// TransferMoney handles transferring money between accounts.
func TransferMoney(c *gin.Context) {
	var transferData struct {
		SenderAccountNumber   string `json:"sender_acn"`
		ReceiverAccountNumber string `json:"receiver_acn"`
		Amount                string `json:"amount"`
	}

	// Bind the incoming JSON to the transferData struct
	if err := c.ShouldBindJSON(&transferData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data", "error": err.Error()})
		return
	}

	// Generate a new reference ID for the transaction
	referenceID := uuid.New().String()

	// Create a new instance of the `Queries` object.
	queries := db.GetQueries()

	// Start a transaction to ensure atomicity.
	tx, err := db.GetDB().Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to begin transaction", "error": err.Error()})
		return
	}
	defer tx.Rollback() // Roll back in case anything goes wrong.

	ctx := context.Background()
	qtx := queries.WithTx(tx)

	// Update the sender's balance.
	err = qtx.UpdateSenderBalance(ctx, sqlc.UpdateSenderBalanceParams{
		AccountNumber: transferData.SenderAccountNumber,
		Balance:       transferData.Amount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update sender balance", "error": err.Error()})
		return
	}

	// Update the receiver's balance.
	err = qtx.UpdateReceiverBalance(ctx, sqlc.UpdateReceiverBalanceParams{
		AccountNumber: transferData.ReceiverAccountNumber,
		Balance:       transferData.Amount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update receiver balance", "error": err.Error()})
		return
	}
	// get sender account id
	sender_account_id, err := qtx.GetAccountIDByAccountNumber(ctx, transferData.SenderAccountNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get sender account id", "error": err.Error()})
		return
	}
	// Create a debit transaction for the sender.
	err = qtx.CreateDebitTransaction(ctx, sqlc.CreateDebitTransactionParams{
		AccountID:   sender_account_id,
		Amount:      transferData.Amount,
		ReferenceID: referenceID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create debit transaction", "error": err.Error()})
		return
	}

	//get receiver account id
	receiver_account_id, err := qtx.GetAccountIDByAccountNumber(ctx, transferData.SenderAccountNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get receiver account id", "error": err.Error()})
		return
	}

	fmt.Printf("reference id: %v", referenceID)
	// Create a credit transaction for the receiver.
	err = qtx.CreateCreditTransaction(ctx, sqlc.CreateCreditTransactionParams{
		AccountID:   receiver_account_id,
		Amount:      transferData.Amount,
		ReferenceID: referenceID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create credit transaction", "error": err.Error()})
		return
	}

	// Commit the transaction if all operations succeeded.
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to commit transaction", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful", "reference_id": referenceID})
}
