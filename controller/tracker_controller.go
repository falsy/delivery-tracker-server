package controller

import (
	"net/http"

	"github.com/falsy/delivery-tracker-server/repository"
	"github.com/gin-gonic/gin"
)

func GetDelivery(c *gin.Context) {
	carrierId := c.Param("carrierId")
	trackingNumber := c.Param("trackingNumber")

	repo := repository.NewTrackerRepository()

	carrier, err := repository.GetCarrierByID(carrierId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB 에러"})
		return
	}

	result, err := repo.GetDelivery(carrier, trackingNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
