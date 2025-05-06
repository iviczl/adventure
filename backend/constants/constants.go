package constants

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text-adventure/models/dbmodels"
	"text-adventure/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gorilla/sessions"
)

const SecretKey = "my_secret_key"
const SessionName = "session"
const SessionSecret = "WYVOEWXHLUDB34DAELL2LDYNJSLEZT5WZRLPLHIQE5JXOGYCGIZQ"

var dbPath = filepath.Join(".", "db", "text-adventure.db")

var SessionStore *utils.InMemoryStore
var DbClient *gorm.DB

func InitSessionStore() {
	SessionStore = utils.NewInMemoryStore([]byte(SessionSecret))
	SessionStore.Options = &sessions.Options{
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
	DbClient, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
	}
	// Migrate the schema (create tables, etc.)
	err = DbClient.AutoMigrate(&dbmodels.User{}, &dbmodels.Play{})
	if err != nil {
		fmt.Println("Failed to migrate the database:", err)
		return
	}
}

// user := dbmodels.User{UserName: "admin", Email: "admin@admin.com", Password: "admin", RegistrationCode: "sasa"}
// result := db.Create(&user)
// if result.Error != nil {
// 	fmt.Println("Failed to create user:", result.Error)
// 	return
// }
