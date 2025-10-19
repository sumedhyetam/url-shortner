package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sumedhyetam/url-shortner/api/database"
	"github.com/sumedhyetam/url-shortner/api/models"
)

func EditURL(c *gin.Context) {
	shortID := c.Param("shortID")
	var body models.Request

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot Parse JSON",
		})
	}

	r := database.CreateRedisClient(0)
	defer r.Close()

	//check if the shortID exists in the DB or not

	val, err := r.Get(database.Ctx, shortID).Result()
	if err != nil || val == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ShortID doesn't exists",
		})
	}

	//update the content of the URL,expiry time with shortID

	err = r.Set(database.Ctx, shortID, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to update the shortend content",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "shortened link updated successfully",
	})
}
