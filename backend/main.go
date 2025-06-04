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
	"text-adventure/models/dbmodels"
	"text-adventure/routes"
	"text-adventure/session"
	"time"

	"encoding/gob"

	"text-adventure/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func InitSessionStore() {
	constants.SessionStore = session.NewInMemoryStore([]byte(constants.SessionSecret))
	constants.SessionStore.Options = &sessions.Options{
		MaxAge:   86400 * 30,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
		// Domain:      os.Getenv("ALLOWED_ORIGIN"),
		Partitioned: true,
	}
}

func InitDbClient() {
	var err error
	constants.DbClient, err = gorm.Open(sqlite.Open(constants.DbPath), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
	}
	// Migrate the schema (create tables, etc.)
	err = constants.DbClient.AutoMigrate(&dbmodels.User{}, &dbmodels.Play{})
	if err != nil {
		fmt.Println("Failed to migrate the database:", err)
		return
	}
	var count int64
	constants.DbClient.Model(&dbmodels.User{}).Count(&count)
	if count == 0 {
		user := dbmodels.User{UserName: "admin", Email: "admin@admin.com", Password: "admin", RegistrationCode: "sasa"}
		result := constants.DbClient.Create(&user)
		if result.Error != nil {
			fmt.Println("Failed to create user:", result.Error)
			return
		}
		user = dbmodels.User{UserName: "alfa", Email: "alfa@admin.com", Password: "alfa", RegistrationCode: "sasa"}
		result = constants.DbClient.Create(&user)
		if result.Error != nil {
			fmt.Println("Failed to create user:", result.Error)
			return
		}
		user = dbmodels.User{UserName: "beta", Email: "beta@admin.com", Password: "beta", RegistrationCode: "sasa"}
		result = constants.DbClient.Create(&user)
		if result.Error != nil {
			fmt.Println("Failed to create user:", result.Error)
			return
		}
	}
}

func init() {
	gob.Register(make(map[string]interface{}))
	gob.Register(&models.Item{})
	gob.Register(&models.Player{})
	gob.Register(&models.Action{})
	gob.Register(&models.Position{})
	gob.Register(&engine.Adventure{})
}

func main() {
	fmt.Println("Service application launched.", time.Now().Format("2006-01-02 15:04:05"))
	godotenv.Load()
	shutDownSignals := make(chan os.Signal, 1)
	signal.Notify(shutDownSignals, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Initializing database client...")
	InitDbClient()

	fmt.Println("Database client initialized.")
	fmt.Println("Initializing HTTP server...")
	InitSessionStore()
	engine := gin.Default()

	engine.Use(CorsMiddleware())

	engine.POST("/login", routes.Login)
	engine.GET("/games", routes.Games)

	engine.Use(middlewares.Authentication())

	engine.POST("/new", routes.New)
	engine.POST("/do", routes.Do)
	engine.POST("/save", routes.Save)
	engine.GET("/plays", routes.Plays)
	engine.POST("/load", routes.Load)

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
