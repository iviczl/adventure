package routes

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"text-adventure/constants"
	"text-adventure/engine"
	"text-adventure/middlewares"
	"text-adventure/models"
	"text-adventure/models/dbmodels"
	"text-adventure/types"
	"text-adventure/utils"

	"github.com/gin-gonic/gin"
)

func Load(c *gin.Context) {
	body, err := utils.RequestBody(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	adventureId, err := types.StringToGuid(body["adventureId"].(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	userId, err := types.StringToGuid(middlewares.GetClaim(c, "user_id").(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	play := &dbmodels.Play{}
	constants.DbClient.Where("user_id = ?", userId).Where("id = ?", adventureId).First(&play)
	if play.Id == types.ZeroGuid() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Cannot find play"})
		return
	}
	adventure := &engine.Adventure{}
	decoder := gob.NewDecoder(bytes.NewBuffer(play.Adventure))
	err = decoder.Decode(&adventure)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	adventure.Id = play.Id
	adventure.ActualPosition.ActualPositionAdventureId = play.Id
	adventure.ActualPosition.AdventureId = play.Id

	session, err := constants.SessionStore.New(c.Request, constants.SessionName)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	actualPosition := models.AdjustedActualPosition(adventure.ActualPosition)
	session.Values[fmt.Sprintf("%v:%v", userId, adventure.Id)] = adventure
	fmt.Println("Session entry key created:", fmt.Sprintf("%v:%v", userId, adventure.Id))
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	playState := models.DtoPlayState{
		AdventureId:    adventure.Id.String(),
		Player:         adventure.Player.Name,
		ActualPosition: models.PositionToDtoPosition(actualPosition),
	}
	c.JSON(200, playState)
}
