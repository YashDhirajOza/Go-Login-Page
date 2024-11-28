package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.etcd.io/bbolt"
)

// Database struct wraps the BoltDB instance
type Database struct {
	DB *bbolt.DB
}

// Open database connection
func NewDatabase(fileName string) (*Database, error) {
	db, err := bbolt.Open(fileName, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &Database{DB: db}, nil
}

// Authenticate user
func (db *Database) AuthenticateUser(username, password string) bool {
	var storedPassword string
	err := db.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("Users"))
		if bucket == nil {
			return nil
		}
		storedPassword = string(bucket.Get([]byte(username)))
		return nil
	})
	if err != nil || storedPassword == "" {
		return false
	}
	return storedPassword == password
}

func main() {
	// Initialize database
	db, err := NewDatabase("user_login.db")
	if err != nil {
		panic("Failed to open database!")
	}
	defer db.DB.Close()

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Next()
	})

	// API for user authentication
	r.POST("/login", func(c *gin.Context) {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		if db.AuthenticateUser(creds.Username, creds.Password) {
			c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		}
	})

	r.Run(":8080")
}
