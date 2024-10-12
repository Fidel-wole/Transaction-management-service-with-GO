package routes

import (
	"github.com/Fidel-wole/Transaction_Management_Service/controllers"
	middleware "github.com/Fidel-wole/Transaction_Management_Service/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine){
server.POST("/signup", controllers.CreateUser)
server.POST("/login", controllers.Login)

authRoutes := server.Group("/auth")
authRoutes.Use(middleware.AuthMiddleware())
{
	authRoutes.GET("/account", controllers.GetAccountByUserId)
	authRoutes.POST("/deposit", controllers.DepositMoney)
	authRoutes.POST("/transfer", controllers.TransferMoney)
}

}