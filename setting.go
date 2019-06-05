package exportlistjson

import "github.com/nekonbu72/sjson/sjson"

type Setting struct {
	Sheet   string `json:"sheet"`
	Columns struct {
		Kata int `json:"Kata"`
		Lot  int `json:"Lot"`
		Qty  int `json:"Qty"`
	} `json:"columns"`
	Rows struct {
		Start int `json:"start"`
	} `json:"rows"`
	Cells struct {
		Date struct {
			Row    int `json:"row"`
			Column int `json:"column"`
		} `json:"date"`
		Invoice struct {
			Row    int `json:"row"`
			Column int `json:"column"`
		} `json:"invoice"`
	} `json:"cells"`
}

const jsonPath = "setting.json"

func NewSetting() (*Setting, error) {
	s := new(Setting)
	if err := sjson.OpenDecode(jsonPath, s); err != nil {
		return nil, err
	}
	return s, nil
}
