package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/falsy/delivery-tracker-server/model"
	"golang.org/x/net/html/charset"
)

func parseDaesinDateTime(value string) string {
	return strings.TrimSpace(value + ":00")
}

func parseDaesinStatus(value string) model.DeliveryState {
	if strings.Contains(value, "배송완료") {
		return model.DeliveryState{ID: "delivered", Name: "배달완료"}
	}
	return model.DeliveryState{ID: "in_transit", Name: "상품이동중"}
}

func DaesinGetTrack(trackingNumber string) (*model.DeliveryResult, error) {
	url := fmt.Sprintf("https://www.ds3211.co.kr/freight/internalFreightSearch.ht?billno=%s", trackingNumber)
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

	utf8Reader, err := charset.NewReader(res.Body, res.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		return nil, err
	}

	content := doc.Find("#printarea")
	tables := content.Find("table")
	if tables.Length() < 2 {
		return nil, errors.New("해당 운송장이 존재하지 않거나 조회할 수 없습니다")
	}

	infoTable := tables.Eq(0)
	progressTable := tables.Eq(1)
	infoTds := infoTable.Find("tbody")

	var progresses []model.DeliveryProgress
	progressTable.Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		td := s.Find("td")
		description := cleanEscapedText(strings.TrimSpace(td.Eq(2).Text()))
		location := strings.TrimSpace(td.Eq(1).Text())
		time := parseDaesinDateTime(td.Eq(3).Text())
		state := parseDaesinStatus(td.Eq(5).Text())

		progress := model.DeliveryProgress{
			Description: description,
			Location:    location,
			Time:        time,
			State:       state,
		}

		// 역순으로 추가 (최신 -> 오래된 순으로)
		progresses = append([]model.DeliveryProgress{progress}, progresses...)
	})

	state := parseDaesinStatus("")
	if len(progresses) > 0 && progresses[0].State.Name == "배달완료" {
		state = progresses[0].State
	}

	from := model.DeliveryLocation{
		Name: parseLocationName(strings.TrimSpace(infoTds.Find("tr").Eq(0).Find("td").Eq(0).Text())),
		Time: func() string {
			if len(progresses) > 0 {
				return progresses[len(progresses)-1].Time
			}
			return ""
		}(),
	}

	to := model.DeliveryLocation{
		Name: parseLocationName(strings.TrimSpace(infoTds.Find("tr").Eq(1).Find("td").Eq(0).Text())),
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
