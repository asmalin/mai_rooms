package handler

import (
	"classrooms/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	router.POST("/login", h.Login)
	router.GET("/checkAuth", h.userIdentity, h.checkAuth)

	router.POST("/tg_login", h.TgLogin)

	api := router.Group("/api", h.userIdentity)
	{
		api.POST("/reserve", h.Reserve)
		api.POST("/cancelReservation", h.CancelReservation)
		api.GET("/qrcodes", h.GetQRCodes)

	}

	api = router.Group("/api")
	{
		api.GET("/buildings", h.GetAllBuildings)
		api.GET("/rooms/:buildingId", h.GetRoomsByBuilding)
		api.GET("/room/:roomId", h.GetRoomNameById)
		api.GET("/schedule", h.GetScheduleByRoomAndDate)
		api.GET("/reserved_lesssons", h.GetReservedLessonsByRoomAndDate)
		api.GET("/all_reserved_lesssons", h.GetAllReservedLessons)
	}
	return router
}
