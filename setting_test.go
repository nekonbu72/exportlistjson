package exportlistjson_test

import (
	"testing"

	"github.com/nekonbu72/exportlistjson"
)

func TestNewSetting(t *testing.T) {
	s, err := exportlistjson.NewSetting()
	if err != nil {
		t.Errorf("NewSetting: %v\n", err)
		return
	}

	if s.Sheet != "Sheet1" {
		t.Errorf("Sheet: %v\n", s.Sheet)
		return
	}

	if s.Columns.Kata != 3 {
		t.Errorf("Columns.Kata: %v\n", s.Columns.Kata)
		return
	}

	if s.Rows.Start != 8 {
		t.Errorf("Rows.Start: %v\n", s.Rows.Start)
		return
	}

	if s.Cells.Date.Row != 0 {
		t.Errorf("Cells.Date.Row: %v\n", s.Cells.Date.Row)
		return
	}
}
