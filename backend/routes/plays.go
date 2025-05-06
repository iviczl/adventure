package routes

import (
	"text-adventure/constants"
	"text-adventure/middlewares"
	"text-adventure/models"
	"text-adventure/models/dbmodels"

	"github.com/gin-gonic/gin"
)

func Plays(c *gin.Context) {
	userId, err := models.StringToGuid(middlewares.GetClaim(c, "user_id").(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var plays []dbmodels.Play
	constants.DbClient.Where("user_id = ?", userId).Find(&plays)
	dtoPlays := []models.DtoPlay{}
	for i := range plays {
		dtoPlays = append(dtoPlays, models.DtoPlay{Id: plays[i].Id, AdventureTitle: plays[i].AdventureTitle, CreatedAt: plays[i].CreatedAt, UpdatedAt: plays[i].UpdatedAt})
	}
	c.JSON(200, dtoPlays)
}
