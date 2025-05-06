package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text-adventure/constants"
	"text-adventure/engine"
	"text-adventure/models"
	"text-adventure/routes"
	"time"

	"encoding/gob"

	"text-adventure/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var port = 8080

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGIN"))
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

func init() {
	gob.Register(&engine.Adventure{})
	gob.Register(&models.Position{})
	gob.Register(&models.Action{})
	gob.Register(&models.Item{})
	gob.Register(&models.Player{})
	gob.Register(make(map[string]interface{}))
}

func main() {
	fmt.Println("Service application launched.", time.Now().Format("2006-01-02 15:04:05"))
	godotenv.Load()
	shutDownSignals := make(chan os.Signal, 1)
	signal.Notify(shutDownSignals, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Initializing database client...")
	constants.InitDbClient()

	fmt.Println("Database client initialized.")
	fmt.Println("Initializing HTTP server...")
	constants.InitSessionStore()
	engine := gin.Default()

	engine.Use(CorsMiddleware())

	engine.POST("/login", routes.Login)
	engine.GET("/games", routes.Games)

	engine.Use(middlewares.Authentication())

	engine.POST("/new", routes.New)
	engine.POST("/do", routes.Do)
	engine.POST("/save", routes.Save)

	go func() {
		fmt.Printf("HTTP server is listening on port %v. %v\n", port, time.Now().Format("2006-01-02 15:04:05"))
		err := engine.Run(fmt.Sprintf("%s:%v", os.Getenv("HOST"), port)) // listen and serve on 0.0.0.0:8080
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("HTTP server error:", err)
		}
	}()

	signal := <-shutDownSignals
	if signal != nil {
		fmt.Println("Received signal:", signal.String())
	}

	fmt.Println("HTTP server has exited.")
}
