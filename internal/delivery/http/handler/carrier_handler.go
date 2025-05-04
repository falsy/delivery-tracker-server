package handler

import (
	"net/http"

	"github.com/falsy/delivery-tracker-server/internal/service"
	"github.com/gin-gonic/gin"
)

type CarrierHandler struct {
	carrierService service.CarrierService
}

func NewCarrierHandler(carrierService service.CarrierService) *CarrierHandler {
	return &CarrierHandler{
		carrierService: carrierService,
	}
}

func (h *CarrierHandler) GetAllCarriers(c *gin.Context) {
	carriers, err := h.carrierService.GetAllCarriers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, carriers)
}

func (h *CarrierHandler) GetCarrierByID(c *gin.Context) {
	id := c.Param("id")

	carrier, err := h.carrierService.GetCarrierByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if carrier == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carrier not found"})
		return
	}

	c.JSON(http.StatusOK, carrier)
}
