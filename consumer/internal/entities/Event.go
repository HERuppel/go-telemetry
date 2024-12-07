package entities

type Event struct {
	Type      string  `json:"type"`
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type EventsResponse struct {
	Page   int64   `json:"page"`
	Limit  int64   `json:"limit"`
	Count  int64   `json:"count"`
	Events []Event `json:"events"`
}
