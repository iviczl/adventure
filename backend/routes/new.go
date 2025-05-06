package routes

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"text-adventure/constants"
	"text-adventure/middlewares"
	"text-adventure/models"
	"text-adventure/models/dbmodels"
	"text-adventure/utils"

	"github.com/gin-gonic/gin"
)

func New(c *gin.Context) {
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
	userId, err := models.StringToGuid(middlewares.GetClaim(c, "user_id").(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	session, err := constants.SessionStore.New(c.Request, constants.SessionName)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	var serialized bytes.Buffer
	encoder := gob.NewEncoder(&serialized)
	err = encoder.Encode(&adventure)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	user := &dbmodels.User{}
	constants.DbClient.Where("id = ?", userId).First(user)
	if user.Id == models.ZeroGuid() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	play := dbmodels.Play{UserId: userId, AdventureTitle: adventure.Title, Adventure: serialized.Bytes()}
	result := constants.DbClient.Create(&play)
	if result.Error != nil {
		fmt.Println("Failed to create a play:", result.Error)
		return
	}
	adventure.Id = play.Id
	adventure.ActualPosition.ActualPositionAdventureId = play.Id
	adventure.ActualPosition.AdventureId = play.Id
	actualPosition := models.AdjustedActualPosition(adventure.ActualPosition)
	session.Values[fmt.Sprintf("%v:%v", userId, adventure.Id)] = adventure
	fmt.Println("Session entry key created:", fmt.Sprintf("%v:%v", userId, adventure.Id))
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, actualPosition)
}
