package exportlistmapping

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

func NewSetting(p string) (*Setting, error) {
	s := new(Setting)
	if err := sjson.OpenDecode(p, s); err != nil {
		return nil, err
	}
	return s, nil
}
