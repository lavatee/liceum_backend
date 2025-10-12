package endpoint

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (e *Endpoint) Middleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	sliceOfHeader := strings.Split(header, " ")
	if len(sliceOfHeader) != 2 || sliceOfHeader[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("invalid header")})
		return
	}
	token := sliceOfHeader[1]
	claims, err := e.services.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if !e.services.CheckIsAdmin(fmt.Sprint(claims["email"])) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("notadmin")})
		return
	}
}
