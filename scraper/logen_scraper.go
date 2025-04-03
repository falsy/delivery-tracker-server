package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/falsy/delivery-tracker-server/model"
)

func parseLogenDateTime(value string) string {
	return strings.TrimSpace(value + ":00")
}

func parseLogenStatus(value string) model.DeliveryState {
	switch {
	case strings.Contains(value, "배송출고"):
		return model.DeliveryState{ID: "out_for_delivery", Name: "배달출발"}
	case strings.Contains(value, "배송완료"):
		return model.DeliveryState{ID: "delivered", Name: "배달완료"}
	default:
		return model.DeliveryState{ID: "in_transit", Name: "상품이동중"}
	}
}

func LogenGetTrack(trackingNumber string) (*model.DeliveryResult, error) {
	url := fmt.Sprintf("https://www.ilogen.com/web/personal/trace/%s", trackingNumber)
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

	content := doc.Find(".tab_contents")
	infoTable := content.Find("table")
	progressTable := content.Find("table").Eq(1)
	infoBody := infoTable.Find("tbody")

	if progressTable.Length() == 0 {
		return nil, errors.New("해당 운송장이 존재하지 않거나 조회할 수 없습니다")
	}

	var progresses []model.DeliveryProgress
	progressTable.Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
		td := s.Find("td")
		description := strings.TrimSpace(td.Eq(3).Text())
		location := strings.TrimSpace(td.Eq(1).Text())
		time := parseLogenDateTime(td.Eq(0).Text())
		state := parseLogenStatus(td.Eq(2).Text())

		progresses = append(progresses, model.DeliveryProgress{
			Description: description,
			Location:    location,
			Time:        time,
			State:       state,
		})
	})

	// reverse progresses
	for i, j := 0, len(progresses)-1; i < j; i, j = i+1, j-1 {
		progresses[i], progresses[j] = progresses[j], progresses[i]
	}

	state := parseLogenStatus("")
	if len(progresses) > 0 && progresses[0].State.Name == "배달완료" {
		state = progresses[0].State
	}

	from := model.DeliveryLocation{
		Name: parseLocationName(strings.TrimSpace(infoBody.Find("tr").Eq(3).Find("td").Eq(1).Text())),
		Time: func() string {
			if len(progresses) > 0 {
				return progresses[len(progresses)-1].Time
			}
			return ""
		}(),
	}

	to := model.DeliveryLocation{
		Name: parseLocationName(strings.TrimSpace(infoBody.Find("tr").Eq(3).Find("td").Eq(3).Text())),
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
