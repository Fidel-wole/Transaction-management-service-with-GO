package routes

import (
	"github.com/Fidel-wole/Transaction_Management_Service/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine){
server.POST("/signup", controllers.CreateUser)
server.POST("/login", controllers.Login)
}