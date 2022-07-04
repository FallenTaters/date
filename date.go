package date

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

// ErrScan is returned when date.Date.Scan() fails
var ErrScan = errors.New(`scan failed`)

// Format is the format used for all conversions
const Format = `2006-01-02`

/*
Date represents a timezone-agnostic date.

Two Dates can be compared using == and will be equal if the represented dates are equal.

Date implements encoding.TextMarshaler and encoding.TextUnmarshaler, making it compatible with
encoding/* packages.

Date implements sql.Scaner and driver.Valuer, making it compatible with database/sql.

It only supports the ISO8601 date format `2006-01-02`, exposed as date.Format
*/
type Date struct {
	t time.Time
}

// New makes a new date with the specified year, month and day.
// The values can be out of range and will be handled by time.Date().
func New(year int, month time.Month, day int) Date {
	return From(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

// From removes the time component of t and returns a Date.
func From(t time.Time) Date {
	return Date{time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)}
}

// Time casts d to a time.Time.
func (d Date) Time() time.Time {
	return d.t
}

// String implements fmt.Stringer.
func (d Date) String() string {
	return d.t.Format(Format)
}

// GoString implements fmt.GoStringer.
func (d Date) GoString() string {
	return fmt.Sprintf(`date.New(%d, %d, %d)`, d.Time().Year(), d.Time().Month(), d.Time().Day())
}

// Scan implements database/sql.Scanner.
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
		return fmt.Errorf(`%w: cannot unmarshal %v into date.Date: %s`, ErrScan, src, err.Error())
	}

	return nil
}

// Value implements database/sql/driver.Valuer.
func (d Date) Value() (driver.Value, error) {
	return d.t, nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (d *Date) UnmarshalText(data []byte) error {
	t, err := time.Parse(Format, string(data)) // will always be UTC
	if err == nil {
		*d = Date{t}
	}

	return err
}

// MarshalText implements encoding.TextMarshaler
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}
