package utils

import (
	"errors"
	"fmt"
	"time"
)

type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", d.Format(time.DateOnly))
	return []byte(stamp), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("Date.UnmarshalJSON: input is not a JSON string")
	}
	data = data[len(`"`) : len(data)-len(`"`)]

	t, err := time.Parse(time.DateOnly, string(data))
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}
