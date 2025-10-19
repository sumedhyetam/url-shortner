package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sumedhyetam/url-shortner/api/database"
	"github.com/sumedhyetam/url-shortner/api/models"
	"github.com/sumedhyetam/url-shortner/api/utils"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func ShortenURL(c *gin.Context) {

	var body models.Request

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot parse JSON",
		})
		return
	}

	r2 := database.CreateRedisClient(1)
	defer r2.Close()

	val, err := r2.Get(database.Ctx, c.ClientIP()).Result()

	//fmt.Println(val, "val")

	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.ClientIP(), "10", 30*60*time.Second)
	} else {
		val, _ = r2.Get(database.Ctx, c.ClientIP()).Result()
		//fmt.Println(val, "val")
		valInt, _ := strconv.Atoi(val)
		//fmt.Println(valInt, "valInt")
		if valInt <= 0 {

			limit, _ := r2.TTL(database.Ctx, c.ClientIP()).Result()
			//fmt.Println(limit)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
			return
		}
	}

	if !govalidator.IsURL(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid URL",
		})
		return
	}

	if !utils.IsDifferentDomain(body.URL) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "You can't hack this System :)",
		})
		return
	}

	body.URL = utils.EnsureHttpPrefix(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateRedisClient(0)
	defer r.Close()

	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "URL Custom Short Already Exists",
		})
		return
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to connect to the redi server",
		})
		return
	}

	resp := models.Response{
		URL:            body.URL,
		CustomShort:    id,
		Expiry:         body.Expiry,
		XRateRemaining: 10,
		XRateLimitRest: 30,
	}

	r2.Decr(database.Ctx, c.ClientIP())

	val, _ = r2.Get(database.Ctx, c.ClientIP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)
	ttl, _ := r2.TTL(database.Ctx, c.ClientIP()).Result()
	resp.XRateLimitRest = ttl / time.Nanosecond / time.Minute
	resp.CustomShort = "localhost:3000" + "/" + id
	c.JSON(http.StatusOK, resp)

	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "URL shortening endpoint - implementation needed",
	// })
}
