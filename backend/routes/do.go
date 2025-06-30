package routes

import (
	"fmt"
	"net/http"
	"text-adventure/constants"
	"text-adventure/engine"
	"text-adventure/middlewares"
	"text-adventure/models"
	"text-adventure/types"
	"text-adventure/utils"

	"github.com/gin-gonic/gin"
)

func Do(c *gin.Context) {
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
	userId, err := types.StringToGuid(middlewares.GetClaim(c, "user_id").(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	session, err := constants.SessionStore.Get(c.Request, constants.SessionName)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	adventureValue := session.Values[fmt.Sprintf("%v:%v", userId, adventureId)]
	if adventureValue == nil {
		c.JSON(404, gin.H{"error": "Adventure not found"})
		return
	}
	adventure := &engine.Adventure{}
	adventure, ok := adventureValue.(*engine.Adventure)
	if !ok {
		c.JSON(500, gin.H{"error": "Adventure unmarshaling failed. "})
		return
	}
	err = adventure.Do(actionCode)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	actualPosition := models.AdjustedActualPosition(adventure.ActualPosition)
	session.Values[fmt.Sprintf("%v:%v", userId, adventure.Id)] = adventure
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, &models.DtoPlayState{
		Player:         models.PlayerToDtoPlayer(adventure.Player),
		ActualPosition: models.PositionToDtoPosition(actualPosition),
	})
}
