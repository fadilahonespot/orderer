package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fadilahonespot/orderer/controller"
	"github.com/fadilahonespot/orderer/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/dgrijalva/jwt-go"
)

func main() {
	godotenv.Load()
	db := database.ConnectDB()

	orderController := controller.NewOrderController(db)
	userController := controller.NewUserController(db)

	router := gin.Default()
	router.POST("/users/register", userController.Register)
	router.POST("/users/login", userController.Login)

	tx := router.Use(authMiddleware)
	tx.POST("/orders", orderController.CreateOrder)
	tx.GET("/orders", orderController.GetOrder)
	tx.PUT("/orders/:orderId", orderController.UpdateOrder)
	tx.DELETE("/orders/:orderId", orderController.DeleteOrder)

	router.Run(fmt.Sprintf(":%v", os.Getenv("PORT")))
}

func authMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	tokenString := authHeader[len("Bearer "):]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	c.Set("token", token)
	c.Next()
}
