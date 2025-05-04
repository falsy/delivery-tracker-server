package service

import (
	"github.com/falsy/delivery-tracker-server/internal/domain/model"
	"github.com/falsy/delivery-tracker-server/internal/domain/repository"
)

type TrackerService interface {
	GetDelivery(carrierID string, trackingNumber string) (*model.DeliveryResult, error)
}

type trackerService struct {
	carrierRepo repository.CarrierRepository
	trackerRepo repository.TrackerRepository
}

func NewTrackerService(
	carrierRepo repository.CarrierRepository,
	trackerRepo repository.TrackerRepository,
) TrackerService {
	return &trackerService{
		carrierRepo: carrierRepo,
		trackerRepo: trackerRepo,
	}
}

func (s *trackerService) GetDelivery(carrierID string, trackingNumber string) (*model.DeliveryResult, error) {
	carrier, err := s.carrierRepo.GetByID(carrierID)
	if err != nil {
		return nil, err
	}
	if carrier == nil {
		return nil, ErrCarrierNotFound
	}

	return s.trackerRepo.GetDelivery(carrier, trackingNumber)
}

var (
	ErrCarrierNotFound = &customError{"택배사를 찾을 수 없음"}
)

type customError struct {
	message string
}

func (e *customError) Error() string {
	return e.message
}
