package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	adventureEngine "text-adventure/engine"
	"text-adventure/models"
	"text-adventure/utils"
	"time"

	sessionHandler "github.com/gorilla/sessions"

	"encoding/gob"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var port = 8080

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
	gob.Register(&adventureEngine.Adventure{})
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
	store := utils.NewInMemoryStore([]byte("WYVOEWXHLUDB34DAELL2LDYNJSLEZT5WZRLPLHIQE5JXOGYCGIZQ"))
	store.Options = &sessionHandler.Options{
		MaxAge:   86400 * 30,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
		// Domain:      os.Getenv("ALLOWED_ORIGIN"),
		Partitioned: true,
	}
	engine := gin.Default()
	engine.Use(CorsMiddleware())

	engine.GET("/games", func(c *gin.Context) { c.JSON(200, utils.AdventureInfos()) })

	engine.GET("/headers", headers)

	engine.POST("/new", func(c *gin.Context) {
		body, err := utils.RequestBody(c)
		if err != nil {
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
		temporaryDescription := ""
		if adventure.ActualPosition.TemporaryDescription != "" {
			temporaryDescription = adventure.ActualPosition.TemporaryDescription
			adventure.ActualPosition.TemporaryDescription = ""
		}
		session, err := store.Get(c.Request, "session") //sessions.Default(c)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		db.Create(adventure)
		actualPosition := models.AdjustedActualPosition(adventure.ActualPosition)
		actualPosition.Description = temporaryDescription
		session.Values[fmt.Sprintf("%v:%v", player, adventure.Id)] = adventure
		fmt.Println("Session entry key created:", fmt.Sprintf("%v:%v", player, adventure.Id))
		fmt.Println(session.Values[fmt.Sprintf("%v:%v", player, adventure.Id)])
		err = session.Save(c.Request, c.Writer)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Session ID:(", session.ID, ")", len(session.ID))
		c.JSON(200, actualPosition)
	})

	engine.POST("/do", func(c *gin.Context) {
		body, err := utils.RequestBody(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		player := body["player"].(string)
		adventureId := body["adventureId"].(string)
		actionCode := body["actionCode"].(string)
		if player == "" || adventureId == "" || actionCode == "" {
			c.JSON(400, gin.H{"error": "Missing required parameters"})
			return
		}
		session, err := store.Get(c.Request, "session")
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		adventureValue := session.Values[fmt.Sprintf("%v:%v", player, adventureId)]
		fmt.Println("Session entry key:", fmt.Sprintf("%v:%v", player, adventureId))
		if adventureValue == nil {
			c.JSON(404, gin.H{"error": "Adventure not found"})
			return
		}
		adventure := &adventureEngine.Adventure{}
		adventure, ok := adventureValue.(*adventureEngine.Adventure)
		if !ok {
			c.JSON(500, gin.H{"error": "Adventure unmarshaling failed. "})
			return
		}
		// Fixing unmarshaling issue with ActualPosition
		if adventure.ActualPosition != nil {
			actualPosition := adventureEngine.GetPositionFromPositionList(adventure.Positions, adventure.ActualPosition.Code)
			if actualPosition != nil {
				adventure.ActualPosition = actualPosition
			}
		}
		err = adventure.Do(actionCode)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		actualPosition := models.AdjustedActualPosition(adventure.ActualPosition)
		session.Values[fmt.Sprintf("%v:%v", player, adventure.Id)] = adventure
		err = session.Save(c.Request, c.Writer)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, actualPosition)
	})

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
