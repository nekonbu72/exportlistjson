package exportlistmapping

import (
	"errors"

	"github.com/nekonbu72/sjson/sjson"
	"github.com/tealeg/xlsx"
)

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

func (s *Setting) isValid(sheet *xlsx.Sheet) error {
	if isWithinMaxRow(s.Date.Row, sheet) == false {
		return errors.New("Data.Row")
	}

	if isWithinMaxRow(s.Invoice.Row, sheet) == false {
		return errors.New("Invoice.Row")
	}

	if isWithinMaxCol(s.Kata, sheet) == false {
		return errors.New("Kata.Col")
	}

	if isWithinMaxCol(s.Lot, sheet) == false {
		return errors.New("Lot.Col")
	}

	if isWithinMaxCol(s.Qty, sheet) == false {
		return errors.New("Qty.Col")
	}

	return nil
}

func isWithinMaxRow(row int, sheet *xlsx.Sheet) bool {
	return row <= sheet.MaxRow
}

func isWithinMaxCol(col int, sheet *xlsx.Sheet) bool {
	return col <= sheet.MaxCol
}
