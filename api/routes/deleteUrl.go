package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sumedhyetam/url-shortner/api/database"
)

func DeleteURL(c *gin.Context) {
	shortId := c.Param("shortID")

	r := database.CreateRedisClient(0)
	defer r.Close()

	err := r.Del(database.Ctx, shortId).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to Delete shortened Link",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Shortened URL Deleted Successfully",
	})
}
