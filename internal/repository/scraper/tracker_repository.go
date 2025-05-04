package scraper

import (
	"fmt"

	"github.com/falsy/delivery-tracker-server/internal/domain/model"
	"github.com/falsy/delivery-tracker-server/internal/domain/repository"
)

type trackerRepository struct{}

func NewTrackerRepository() repository.TrackerRepository {
	return &trackerRepository{}
}

func (r *trackerRepository) GetDelivery(carrier *model.Carrier, trackingNumber string) (*model.DeliveryResult, error) {
	switch carrier.Name {
	case "cjlogistics":
		return CJLogisticsGetTrack(trackingNumber)
	case "epost":
		return EPostGetTrack(trackingNumber)
	case "daesin":
		return DaesinGetTrack(trackingNumber)
	case "hanjin":
		return HanjinGetTrack(trackingNumber)
	case "lotte":
		return LotteGetTrack(trackingNumber)
	case "kdexp":
		return KDExpGetTrack(trackingNumber)
	case "logen":
		return LogenGetTrack(trackingNumber)
	default:
		return nil, &UnknownCarrierError{Carrier: carrier.Name}
	}
}

type UnknownCarrierError struct {
	Carrier string
}

func (e *UnknownCarrierError) Error() string {
	return fmt.Sprintf("지원하지 않는 택배사: %s", e.Carrier)
}
