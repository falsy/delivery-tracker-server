package middleware

import (
	"fmt"
	"time"

	"github.com/falsy/delivery-tracker-server/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupGlobalMiddleware(r *gin.Engine, cfg *config.Config) {
	r.Use(gzip.Gzip(gzip.BestSpeed))

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if cfg.Environment == "RELEASE" {
		r.Use(cors.New(cors.Config{
			AllowOriginFunc: func(origin string) bool {
				return origin == fmt.Sprintf("chrome-extension://%s", cfg.ExtensionID) ||
					origin == cfg.DevClientURL
			},
			AllowCredentials: true,
			ExposeHeaders:    []string{"ETag"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			MaxAge:           12 * time.Hour,
		}))
	} else {
		r.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowCredentials: true,
			ExposeHeaders:    []string{"ETag"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			MaxAge:           12 * time.Hour,
		}))
	}

	r.Use(gin.Logger())

	r.Use(gin.Recovery())
}
