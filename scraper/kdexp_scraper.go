package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/falsy/delivery-tracker-server/model"
)

type kdexpResponse struct {
	Result string `json:"result"`
	Info   struct {
		SendName string `json:"send_name"`
		ReName   string `json:"re_name"`
	} `json:"info"`
	Items []struct {
		Tel      string `json:"tel"`
		Location string `json:"location"`
		RegDate  string `json:"reg_date"`
		Stat     string `json:"stat"`
	} `json:"items"`
}

func parseKDExpDateTime(value string) string {
	parts := strings.Split(value, ".")
	return parts[0]
}

func parseKDExpStatus(value string) model.DeliveryState {
	switch {
	case strings.Contains(value, "접수완료"):
		return model.DeliveryState{ID: "item_received", Name: "상품인수"}
	case strings.Contains(value, "배송완료"):
		return model.DeliveryState{ID: "delivered", Name: "배달완료"}
	default:
		return model.DeliveryState{ID: "in_transit", Name: "상품이동중"}
	}
}

func KDExpGetTrack(trackingNumber string) (*model.DeliveryResult, error) {
	url := fmt.Sprintf("https://kdexp.com/service/delivery/ajax_basic.do?barcode=%s", trackingNumber)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("운송장 조회에 실패하였습니다")
	}

	var data kdexpResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if data.Result != "suc" {
		return nil, errors.New("해당 운송장이 존재하지 않거나 조회할 수 없습니다")
	}

	progresses := make([]model.DeliveryProgress, 0, len(data.Items))
	for i := len(data.Items) - 1; i >= 0; i-- {
		row := data.Items[i]
		progresses = append(progresses, model.DeliveryProgress{
			Description: "연락처: " + row.Tel,
			Location:    row.Location,
			Time:        parseKDExpDateTime(row.RegDate),
			State:       parseKDExpStatus(row.Stat),
		})
	}

	state := parseKDExpStatus("")
	if len(progresses) > 0 {
		state = progresses[0].State
	}

	from := model.DeliveryLocation{
		Name: parseLocationName(data.Info.SendName),
		Time: func() string {
			if len(progresses) > 0 {
				return progresses[len(progresses)-1].Time
			}
			return ""
		}(),
	}

	to := model.DeliveryLocation{
		Name: parseLocationName(data.Info.ReName),
		Time: func() string {
			if state.Name == "배달완료" && len(progresses) > 0 {
				return progresses[0].Time
			}
			return ""
		}(),
	}

	return &model.DeliveryResult{
		From:       from,
		To:         to,
		Progresses: progresses,
		State:      state,
	}, nil
}
