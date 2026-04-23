package tool

import "time"

func ParseStringToTime(t string) (time.Time, error) {
	deadline, err := time.Parse("2006-01-02 15:04:05", t)
	if err != nil {
		return time.Time{}, err
	}
	return deadline, nil
}

func ParseTimeToString(t time.Time) string {
	return t.Format("2006-01-02")
}
