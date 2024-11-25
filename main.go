package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Add CORS Middleware
	r.Use(setupCors())

	// Mock user data
	users := map[string]string{
		"john_doe": "password123",
	}

	r.POST("/api/login", func(c *gin.Context) {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "status": "error"})
			return
		}

		password, exists := users[credentials.Username]
		if !exists || password != credentials.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials", "status": "error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "status": "ok"})
	})

	r.Run(":8080") // Listen and serve on localhost:8080
}

// CORS Middleware
func setupCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}
