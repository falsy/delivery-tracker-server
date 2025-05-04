package router

import (
	"net/http"

	"github.com/falsy/delivery-tracker-server/internal/config"
	"github.com/falsy/delivery-tracker-server/internal/delivery/http/handler"
	"github.com/falsy/delivery-tracker-server/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	carrierHandler *handler.CarrierHandler,
	trackerHandler *handler.TrackerHandler,
	cfg *config.Config,
) *gin.Engine {
	r := gin.New()

	middleware.SetupGlobalMiddleware(r, cfg)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Delivery Tracker API Server"})
	})

	r.GET("/carriers", carrierHandler.GetAllCarriers)
	r.GET("/carrier/:id", carrierHandler.GetCarrierByID)
	r.GET("/delivery/:carrierId/:trackingNumber", trackerHandler.GetDelivery)

	return r
}
