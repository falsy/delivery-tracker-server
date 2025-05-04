package model

type DeliveryLocation struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

type DeliveryState struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DeliveryProgress struct {
	Description string        `json:"description"`
	Location    string        `json:"location"`
	Time        string        `json:"time"`
	State       DeliveryState `json:"state"`
}

type DeliveryResult struct {
	From       DeliveryLocation   `json:"from"`
	To         DeliveryLocation   `json:"to"`
	State      DeliveryState      `json:"state"`
	Progresses []DeliveryProgress `json:"progresses"`
}
