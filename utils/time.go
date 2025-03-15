package utils

import (
	"time"
)

type Time time.Time

func (t Time) Unix() int64 {
	return time.Time(t).Unix()
}

var (
	Now     = timeNow
	NowUnix = timeNowUnix
	NowISO  = timeNowISO
)

func timeNow() Time {
	return Time(time.Now().UTC())
}

func timeNowUnix() int64 {
	return time.Time(timeNow()).Unix()
}

func timeNowISO() string {
	return time.Time(timeNow()).Format(time.RFC3339)
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format(time.RFC3339) + `"`), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	parsed, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return err
	}
	time := Time(parsed)
	t = &time
	return nil
}
