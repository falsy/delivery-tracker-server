package repository

import "github.com/falsy/delivery-tracker-server/internal/domain/model"

type CarrierRepository interface {
	GetAll() ([]model.Carrier, error)
	GetByID(id string) (*model.Carrier, error)
}
