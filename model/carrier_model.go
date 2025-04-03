package model

type Carrier struct {
	ID             int    `json:"id"`
	UID            string `json:"uid"`
	No             int    `json:"no"`
	Name           string `json:"name"`
	DisplayName    string `json:"displayName"`
	IsCrawlable    bool   `json:"isCrawlable"`
	IsPopupEnabled bool   `json:"isPopupEnabled"`
	PopupURL       string `json:"popupURL"`
}
