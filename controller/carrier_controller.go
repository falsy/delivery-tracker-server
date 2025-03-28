package controller

import (
	"net/http"
	"strconv"

	"github.com/falsy/delivery-tracker-server/repository"
	"github.com/gin-gonic/gin"
)

func GetAllCarriers(c *gin.Context) {
	carriers, err := repository.GetAllCarriers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB 에러"})
		return
	}
	c.JSON(http.StatusOK, carriers)
}

func GetCarrierByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID"})
		return
	}

	carrier, err := repository.GetCarrierByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB 에러"})
		return
	}
	if carrier == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "존재하지 않는 택배사"})
		return
	}

	c.JSON(http.StatusOK, carrier)
}
