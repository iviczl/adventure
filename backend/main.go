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

	"github.com/gin-contrib/sessions"
	// "github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memstore"

	"encoding/gob"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	gob.Register(adventureEngine.Adventure{})
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
	// store := cookie.NewStore([]byte("secret"))
	store := memstore.NewStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 0, SameSite: http.SameSiteNoneMode, Secure: true})
	engine := gin.Default()
	engine.Use(CorsMiddleware())
	engine.Use(sessions.Sessions("session", store))

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

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
		session := sessions.Default(c)
		db.Create(adventure)
		actualPosition := models.AdjustedActualPosition(adventure.ActualPosition)
		actualPosition.Description = temporaryDescription
		value, _ := json.Marshal(adventure)
		session.Set(fmt.Sprintf("%v:%v", player, adventure.Id), value)
		fmt.Println("Session entry key created:", fmt.Sprintf("%v:%v", player, adventure.Id))
		session.Save()
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
		session := sessions.Default(c)
		adventureValue := session.Get(fmt.Sprintf("%v:%v", player, adventureId))
		fmt.Println("Session entry key:", fmt.Sprintf("%v:%v", player, adventureId))
		if adventureValue == nil {
			c.JSON(404, gin.H{"error": "Adventure not found"})
			return
		}
		adventure := &adventureEngine.Adventure{}
		err = json.Unmarshal(adventureValue.([]byte), &adventure)
		if err != nil {
			c.JSON(500, gin.H{"error": "Adventure unmarshaling failed. " + err.Error()})
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
		value, _ := json.Marshal(adventure)
		session.Set(fmt.Sprintf("%v:%v", player, adventure.Id), value)
		// session.Set(fmt.Sprintf("%v:%v", player, adventure.Id), adventure)
		session.Save()
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
