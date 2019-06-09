package exportlistjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/nekonbu72/xemlsx"
)

type Data struct {
	FileName string
	Date     string
	Invoice  string
	Kata     string
	Lot      string
	Qty      int
}

const (
	errLimit = 3
)

func toJSON(s *Setting, x *xemlsx.XLSX) (string, error) {
	ds, err := toData(s, x)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	b, err := json.Marshal(ds)
	if err != nil {
		return "", err
	}

	if _, err := buf.Write(b); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func toData(s *Setting, x *xemlsx.XLSX) ([]*Data, error) {
	var datas []*Data

	sheet, ok := x.File.Sheet[s.Sheet]
	if !ok {
		return nil, errors.New("Sheet not found: " + s.Sheet)
	}

	date := strings.Trim(
		sheet.Cell(
			s.Date.Row,
			s.Date.Column,
		).Value, s.Date.Remove)
	if date == "" {
		return nil, errors.New("Empty date")
	}

	// strings.Trim だとうまくいかなかった
	inv := strings.Replace(
		sheet.Cell(
			s.Invoice.Row,
			s.Invoice.Column,
		).Value, s.Invoice.Remove, "", 1)
	if inv == "" {
		return nil, errors.New("Empty invoice")
	}

	for r := s.Start; r <= sheet.MaxRow; r++ {
		kata := sheet.Cell(r, s.Kata).Value
		if kata == "" {
			r++
			continue
		}

		lot := sheet.Cell(r, s.Lot).Value
		if lot == "" {
			r++
			continue
		}

		qty, err := strconv.Atoi(sheet.Cell(r, s.Qty).Value)
		if err != nil {
			r++
			continue
		}
		if qty <= 0 {
			r++
			continue
		}

		data := &Data{
			FileName: x.FileName,
			Date:     date,
			Invoice:  inv,
			Kata:     kata,
			Lot:      lot,
			Qty:      qty,
		}

		datas = append(datas, data)
	}
	return datas, nil
}

func ToJSON(
	done <-chan interface{},
	xlsxStream <-chan *xemlsx.XLSX,
) <-chan string {

	jsonStream := make(chan string)
	go func() {
		defer close(jsonStream)

		setting, err := NewSetting()
		if err != nil {
			log.Printf("NewSetting: %v\n", err)
			return
		}
		errCount := 0
		for x := range xlsxStream {
			s, err := toJSON(setting, x)
			if err != nil {
				log.Printf("error: %v", err)
				errCount++
				if errCount >= errLimit {
					log.Println("Too many errors, breaking!")
					break
				}
				continue
			}

			select {
			case <-done:
				return
			case jsonStream <- s:
			}
		}
	}()
	return jsonStream
}
