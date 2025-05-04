package service

import (
	"github.com/falsy/delivery-tracker-server/internal/domain/model"
	"github.com/falsy/delivery-tracker-server/internal/domain/repository"
)

type CarrierService interface {
	GetAllCarriers() ([]model.Carrier, error)
	GetCarrierByID(id string) (*model.Carrier, error)
}

type carrierService struct {
	carrierRepo repository.CarrierRepository
}

func NewCarrierService(repo repository.CarrierRepository) CarrierService {
	return &carrierService{
		carrierRepo: repo,
	}
}

func (s *carrierService) GetAllCarriers() ([]model.Carrier, error) {
	return s.carrierRepo.GetAll()
}

func (s *carrierService) GetCarrierByID(id string) (*model.Carrier, error) {
	return s.carrierRepo.GetByID(id)
}
