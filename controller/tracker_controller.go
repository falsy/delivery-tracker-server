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
	if carrier == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "존재하지 않는 택배사입니다"})
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
