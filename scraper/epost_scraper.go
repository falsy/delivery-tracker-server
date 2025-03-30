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

func parseEPostDateTime(value string) string {
	parts := strings.Fields(value)
	if len(parts) == 0 {
		return ""
	}
	date := strings.ReplaceAll(parts[0], ".", "-")
	if len(parts) > 1 {
		return date + " " + parts[1] + ":00"
	}
	return date
}

func parseEPostStatus(value string) model.DeliveryState {
	switch {
	case strings.Contains(value, "상품준비중"):
		return model.DeliveryState{ID: "preparing_item", Name: "상품준비중"}
	case strings.Contains(value, "접수"):
		return model.DeliveryState{ID: "item_received", Name: "상품인수"}
	case strings.Contains(value, "배달준비"):
		return model.DeliveryState{ID: "out_for_delivery", Name: "배달출발"}
	case strings.Contains(value, "배달완료"):
		return model.DeliveryState{ID: "delivered", Name: "배달완료"}
	default:
		return model.DeliveryState{ID: "in_transit", Name: "상품이동중"}
	}
}

func EPostGetTrack(trackingNumber string) (*model.DeliveryResult, error) {
	url := fmt.Sprintf("https://service.epost.go.kr/trace.RetrieveDomRigiTraceList.comm?sid1=%s", trackingNumber)
	res, err := http.Get(url)
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

	infoTable := doc.Find("#print").Find("table")
	progressTable := doc.Find("#processTable")
	infoTds := infoTable.Find("td")

	if infoTds.Length() == 0 {
		return nil, errors.New("해당 운송장이 존재하지 않거나 조회할 수 없습니다")
	}

	var progresses []model.DeliveryProgress
	progressTable.Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
		td := s.Find("td")
		descriptionRaw := td.Eq(3).Text()
		descriptionText := strings.Join(strings.Fields(descriptionRaw), " ")
		description := descriptionText
		if strings.Contains(description, "소포 물품 사진") {
			description = "접수"
		}
		location := td.Eq(2).Find("a").Eq(0).Text()
		rawTime := td.Eq(0).Text() + " " + td.Eq(1).Text()
		time := parseEPostDateTime(rawTime)
		state := parseEPostStatus(descriptionText)

		progress := model.DeliveryProgress{
			Description: description,
			Location:    location,
			Time:        time,
			State:       state,
		}
		// 역순으로 추가 (최신 -> 오래된 순으로)
		progresses = append([]model.DeliveryProgress{progress}, progresses...)
	})

	state := model.DeliveryState{ID: "in_transit", Name: "상품이동중"}
	if len(progresses) > 0 {
		state = progresses[0].State
	}

	fromHtml, err := infoTds.Eq(0).Html()
	if err != nil {
		return nil, err
	}

	toHtml, err := infoTds.Eq(1).Html()
	if err != nil {
		return nil, err
	}

	fromParts := strings.Split(strings.TrimSpace(fromHtml), "<br/>")
	toParts := strings.Split(strings.TrimSpace(toHtml), "<br/>")

	fromName := parseLocationName(strings.TrimSpace(fromParts[0]))
	fromTime := ""
	if len(progresses) > 1 {
		fromTime = progresses[len(progresses)-1].Time
	}

	toName := parseLocationName(strings.TrimSpace(toParts[0]))
	toTime := ""
	if state.Name == "배달완료" && len(progresses) > 0 {
		toTime = progresses[0].Time
	}

	return &model.DeliveryResult{
		From: model.DeliveryLocation{
			Name: fromName,
			Time: fromTime,
		},
		To: model.DeliveryLocation{
			Name: toName,
			Time: toTime,
		},
		Progresses: progresses,
		State:      state,
	}, nil
}
