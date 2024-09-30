package main

import (
	"github.com/Fidel-wole/Transaction_Management_Service/db"
	"github.com/Fidel-wole/Transaction_Management_Service/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	r := gin.Default()
	routes.RegisterRoutes(r)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Transaction Management Service",
		})
	})
	r.Run(":8080")
}