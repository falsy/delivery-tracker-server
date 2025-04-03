package repository

import (
	"github.com/falsy/delivery-tracker-server/model"
	"github.com/falsy/delivery-tracker-server/scraper"
)

type TrackerRepository struct{}

func NewTrackerRepository() *TrackerRepository {
	return &TrackerRepository{}
}

func (r *TrackerRepository) GetDelivery(carrier *model.Carrier, trackingNumber string) (*model.DeliveryResult, error) {
	switch carrier.Name {
	case "cjlogistics":
		return scraper.CJLogisticsGetTrack(trackingNumber)
	case "epost":
		return scraper.EPostGetTrack(trackingNumber)
	case "daesin":
		return scraper.DaesinGetTrack(trackingNumber)
	case "hanjin":
		return scraper.HanjinGetTrack(trackingNumber)
	case "lotte":
		return scraper.LotteGetTrack(trackingNumber)
	case "kdexp":
		return scraper.KDExpGetTrack(trackingNumber)
	case "logen":
		return scraper.LogenGetTrack(trackingNumber)
	default:
		return nil, &UnknownCarrierError{Carrier: carrier.Name}
	}
}

type UnknownCarrierError struct {
	Carrier string
}

func (e *UnknownCarrierError) Error() string {
	return "지원하지 않는 택배사입니다: " + e.Carrier
}
