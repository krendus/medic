package database

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	//specify the header key to hold the value which is the token for the user
	clientToken := c.Request.Header.Get("token")
	if clientToken == "" {
		msg := "No Authorization header provided"
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		c.Abort()
		return
	}

	claims, err := ValidateToken(clientToken)
	if err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}

	c.Set("uid", claims.Uid)

	c.Next()

}
