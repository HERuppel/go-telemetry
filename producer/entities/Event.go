package entities

type Event struct {
	Type      string  `json:"type"`
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}
