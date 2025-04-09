package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"text-adventure/models"
	"text-adventure/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var port = 8080

// func headers(w http.ResponseWriter, req *http.Request) {
// 	for name, headers := range req.Header {
// 		for _, h := range headers {
// 			fmt.Fprintf(w, "%v: %v\n", name, h)
// 		}
// 	}
// }

// type ListenHandler struct{}

// func (a *ListenHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	defer fmt.Println("Server has exited.")
// 	ctx := req.Context()
// 	<-ctx.Done()
// 	err := ctx.Err()
// 	fmt.Println("server:", err)
// 	internalError := http.StatusInternalServerError
// 	http.Error(w, err.Error(), internalError)
// }

// func main() {
// 	fmt.Println("Initializing routing...")
// 	http.HandleFunc("/headers", headers)
// 	fmt.Println("Routing initialized.")

// 	fmt.Printf("Server is listening on port %v.\n", port)
// 	http.ListenAndServe(fmt.Sprintf(":%v", port), &ListenHandler{})
// }

func headers(c *gin.Context) {
	headerMap := gin.H{}
	for name, headers := range c.Request.Header {
		for _, h := range headers {
			headerMap[name] = h
		}
	}
	c.JSON(200, headerMap)
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PATCH, DELETE, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	fmt.Println("Service application launched.", time.Now().Format("2006-01-02 15:04:05"))
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Initializing database client...")
	db, err := gorm.Open(sqlite.Open(filepath.Join(".", "db", "text-adventure.db")), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}

	// Migrate the schema (create tables, etc.)
	err = db.AutoMigrate(&models.Item{}, &models.Player{}, &models.Action{}, &models.Position{}, &models.Adventure{})
	if err != nil {
		fmt.Println("Failed to migrate the database:", err)
		return
	}
	fmt.Println("Database client initialized.")
	fmt.Println("Initializing HTTP server...")
	engine := gin.Default()
	engine.Use(CorsMiddleware())

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	engine.GET("/games", func(c *gin.Context) { c.JSON(200, utils.AdventureInfos()) })
	engine.GET("/headers", headers)
	engine.POST("/new", func(c *gin.Context) {
		rawBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var body map[string]interface{}
		if err := json.Unmarshal(rawBody, &body); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		player := body["player"].(string)
		adventureCode := body["gameId"].(string)

		adventure, err := utils.Load(adventureCode, nil)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		adventure.Player.Name = player
		if err := adventure.Start(); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		db.Create(adventure)
		c.JSON(200, adventure.ActualPosition)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%v", port),
		Handler: engine,
	}

	go func() {
		fmt.Printf("HTTP server is listening on port %v. %v\n", port, time.Now().Format("2006-01-02 15:04:05"))
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
		}
	}()

	// router.Run(fmt.Sprintf(":%v", port)) // listen and serve on 0.0.0.0:8080
	signal := <-osSignals
	if signal != nil {
		fmt.Println("Received signal:", signal)
	}
	fmt.Println("Server has exited.")
}
