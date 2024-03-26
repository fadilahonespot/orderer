package main

import (
	"github.com/fadilahonespot/orderer/controller"
	"github.com/fadilahonespot/orderer/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db := database.ConnectDB()

	orderController := controller.NewOrderController(db)

	router := gin.Default()
	router.POST("/orders", orderController.CreateOrder)
	router.GET("/orders", orderController.GetOrder)
	router.PUT("/orders/:orderId", orderController.UpdateOrder)
	router.DELETE("/orders/:orderId", orderController.DeleteOrder)
}
