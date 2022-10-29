package utils

import (
	"strings"
	"time"
)

type DateTime struct {
	time.Time
}

func (t *DateTime) MarshalJSON() ([]byte, error) {
	stamp := time.Now().Format(time.RFC3339)
	return []byte("\"" + stamp + "\""), nil
}

func (t *DateTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	date, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	t.Time = date
	return
}
