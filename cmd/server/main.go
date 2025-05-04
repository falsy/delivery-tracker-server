package main

import (
	"fmt"
	"log"

	"github.com/falsy/delivery-tracker-server/internal/config"
	"github.com/falsy/delivery-tracker-server/internal/delivery/http/handler"
	"github.com/falsy/delivery-tracker-server/internal/delivery/http/router"
	"github.com/falsy/delivery-tracker-server/internal/repository/mysql"
	"github.com/falsy/delivery-tracker-server/internal/repository/scraper"
	"github.com/falsy/delivery-tracker-server/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// 설정 로드
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatal("잘못된 설정:", err)
	}

	if cfg.Environment == "RELEASE" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 데이터베이스 초기화
	mysql.Init()
	defer mysql.Close()

	carrierRepo := mysql.NewCarrierRepository(mysql.DB)
	trackerRepo := scraper.NewTrackerRepository()

	carrierService := service.NewCarrierService(carrierRepo)
	trackerService := service.NewTrackerService(carrierRepo, trackerRepo)

	carrierHandler := handler.NewCarrierHandler(carrierService)
	trackerHandler := handler.NewTrackerHandler(trackerService)

	// 라우터 설정
	r := router.SetupRouter(carrierHandler, trackerHandler, cfg)

	// 서버 시작
	address := fmt.Sprintf(":%s", cfg.Port)
	log.Println("서버 시작", address)
	if err := r.Run(address); err != nil {
		log.Fatal("서버 실행 실패:", err)
	}
}
