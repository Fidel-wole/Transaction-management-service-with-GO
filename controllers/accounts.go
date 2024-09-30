package controllers

import (
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