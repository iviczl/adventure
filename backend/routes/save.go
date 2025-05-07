package routes

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"text-adventure/constants"
	"text-adventure/engine"
	"text-adventure/middlewares"
	"text-adventure/models/dbmodels"
	"text-adventure/types"
	"text-adventure/utils"

	"github.com/gin-gonic/gin"
)

func Save(c *gin.Context) {
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
	session, err := constants.SessionStore.Get(c.Request, constants.SessionName)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	userId, err := types.StringToGuid(middlewares.GetClaim(c, "user_id").(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
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
	var serialized bytes.Buffer
	encoder := gob.NewEncoder(&serialized)
	err = encoder.Encode(&adventure)
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
	play.Adventure = serialized.Bytes()
	result := constants.DbClient.Save(&play)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Adventure saved successfully."})
}
