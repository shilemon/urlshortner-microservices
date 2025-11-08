package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const pythonServiceURL = "http://localhost:5000"

type ShortenRequest struct {
	LongURL string `json:"long_url" binding:"required"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
	LongURL   string `json:"long_url"`
}

type ClickEvent struct {
	ShortCode string `json:"short_code"`
	ClickedAt string `json:"clicked_at"`
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./go.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_code TEXT UNIQUE NOT NULL,
		long_url TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database initialized successfully")
}

func generateShortCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	encoded := base64.URLEncoding.EncodeToString(b)
	// Take first 6 characters and remove any special chars
	shortCode := encoded[:6]
	return shortCode
}

func createShortURL(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortCode := generateShortCode()

	// Check if short code already exists (unlikely but possible)
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM urls WHERE short_code = ?", shortCode).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Regenerate if exists (very rare)
	for exists > 0 {
		shortCode = generateShortCode()
		db.QueryRow("SELECT COUNT(*) FROM urls WHERE short_code = ?", shortCode).Scan(&exists)
	}

	_, err = db.Exec("INSERT INTO urls (short_code, long_url) VALUES (?, ?)", shortCode, req.LongURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		return
	}

	response := ShortenResponse{
		ShortCode: shortCode,
		ShortURL:  "http://localhost:8000/" + shortCode,
		LongURL:   req.LongURL,
	}

	log.Printf("Created short URL: %s -> %s", shortCode, req.LongURL)
	c.JSON(http.StatusOK, response)
}

func redirect(c *gin.Context) {
	shortCode := c.Param("code")

	var longURL string
	err := db.QueryRow("SELECT long_url FROM urls WHERE short_code = ?", shortCode).Scan(&longURL)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Send event to Python service asynchronously
	go sendClickEvent(shortCode)

	// Redirect to the long URL
	c.Redirect(http.StatusMovedPermanently, longURL)
}

func sendClickEvent(shortCode string) {
	event := ClickEvent{
		ShortCode: shortCode,
		ClickedAt: time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Post(pythonServiceURL+"/api/events", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending event to Python service: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Python service returned status: %d", resp.StatusCode)
	} else {
		log.Printf("Click event sent for short code: %s", shortCode)
	}
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	// Routes
	r.POST("/api/shorten", createShortURL)
	r.GET("/:code", redirect)

	log.Println("Go service starting on :8000")
	r.Run(":8000")
}
