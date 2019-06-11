package exportlistmapping

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/nekonbu72/xemlsx"
	"github.com/tealeg/xlsx"
)

const (
	errLimit = 3
)

type Data struct {
	FileName string
	Date     string
	Invoice  string
	Kata     string
	Lot      string
	Qty      int
}

type XLSXData struct {
	FileName string
	*SheetData
}

type SheetData struct {
	Date    string
	Invoice string
	RowData []*RowData
}

type RowData struct {
	Kata string
	Lot  string
	Qty  int
}

func Fetch(
	setting *Setting,
	xlsxStream <-chan *xemlsx.XLSX,
) ([]*Data, error) {
	done := make(chan interface{})
	defer close(done)

	xlsxDataStream := toXLSXData(done, setting, xlsxStream)
	var data []*Data
	for d := range toData(done, xlsxDataStream) {
		data = append(data, d)
	}
	return data, nil
}

func ToData(
	done <-chan interface{},
	setting *Setting,
	xlsxStream <-chan *xemlsx.XLSX,
) <-chan *Data {
	xlsxDataStream := toXLSXData(done, setting, xlsxStream)
	return toData(done, xlsxDataStream)
}

func toData(
	done <-chan interface{},
	xlsxDataStream <-chan *XLSXData,
) <-chan *Data {
	dataStream := make(chan *Data)
	go func() {
		defer close(dataStream)

		for xd := range xlsxDataStream {
			sd := xd.SheetData
			for _, rd := range sd.RowData {
				select {
				case <-done:
					return
				case dataStream <- &Data{
					Date:     sd.Date,
					FileName: xd.FileName,
					Invoice:  sd.Invoice,
					Kata:     rd.Kata,
					Lot:      rd.Lot,
					Qty:      rd.Qty,
				}:
				}
			}
		}
	}()
	return dataStream
}

func toXLSXData(
	done <-chan interface{},
	setting *Setting,
	xlsxStream <-chan *xemlsx.XLSX,
) <-chan *XLSXData {
	xlsxDataStream := make(chan *XLSXData)
	go func() {
		defer close(xlsxDataStream)

		errCount := 0
		for x := range xlsxStream {
			xd, err := xlsxData(setting, x)
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
			case xlsxDataStream <- xd:
			}
		}
	}()
	return xlsxDataStream
}

func xlsxData(setting *Setting, x *xemlsx.XLSX) (*XLSXData, error) {
	sheet, ok := x.Sheet[setting.Sheet]
	if ok == false {
		return nil, errors.New("Sheet not found")
	}

	sd, err := sheetData(setting, sheet)
	if err != nil {
		return nil, err
	}

	return &XLSXData{
		FileName:  x.FileName,
		SheetData: sd,
	}, nil
}

func sheetData(setting *Setting, sheet *xlsx.Sheet) (*SheetData, error) {
	done := make(chan interface{})
	defer close(done)
	rowStream := generateRow(done, setting, sheet)
	rowDataStream := toRowData(done, setting, rowStream)

	date := strings.Trim(
		sheet.Cell(
			setting.Date.Row,
			setting.Date.Column,
		).Value, setting.Date.Remove)
	if date == "" {
		return nil, errors.New("Empty date")
	}

	// strings.Trim だとうまくいかなかった
	invoice := strings.Replace(
		sheet.Cell(
			setting.Invoice.Row,
			setting.Invoice.Column,
		).Value, setting.Invoice.Remove, "", 1)
	if invoice == "" {
		return nil, errors.New("Empty invoice")
	}

	var rds []*RowData
	for rd := range rowDataStream {
		rds = append(rds, rd)
	}
	return &SheetData{
		Date:    date,
		Invoice: invoice,
		RowData: rds,
	}, nil
}

func generateRow(
	done <-chan interface{},
	setting *Setting,
	sheet *xlsx.Sheet,
) <-chan *xlsx.Row {
	rowStearm := make(chan *xlsx.Row)
	go func() {
		defer close(rowStearm)

		for r := setting.Start; r <= sheet.MaxRow; r++ {
			select {
			case <-done:
				return
			case rowStearm <- sheet.Row(r):
			}
		}
	}()
	return rowStearm
}

func toRowData(
	done <-chan interface{},
	setting *Setting,
	rowStream <-chan *xlsx.Row,
) <-chan *RowData {
	rowDataStream := make(chan *RowData)
	go func() {
		defer close(rowDataStream)

		for r := range rowStream {
			rd, err := rowData(setting, r)
			if err != nil {
				break
			}

			select {
			case <-done:
				return
			case rowDataStream <- rd:
			}
		}
	}()
	return rowDataStream
}

func rowData(setting *Setting, row *xlsx.Row) (*RowData, error) {
	kata := row.Cells[setting.Kata].Value
	if kata == "" {
		return nil, errors.New("Empty kata")
	}

	lot := row.Cells[setting.Lot].Value
	if lot == "" {
		return nil, errors.New("Empty lot")
	}

	qty, err := strconv.Atoi(row.Cells[setting.Qty].Value)
	if err != nil {
		return nil, err
	}
	if qty <= 0 {
		return nil, errors.New("Zero or minus qty")
	}

	return &RowData{
		Kata: kata,
		Lot:  lot,
		Qty:  qty,
	}, nil
}
