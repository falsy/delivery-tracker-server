package main

import (
	"log"
	"net/http"
	"os"

	"github.com/falsy/delivery-tracker-server/controller"
	"github.com/falsy/delivery-tracker-server/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	db.Init()

	port := os.Getenv("PORT")
	//extensionID := os.Getenv("EXTENSION_ID")
	//devClientURL := os.Getenv("DEV_CLIENT_URL")

	r := gin.Default()

	r.Use(gzip.Gzip(gzip.BestSpeed))

	r.Use(cors.New(cors.Config{
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == fmt.Sprintf("chrome-extension://%s", extensionID) || origin == devClientURL
		//},
		AllowAllOrigins:  true,
		AllowCredentials: true,
		ExposeHeaders:    []string{"ETag"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
	})

	r.GET("/carriers", controller.GetAllCarriers)
	r.GET("/carrier/:id", controller.GetCarrierByID)
	r.GET("/delivery/:carrierId/:trackingNumber", controller.GetDelivery)

	log.Printf("Server running on port %s\n", port)
	r.Run(":" + port)
}
