package date

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

// Based of https://github.com/go-gorm/datatypes/blob/e8a383d1ba59f52be2c6cd24c25b880ffdb9c64c/date.go

type Date time.Time

func (date *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = Date(nullTime.Time)
	return
}

func (date Date) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date).Location()), nil
}

// GormDataType gorm common data type
func (date Date) GormDataType() string {
	return "date"
}

func (date Date) MarshalJSON() ([]byte, error) {
	t := time.Time(date)
	s := fmt.Sprintf(`"%s"`, t.Format(time.DateOnly))

	return []byte(s), nil
}

func (date *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	// TODO(https://go.dev/issue/47353): Properly unescape a JSON string.
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("Date.UnmarshalJSON: input is not a JSON string")
	}

	data = data[len(`"`) : len(data)-len(`"`)]

	t, err := time.Parse(time.DateOnly, string(data))
	if err != nil {
		return err
	}

	*date = Date(t)

	return nil

}
