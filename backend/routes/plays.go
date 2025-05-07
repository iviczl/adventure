package routes

import (
	"text-adventure/constants"
	"text-adventure/middlewares"
	"text-adventure/models"
	"text-adventure/models/dbmodels"
	"text-adventure/types"

	"github.com/gin-gonic/gin"
)

func Plays(c *gin.Context) {
	userId, err := types.StringToGuid(middlewares.GetClaim(c, "user_id").(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var plays []dbmodels.Play
	constants.DbClient.Where("user_id = ?", userId).Find(&plays)
	dtoPlays := []models.DtoPlay{}
	for i := range plays {
		dtoPlays = append(dtoPlays, models.PlayToDtoPlay(plays[i]))
	}
	c.JSON(200, dtoPlays)
}
