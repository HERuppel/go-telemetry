package utils

import "time"

func GetUnixStartAndEndOfDay(dateString string) (int64, int64, error) {
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return 0, 0, err
	}

	start := date.Unix()

	end := date.Add(24*time.Hour - time.Second).Unix()

	return start, end, nil
}
