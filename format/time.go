package format

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

// ErrScan is returned when format.Time.Scan() fails
var ErrScan = errors.New(`scan failed`)

// ISO8601Date accepts a date in the format `2006-01-02`
type ISO8601Date struct{}

// TimeFormat implements format.TimeFormat
func (ISO8601Date) TimeFormat() string {
	return `2006-01-02`
}

// RFC3339 is a format.TimeFormat for RFC3339
type RFC3339 struct{}

func (RFC3339) TimeFormat() string {
	return time.RFC3339
}

type TimeFormat interface {
	TimeFormat() string
}

type Time[F TimeFormat] time.Time

func (t Time[F]) format() string {
	var f F
	return (f).TimeFormat()
}

func TimeFrom[T TimeFormat](t time.Time) Time[T] {
	return Time[T](t)
}

func (t Time[F]) Time() time.Time {
	return time.Time(t)
}

func (t Time[F]) String() string {
	return time.Time(t).Format(t.format())
}

func (t *Time[F]) Scan(src any) error {
	var err error

	switch v := src.(type) {
	case []byte:
		err = t.UnmarshalText(v)
	case string:
		err = t.UnmarshalText([]byte(v))
	case time.Time:
		*t = TimeFrom[F](v)
	default:
		return fmt.Errorf(`%w: cannot unmarshal variable of type %T into format.Time`, ErrScan, src)
	}

	if err != nil {
		return fmt.Errorf(`%w: unable to unmarshal %v into format.Time[%T]: %s`, ErrScan, src, *new(F), err.Error())
	}

	return nil
}

func (t Time[F]) Value() (driver.Value, error) {
	return t.Time(), nil
}

func (t *Time[F]) UnmarshalText(data []byte) error {
	tim, err := time.Parse(t.format(), string(data))
	if err == nil {
		*t = Time[F](tim)
	}

	return err
}

func (t *Time[F]) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}
