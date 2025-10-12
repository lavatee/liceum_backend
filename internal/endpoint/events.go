package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *Endpoint) GetCurrentEvents(c *gin.Context) {
	events, err := e.services.Events.GetCurrentEvents()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (e *Endpoint) GetAllEvents(c *gin.Context) {
	events, err := e.services.Events.GetAllEvents()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}
