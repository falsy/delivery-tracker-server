package handler

import (
	"github.com/falsy/delivery-tracker-server/pkg/response"
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
			response.Error(c, http.StatusInternalServerError, "택배사를 찾을 수 없음")
		default:
			response.Error(c, http.StatusInternalServerError, "서버 내부 오류")
		}
		return
	}

	response.Success(c, result, "운송장 정보 조회 성공")
}
