package date_test

import (
	"testing"
	"time"

	"github.com/FallenTaters/date"
)

var times = []struct {
	name string
	t    time.Time
}{
	{`zero time`, time.Time{}},
	{`negative time (out of range values)`, time.Date(-1, -1, -1, 0, 0, 0, 0, time.UTC)},
	{`far future time`, time.Date(100_000, time.December, 25, 0, 0, 0, 0, time.UTC)},
	{`time included`, time.Date(2020, time.February, 30, 12, 13, 14, 151617, time.Local)},
	{`now`, time.Now()},
}

func TestNew(t *testing.T) {
	t.Run(`normal date`, func(t *testing.T) {
		if d := date.New(2020, time.December, 13); d.String() != `2020-12-13` {
			t.Error(d)
		}
	})

	t.Run(`values out of range -> should be normalized`, func(t *testing.T) {
		if d := date.New(-1, -1, -1); d.String() != `-0002-10-30` {
			t.Error(d)
		}
	})
}

func TestFrom(t *testing.T) {
	for _, c := range times {
		t.Run(c.name, func(t *testing.T) {
			if d := date.From(c.t); d.String() != c.t.Format(date.Format) {
				t.Errorf(`date %s should be %q but is %q`, c.name, d.String(), c.t.Format(date.Format))
			}
		})
	}
}

func TestGoString(t *testing.T) {
	expected := `date.New(1, 1, 1)`
	if actual := date.From(times[0].t).GoString(); actual != expected {
		t.Errorf(`expected %q, but got %q`, expected, actual)
	}

	expected = `date.New(-2, 10, 30)`
	if actual := date.From(times[1].t).GoString(); actual != expected {
		t.Errorf(`expected %q, but got %q`, expected, actual)
	}

	expected = `date.New(100000, 12, 25)`
	if actual := date.From(times[2].t).GoString(); actual != expected {
		t.Errorf(`expected %q, but got %q`, expected, actual)
	}

	expected = `date.New(2020, 3, 1)`
	if actual := date.From(times[3].t).GoString(); actual != expected {
		t.Errorf(`expected %q, but got %q`, expected, actual)
	}
}

func TestComparable(t *testing.T) {
	actual, expected := date.From(times[3].t), date.New(2020, 3, 1)
	if expected != actual {
		t.Error(actual, expected)
	}
}

func TestScan(t *testing.T) {
	cases := []struct {
		name      string
		input     any
		expectErr string
		expected  date.Date
	}{
		{
			name:     `string date`,
			input:    `2006-10-12`,
			expected: date.New(2006, 10, 12),
		},
		{
			name:     `byte slice`,
			input:    []byte(`2006-10-12`),
			expected: date.New(2006, 10, 12),
		},
		{
			name:     `time.Time`,
			input:    time.Date(2006, 10, 12, 1, 2, 3, 4, time.Local),
			expected: date.New(2006, 10, 12),
		},
		{
			name:      `wrong type`,
			input:     123345,
			expectErr: `scan failed: cannot unmarshal variable of type int into date.Date`,
		},
		{
			name:      `bad format`,
			input:     `213-4=vdf&`,
			expectErr: `scan failed: cannot unmarshal 213-4=vdf& into date.Date: parsing time "213-4=vdf&" as "2006-01-02": cannot parse "4=vdf&" as "2006"`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var d date.Date
			err := d.Scan(c.input)
			if c.expectErr != `` {
				if err == nil {
					t.Errorf(`err should be %q, but its nil`, c.expectErr)
				}

				if err != nil && err.Error() != c.expectErr {
					t.Errorf(`err should be %q but it is %q`, c.expectErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestValue(t *testing.T) {
	d := date.From(times[3].t)
	actual, err := d.Value()
	if err != nil {
		t.Error(err)
	}
	expected := time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)
	if actual != expected {
		t.Errorf(`should be %#v, but its %#v`, expected, actual)
	}
}

func TestUnmarshalText(t *testing.T) {
	t.Run(`successful unmarshal`, func(t *testing.T) {
		var actual date.Date
		err := actual.UnmarshalText([]byte(`2020-03-01`))
		if err != nil {
			t.Error(err)
		}
		expected := date.New(2020, 3, 1)
		if actual != expected {
			t.Errorf(`should be %#v, but is %#v`, expected, actual)
		}
	})

	t.Run(`failed unmarshal`, func(t *testing.T) {
		var d date.Date
		err := d.UnmarshalText([]byte(`20%0-0.3-01`))
		if err == nil {
			t.Error(`err should be non-nil, but is nil`)
			return
		}
		expected := `parsing time "20%0-0.3-01" as "2006-01-02": cannot parse "-0.3-01" as "2006"`
		if err.Error() != expected {
			t.Errorf(`should be %#v, but is %#v`, expected, err.Error())
		}
	})
}

func TestMarshalText(t *testing.T) {
	d := date.New(2020, 3, 1)
	actual, err := d.MarshalText()
	if err != nil {
		t.Error(err)
	}
	expected := `2020-03-01`
	if string(actual) != expected {
		t.Errorf(`should be %s but is %s`, expected, actual)
	}
}
