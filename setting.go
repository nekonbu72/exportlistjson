package exportlistjson

import "github.com/nekonbu72/sjson/sjson"

type Setting struct {
	Sheet string `json:"sheet"`

	Start int `json:"start"`

	Filename struct {
		Name string `json:"name"`
	} `json:"filename"`

	Date struct {
		Name   string `json:"name"`
		Remove string `json:"remove"`
		Row    int    `json:"row"`
		Column int    `json:"column"`
	} `json:"date"`

	Invoice struct {
		Name   string `json:"name"`
		Remove string `json:"remove"`
		Row    int    `json:"row"`
		Column int    `json:"column"`
	} `json:"invoice"`

	Kata struct {
		Name  string `json:"name"`
		Index int    `json:"index"`
	} `json:"kata"`

	Lot struct {
		Name  string `json:"name"`
		Index int    `json:"index"`
	} `json:"lot"`

	Qty struct {
		Name  string `json:"name"`
		Index int    `json:"index"`
	} `json:"qty"`
}

const jsonPath = "setting.json"

func NewSetting() (*Setting, error) {
	s := new(Setting)
	if err := sjson.OpenDecode(jsonPath, s); err != nil {
		return nil, err
	}
	return s, nil
}
