package repository

import "github.com/falsy/delivery-tracker-server/internal/domain/model"

type TrackerRepository interface {
	GetDelivery(carrier *model.Carrier, trackingNumber string) (*model.DeliveryResult, error)
}
