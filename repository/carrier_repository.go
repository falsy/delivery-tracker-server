package repository

import (
	"database/sql"

	"github.com/falsy/delivery-tracker-server/db"
	"github.com/falsy/delivery-tracker-server/model"
)

func GetAllCarriers() ([]model.Carrier, error) {
	rows, err := db.DB.Query(`
    SELECT uid, no, name, displayName, isCrawlable, isPopupEnabled, popupURL
    FROM CarrierModels
  `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carriers []model.Carrier
	for rows.Next() {
		var c model.Carrier
		err := rows.Scan(&c.UID, &c.No, &c.Name, &c.DisplayName, &c.IsCrawlable, &c.IsPopupEnabled, &c.PopupURL)
		if err != nil {
			return nil, err
		}
		carriers = append(carriers, c)
	}

	return carriers, nil
}

func GetCarrierByID(id string) (*model.Carrier, error) {
	var c model.Carrier
	err := db.DB.QueryRow(`
    SELECT uid, no, name, displayName, isCrawlable, isPopupEnabled, popupURL
    FROM CarrierModels WHERE uid = ?
  `, id).Scan(&c.UID, &c.No, &c.Name, &c.DisplayName, &c.IsCrawlable, &c.IsPopupEnabled, &c.PopupURL)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &c, nil
}
