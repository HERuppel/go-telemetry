package entities

type Metrics struct {
	EventType    string  `json:"eventType"`
	Count        int64   `json:"count"`
	AverageValue float64 `json:"averageValue"`
}

type MetricsSinceDayOne struct {
	EventType    string  `json:"eventType"`
	Count        int64   `json:"count"`
	Sum          float64 `json:"sum"`
	AverageValue float64 `json:"averageValue"`
}
