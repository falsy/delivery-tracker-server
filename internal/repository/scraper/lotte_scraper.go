package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/falsy/delivery-tracker-server/internal/domain/model"
)

func parseLotteDateTime(value string) string {
	cleaned := strings.TrimSpace(value)
	parts := strings.Fields(cleaned)
	if len(parts) < 2 {
		return strings.TrimSpace(value)
	}

	date := parts[0]
	time := parts[1]
	if time == "--:--" {
		time += ":--"
	} else {
		time += ":00"
	}

	return strings.TrimSpace(date) + " " + strings.TrimSpace(time)
}

func parseLotteStatus(value string) model.DeliveryState {
	switch {
	case strings.Contains(value, "상품접수"):
		return model.DeliveryState{ID: "item_received", Name: "상품인수"}
	case strings.Contains(value, "배송 출발"):
		return model.DeliveryState{ID: "out_for_delivery", Name: "배달출발"}
	case strings.Contains(value, "배달 완료"):
		return model.DeliveryState{ID: "delivered", Name: "배달완료"}
	default:
		return model.DeliveryState{ID: "in_transit", Name: "상품이동중"}
	}
}

func LotteGetTrack(trackingNumber string) (*model.DeliveryResult, error) {
	url := fmt.Sprintf("https://www.lotteglogis.com/home/reservation/tracking/linkView?InvNo=%s", trackingNumber)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("운송장 조회에 실패하였습니다")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	wrap := doc.Find(".contArea")
	infoTable := wrap.Find("table").Eq(0)
	progressTable := wrap.Find("table").Eq(1)
	infoTds := infoTable.Find("tbody").Find("td")

	if infoTds.Length() == 1 {
		return nil, errors.New("해당 운송장이 존재하지 않거나 조회할 수 없습니다")
	}

	var progresses []model.DeliveryProgress
	progressTable.Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
		td := s.Find("td")
		description := strings.TrimSpace(td.Eq(3).Text())
		location := strings.TrimSpace(td.Eq(2).Text())
		timeHtml, _ := td.Eq(1).Html()
		time := parseLotteDateTime(timeHtml)
		state := parseLotteStatus(td.Eq(0).Text())

		progresses = append(progresses, model.DeliveryProgress{
			Description: description,
			Location:    location,
			Time:        time,
			State:       state,
		})
	})

	state := parseLotteStatus("상품이동중")
	if len(progresses) > 0 {
		state = progresses[0].State
	}

	from := model.DeliveryLocation{
		Name: strings.TrimSpace(infoTds.Eq(1).Text()),
		Time: func() string {
			if len(progresses) > 0 {
				return progresses[len(progresses)-1].Time
			}
			return ""
		}(),
	}

	to := model.DeliveryLocation{
		Name: strings.TrimSpace(infoTds.Eq(2).Text()),
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
