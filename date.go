package date

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

// ErrScan is returned when date.Date.Scan() fails
var ErrScan = errors.New(`scan failed`)

const Layout = `2006-01-02`

/*
Date represents a timezone-agnostic date.

Two Dates can be compared and will be equal if the represented dates are equal.

Date implements encoding.TextMarshaler and encoding.TestUnmarshaler, making it compatible with
encoding/json, encoding/xml, etc.

It only supports the ISO8601 date format `2006-01-02`
*/
type Date time.Time

func New(year int, month time.Month, day int) Date {
	return From(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

func From(t time.Time) Date {
	return Date(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC))
}

func (d Date) Time() time.Time {
	return time.Time(d)
}

func (d Date) String() string {
	return d.Time().Format(Layout)
}

func (d *Date) Scan(src any) error {
	var err error
	switch v := src.(type) {
	case []byte:
		err = d.UnmarshalText(v)
	case string:
		err = d.UnmarshalText([]byte(v))
	case time.Time:
		*d = From(v)
	default:
		return fmt.Errorf(`%w: cannot unmarshal variable of type %T into date.Date`, ErrScan, src)
	}

	if err != nil {
		return fmt.Errorf(`%w: unable to unmarshal %v into date.Date: %s`, ErrScan, src, err.Error())
	}

	return nil
}

func (d Date) Value() (driver.Value, error) {
	return time.Time(d), nil
}

func (d *Date) UnmarshalText(data []byte) error {
	t, err := time.Parse(Layout, string(data))
	if err == nil {
		*d = Date(t)
	}

	return err
}

func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.Time().Format(`2006-01-02`)), nil
}
