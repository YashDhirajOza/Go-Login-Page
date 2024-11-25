package main

import (
	
	"fmt"
	"log"
	"net/http"

	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var db *badger.DB

func main() {
	var err error
	db, err = badger.Open(badger.DefaultOptions("user_db"))
	if err != nil {
		log.Fatal("Failed to open Badger database:", err)
	}
	defer db.Close()

	r := gin.Default()
	r.Use(setupCors())

	r.POST("/api/login", login)
	r.POST("/api/register", register)

	r.Run(":8080")
}

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

func login(c *gin.Context) {
	var credentials User
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "status": "error"})
		return
	}

	var storedPassword string
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(credentials.Username))
		if err != nil {
			return fmt.Errorf("user not found")
		}
		return item.Value(func(v []byte) error {
			storedPassword = string(v)
			return nil
		})
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials", "status": "error"})
		return
	}

	if storedPassword != credentials.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials", "status": "error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "status": "ok"})
}

func register(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "status": "error"})
		return
	}

	err := db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(newUser.Username))
		if err == nil {
			return fmt.Errorf("user already exists")
		}
		return txn.Set([]byte(newUser.Username), []byte(newUser.Password))
	})

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": err.Error(), "status": "error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful", "status": "ok"})
}
