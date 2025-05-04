package mysql

import (
	"database/sql"

	"github.com/falsy/delivery-tracker-server/internal/domain/model"
	"github.com/falsy/delivery-tracker-server/internal/domain/repository"
)

type carrierRepository struct {
	db *sql.DB
}

func NewCarrierRepository(db *sql.DB) repository.CarrierRepository {
	return &carrierRepository{db: db}
}

func (r *carrierRepository) GetAll() ([]model.Carrier, error) {
	rows, err := r.db.Query(`
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return carriers, nil
}

func (r *carrierRepository) GetByID(id string) (*model.Carrier, error) {
	var c model.Carrier
	err := r.db.QueryRow(`
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
