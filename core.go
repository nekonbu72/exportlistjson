package exportlistjson

import "time"

type Data struct {
	Filename string
	Date     time.Time
	Invoice  string
	Kata     string
	Lot      string
	Qty      int
}
