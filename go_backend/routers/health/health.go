package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHealth(c *gin.Context) {
	msg, key := c.GetQuery("msg")
	if !key {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid request param for playing ping pong with GO",
		})
		c.Abort()
	}
	c.JSON(http.StatusOK, gin.H{
		"you_said": msg,
		"I_said":   "pong",
	})

}
