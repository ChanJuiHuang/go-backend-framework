package util

import (
	"strings"
	"time"
)

type TimeForJson time.Time

func (t *TimeForJson) MarshalJSON() ([]byte, error) {
	if time.Time(*t).IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + time.Time(*t).Format(time.RFC3339) + `"`), nil
}

func (t *TimeForJson) UnmarshalJSON(b []byte) error {
	dateTimeString := strings.Trim(string(b), `"`)
	if dateTimeString == "" || dateTimeString == "null" {
		return nil
	}

	dateTime, err := time.Parse(time.RFC3339, dateTimeString)
	if err != nil {
		return err
	}
	*t = TimeForJson(dateTime)

	return nil
}
