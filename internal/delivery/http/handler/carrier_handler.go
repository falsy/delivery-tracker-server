package handler

import (
	"github.com/falsy/delivery-tracker-server/pkg/response"
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
		response.Error(c, http.StatusInternalServerError, "DB 에러")
		return
	}
	response.Success(c, carriers, "택배사 목록 조회 성공")
}

func (h *CarrierHandler) GetCarrierByID(c *gin.Context) {
	id := c.Param("id")

	carrier, err := h.carrierService.GetCarrierByID(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "DB 에러")
		return
	}
	if carrier == nil {
		response.Error(c, http.StatusNotFound, "존재하지 않는 택배사")
		return
	}

	response.Success(c, carrier, "택배사 조회 성공")
}
