package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/falsy/delivery-tracker-server/model"
)

func parseHanjinDateTime(date, time string) string {
	return strings.TrimSpace(date) + " " + strings.TrimSpace(time) + ":00"
}

func parseHanjinStatus(value string) model.DeliveryState {
	switch {
	case strings.Contains(value, "집하"):
		return model.DeliveryState{ID: "item_received", Name: "상품인수"}
	case strings.Contains(value, "배송출발"):
		return model.DeliveryState{ID: "out_for_delivery", Name: "배달출발"}
	case strings.Contains(value, "배송완료"):
		return model.DeliveryState{ID: "delivered", Name: "배달완료"}
	default:
		return model.DeliveryState{ID: "in_transit", Name: "상품이동중"}
	}
}

func HanjinGetTrack(trackingNumber string) (*model.DeliveryResult, error) {
	url := fmt.Sprintf("https://www.hanjin.com/kor/CMS/DeliveryMgr/WaybillResult.do?wblnum=%s&mCode=MN038&schLang=KR", trackingNumber)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
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

	wrap := doc.Find("#delivery-wr")
	if wrap.Length() == 0 {
		return nil, errors.New("해당 운송장이 존재하지 않거나 조회할 수 없습니다")
	}

	infoTds := wrap.Find(".delivery-tbl").Find("tbody").Find("td")
	progressTable := wrap.Find(".waybill-tbl").Find("table")

	var progresses []model.DeliveryProgress
	progressTable.Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
		td := s.Find("td")
		description := strings.TrimSpace(td.Eq(3).Text())
		location := strings.TrimSpace(td.Eq(2).Text())
		time := parseHanjinDateTime(td.Eq(0).Text(), td.Eq(1).Text())
		state := parseHanjinStatus(description)

		progresses = append([]model.DeliveryProgress{{
			Description: description,
			Location:    location,
			Time:        time,
			State:       state,
		}}, progresses...)
	})

	state := parseHanjinStatus("상품이동중")
	if len(progresses) > 0 {
		state = progresses[0].State
	}

	from := model.DeliveryLocation{
		Name: parseLocationName(strings.TrimSpace(infoTds.Eq(1).Text())),
		Time: func() string {
			if len(progresses) > 0 {
				return progresses[len(progresses)-1].Time
			}
			return ""
		}(),
	}

	to := model.DeliveryLocation{
		Name: parseLocationName(strings.TrimSpace(infoTds.Eq(2).Text())),
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
