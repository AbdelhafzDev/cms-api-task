package timeutil

import (
	"encoding/json"
	"time"
)

//  RFC3339 format with millisecond precision
const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

type Milli time.Time


func (t Milli) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(RFC3339Milli))
}


func (t *Milli) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		parsed, err = time.Parse(RFC3339Milli, s)
		if err != nil {
			return err
		}
	}
	*t = Milli(parsed)
	return nil
}


func (t Milli) Time() time.Time {
	return time.Time(t)
}


func ToMilli(t time.Time) Milli {
	return Milli(t)
}


func ToMilliPtr(t *time.Time) *Milli {
	if t == nil {
		return nil
	}
	m := Milli(*t)
	return &m
}
