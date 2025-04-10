package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text-adventure/engine"
	"text-adventure/models"

	"github.com/gin-gonic/gin"
)

var adventureFilePath = "adventures"

func AdventureFilePath() string {
	return adventureFilePath
}

type AdventureInfo struct {
	Code        string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func AdventureInfos() []AdventureInfo {
	content, err := os.ReadDir(adventureFilePath)
	if err != nil {
		panic(err)
	}
	adventureInfos := []AdventureInfo{}
	for _, file := range content {
		fileName := file.Name()
		rawData, err := os.ReadFile(filepath.Join(adventureFilePath, fileName))
		if err != nil {
			panic(err)
		}
		adventure := &models.Adventure{}
		jsonErr := json.Unmarshal(rawData, &adventure)
		if jsonErr != nil {
			panic(jsonErr)
		}
		adventureInfos = append(adventureInfos, AdventureInfo{
			Code:        adventure.Code,
			Title:       adventure.Title,
			Description: adventure.Description,
		})

		// adventureCodes := []string{}
		// ext := filepath.Ext(fileName)
		// if !file.IsDir() && ext == ".json" {
		// 	adventureCodes = append(adventureCodes, strings.TrimSuffix(fileName, ext))
		// }
	}
	return adventureInfos
}

// Load loads an Adventure from a JSON file.
func Load(code string, adventure *engine.Adventure) (*engine.Adventure, error) {
	rawData, err := os.ReadFile(filepath.Join(AdventureFilePath(), fmt.Sprintf("%v.json", code)))
	if err != nil {
		return nil, err
	}
	if adventure == nil {
		adventure = &engine.Adventure{}
	}
	jsonErr := json.Unmarshal(rawData, adventure)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return adventure, jsonErr
	// return Build(rawData, adventure)
}

func RequestBody(c *gin.Context) (map[string]interface{}, error) {
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return nil, err
	}
	return body, nil
}
