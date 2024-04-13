package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllBuildings(c *gin.Context) {

	buildings, err := h.services.Room.GetAllBuildings()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, err)
		return
	}

	c.JSON(http.StatusOK, buildings)

}

func (h *Handler) GetRoomsByBuilding(c *gin.Context) {
	buildingIDstr := c.Param("buildingId")

	buildingID, err := strconv.Atoi(buildingIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	rooms, err := h.services.Room.GetRoomsByBuildingId(buildingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *Handler) GetRoomNameById(c *gin.Context) {
	roomIdStr := c.Param("roomId")
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	room, err := h.services.Room.GetRoomById(roomId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, room.Name)

}
