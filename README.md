# Date

A fully comparable date object that always formats to ISO8601 format YYYY-MM-DD. Meant to be used for JSON or other encoding.

## Example

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/FallenTaters/date"
)

type Payload struct {
	Date date.Date `json:"date"`
}

func main() {
	d := date.New(2022, 7, 4)
	fmt.Println(d) // 2022-07-04

	// json
	jsonText := `{"date":"2022-07-04"}`
	var p Payload
	err := json.Unmarshal([]byte(jsonText), &p)
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data) == jsonText) // true

	// comparable
	d1 := time.Date(2022, 7, 4, 0, 0, 0, 0, time.UTC)
	d2 := d1.Add(time.Hour + time.Minute + time.Second)
	fmt.Println(date.From(d1) == date.From(d2)) // true
}
```

## Description

* Compatible with
    * `encoding/*`
        * through `encoding.TextUnmarshaler` and `encoding.TextMarshaler`
    * `database/sql` and `database/sql/driver`
        * through `sql.Scanner` and `driver.Valuer`
* Comparable
    * hour, minute, second, nanoseconds and timezone are stripped when using `From`

## Notes

* For more formats, see https://github.com/FallenTaters/timefmt
