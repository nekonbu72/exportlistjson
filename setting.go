package exportlistjson

import "github.com/nekonbu72/sjson/sjson"

type Setting struct {
	Sheet string `json:"sheet"`

	Start int `json:"start"`

	Date struct {
		Remove string `json:"remove"`
		Row    int    `json:"row"`
		Column int    `json:"column"`
	} `json:"date"`

	Invoice struct {
		Remove string `json:"remove"`
		Row    int    `json:"row"`
		Column int    `json:"column"`
	} `json:"invoice"`

	Kata int `json:"kata"`

	Lot int `json:"lot"`

	Qty int `json:"qty"`
}

const jsonPath = "setting.json"

func NewSetting() (*Setting, error) {
	s := new(Setting)
	if err := sjson.OpenDecode(jsonPath, s); err != nil {
		return nil, err
	}
	return s, nil
}
