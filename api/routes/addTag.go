package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sumedhyetam/url-shortner/api/database"
)

type TagRequest struct {
	ShortID string `json:"shortID"`
	Tag     string `json:"tag"`
}

func AddTag(c *gin.Context) {
	var tagRequest TagRequest
	if err := c.ShouldBind(&tagRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Request Body",
		})
		return
	}

	shortId := tagRequest.ShortID
	tag := tagRequest.Tag
	r := database.CreateRedisClient(0)
	defer r.Close()

	val, err := r.Get(database.Ctx, shortId).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data not found for the given ShortID",
		})
		return
	}

	var data map[string]interface{}

	if err := json.Unmarshal([]byte(val), &data); err != nil {
		//if the data is not a JSON object ,assume it as plan string
		data = make(map[string]interface{})
		data["data"] = val
	}

	//check if "tags" field already exists and it's a slice of strings
	var tags []string
	if existingTags, ok := data["tags"].([]interface{}); ok {
		for _, t := range existingTags {
			if strTag, ok := t.(string); ok {
				tags = append(tags, strTag)
			}
		}
	}

	//check for duplicte tags
	for _, existingTag := range tags {
		if existingTag == tag {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Tag ALready Exists",
			})
			return
		}
	}

	// Add the new tag to the tags slice
	tags = append(tags, tag)
	data["tags"] = tags

	//Marshal the updated data back to JSON
	updatedData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to Marshal updated data",
		})
	}

	err = r.Set(database.Ctx, shortId, updatedData, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update the Database",
		})
		return
	}

	//Response with the updated data
	c.JSON(http.StatusOK, data)
}
