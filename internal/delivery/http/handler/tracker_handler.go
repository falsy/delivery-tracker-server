package handler

import (
	"net/http"

	"github.com/falsy/delivery-tracker-server/internal/service"
	"github.com/gin-gonic/gin"
)

type TrackerHandler struct {
	trackerService service.TrackerService
}

func NewTrackerHandler(trackerService service.TrackerService) *TrackerHandler {
	return &TrackerHandler{
		trackerService: trackerService,
	}
}

func (h *TrackerHandler) GetDelivery(c *gin.Context) {
	carrierID := c.Param("carrierId")
	trackingNumber := c.Param("trackingNumber")

	result, err := h.trackerService.GetDelivery(carrierID, trackingNumber)
	if err != nil {
		switch err {
		case service.ErrCarrierNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Carrier not found"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}
