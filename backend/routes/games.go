package routes

import (
	"text-adventure/utils"

	"github.com/gin-gonic/gin"
)

func Games(c *gin.Context) { c.JSON(200, utils.AdventureInfos()) }
