package date_test

import (
	"encoding/json"
	"encoding/xml"
	"testing"
	"time"

	"github.com/FallenTaters/date"
)

var times = []time.Time{
	{},
	time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
	time.Date(-1, 0, 0, 0, 0, 0, 0, time.UTC),
	time.Date(0, -1, 0, 0, 0, 0, 0, time.UTC),
	time.Date(0, 0, -1, 0, 0, 0, 0, time.UTC),
	time.Now(),
}

func TestDate(t *testing.T) {
	emptyTimeStr := (time.Time{}).Format(date.Layout)
	emptyDateStr := (date.Date{}).String()
	if emptyDateStr != emptyTimeStr {
		t.Errorf(`%s should be %s`, emptyDateStr, emptyTimeStr)
	}

	for _, tim := range times {
		expected := tim.Format(date.Layout)
		actual := date.From(tim).String()
		if expected != actual {
			t.Errorf(`expected %s, got %s`, expected, actual)
		}
	}

	t.Run(`json`, func(t *testing.T) {
		input := `"1996-02-14"`
		var d date.Date
		err := json.Unmarshal([]byte(input), &d)
		if err != nil {
			t.Error(err)
		}

		data, err := json.Marshal(d)
		if err != nil {
			t.Error(err)
		}
		if string(data) != input {
			t.Errorf(`%q should be "1996-02-14"`, data)
		}
	})

	t.Run(`xml`, func(t *testing.T) {
		input := `<xml><date>1996-02-14</date></xml>`
		dst := struct {
			XMLName xml.Name  `xml:"xml"`
			Date    date.Date `xml:"date"`
		}{}

		err := xml.Unmarshal([]byte(input), &dst)
		if err != nil {
			t.Error(err)
		}

		dateStr := dst.Date.String()
		if dateStr != `1996-02-14` {
			t.Error(dateStr)
		}

		data, err := xml.Marshal(dst)
		if err != nil {
			t.Error(err)
		}

		if string(data) != input {
			t.Error(string(data))
		}
	})

	t.Run(`sql`, func(t *testing.T) {
		testTime := time.Date(1996, 2, 14, 0, 0, 0, 0, time.UTC)
		expected := date.From(testTime)

		v, err := expected.Value()
		if err != nil {
			t.Error(err)
		}
		if v != testTime {
			t.Error(v)
		}

		for _, v := range [...]any{`1996-02-14`, []byte(`1996-02-14`), testTime} {
			var actual date.Date
			err = actual.Scan(v)
			if err != nil {
				t.Error(err)
			}
			if actual != expected {
				t.Error(actual, expected)
			}
		}

		var actual date.Date
		err = actual.Scan(true)
		if err == nil {
			t.Error(`should not be nil`)
			t.FailNow()
		}
		if err.Error() != `scan failed: cannot unmarshal variable of type bool into date.Date` {
			t.Error(err.Error())
		}
	})
}
