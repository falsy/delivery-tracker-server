package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/falsy/delivery-tracker-server/internal/domain/model"
)

func parseCSRFToken(html string) (string, error) {
	tokenStart := strings.Index(html, `name="_csrf" value="`)
	if tokenStart == -1 {
		return "", errors.New("CSRF token not found")
	}
	tokenStart += len(`name="_csrf" value="`)
	tokenEnd := strings.Index(html[tokenStart:], `"`)
	if tokenEnd == -1 {
		return "", errors.New("invalid csrf token format")
	}
	return html[tokenStart : tokenStart+tokenEnd], nil
}

func parseStatus(code string) model.DeliveryState {
	switch code {
	case "41", "42", "44":
		return model.DeliveryState{ID: "in_transit", Name: "상품이동중"}
	case "11":
		return model.DeliveryState{ID: "item_received", Name: "상품인수"}
	case "82":
		return model.DeliveryState{ID: "out_for_delivery", Name: "배달출발"}
	case "91":
		return model.DeliveryState{ID: "delivered", Name: "배달완료"}
	default:
		return model.DeliveryState{ID: "preparing_item", Name: "상품준비중"}
	}
}

func CJLogisticsGetTrack(trackingNumber string) (*model.DeliveryResult, error) {
	getRes, err := http.Get("https://www.cjlogistics.com/ko/tool/parcel/tracking")
	if err != nil {
		return nil, err
	}
	defer getRes.Body.Close()

	bodyBytes, err := io.ReadAll(getRes.Body)
	if err != nil {
		return nil, err
	}
	csrf, err := parseCSRFToken(string(bodyBytes))
	if err != nil {
		return nil, err
	}

	rawCookies := getRes.Header.Values("Set-Cookie")
	var cookieStrings []string
	for _, c := range rawCookies {
		parts := strings.Split(c, ";")
		if len(parts) > 0 {
			cookieStrings = append(cookieStrings, strings.TrimSpace(parts[0]))
		}
	}
	cookies := strings.Join(cookieStrings, "; ")

	postUrl := fmt.Sprintf("https://www.cjlogistics.com/ko/tool/parcel/tracking-detail?paramInvcNo=%s&_csrf=%s", trackingNumber, csrf)
	req, err := http.NewRequest("POST", postUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", cookies)

	postRes, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer postRes.Body.Close()

	if postRes.StatusCode != 200 {
		return nil, errors.New("운송장 조회에 실패하였습니다")
	}

	var jsonData struct {
		ParcelResultMap struct {
			ResultList []map[string]string `json:"resultList"`
		} `json:"parcelResultMap"`

		ParcelDetailResultMap struct {
			ResultList []map[string]string `json:"resultList"`
		} `json:"parcelDetailResultMap"`
	}

	err = json.NewDecoder(postRes.Body).Decode(&jsonData)
	if err != nil {
		return nil, err
	}

	if len(jsonData.ParcelResultMap.ResultList) == 0 {
		return nil, errors.New("해당 운송장이 존재하지 않거나 조회할 수 없습니다")
	}

	progresses := make([]model.DeliveryProgress, 0, len(jsonData.ParcelDetailResultMap.ResultList))
	for i := len(jsonData.ParcelDetailResultMap.ResultList) - 1; i >= 0; i-- {
		row := jsonData.ParcelDetailResultMap.ResultList[i]
		progresses = append(progresses, model.DeliveryProgress{
			Description: row["crgNm"],
			Location:    row["regBranNm"],
			Time:        row["dTime"],
			State:       parseStatus(row["crgSt"]),
		})
	}

	state := parseStatus("")
	if len(progresses) > 0 {
		state = progresses[0].State
	}

	from := model.DeliveryLocation{
		Name: parseLocationName(jsonData.ParcelResultMap.ResultList[0]["sendrNm"]),
		Time: jsonData.ParcelDetailResultMap.ResultList[0]["dTime"],
	}

	to := model.DeliveryLocation{
		Name: parseLocationName(jsonData.ParcelResultMap.ResultList[0]["rcvrNm"]),
		Time: func() string {
			if state.Name == "배달완료" {
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
