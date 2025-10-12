package endpoint

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lavatee/liceum_backend/internal/service"
)

type Endpoint struct {
	services *service.Service
}

func NewEndpoint(services *service.Service) *Endpoint {
	return &Endpoint{
		services: services,
	}
}

func (e *Endpoint) InitRoutes() *gin.Engine {
	router := gin.New()
	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))
	users := router.Group("/users")
	{
		users.GET("/current-events", e.GetCurrentEvents)
		users.GET("/all-events", e.GetAllEvents)
		users.GET("/event/:id", e.GetOneEvent)
		users.GET("/block/:id", e.GetOneBlock)
		users.POST("/send-code", e.SendAuthCode)
		users.POST("/verify-code", e.VerifyCode)
	}
	admins := router.Group("/admins", e.Middleware)
	{
		admins.POST("/events", e.PostEvent)
		admins.DELETE("/events/:id", e.DeleteEvent)
		admins.PUT("/events/:id", e.PutEvent)
		admins.POST("/blocks", e.PostEventBlock)
		admins.DELETE("/blocks/:id", e.DeleteEventBlock)
		admins.PUT("/blocks/:id", e.PutEventBlock)
	}
}
