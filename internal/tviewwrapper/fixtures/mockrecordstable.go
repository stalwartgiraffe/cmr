package fixtures

import (
	tw "github.com/stalwartgiraffe/cmr/internal/tviewwrapper"
)

type MockRecordsTable struct {
}

type FakeRecord struct {
	Row int
}

func (m *MockRecordsTable) GetRowRecord(row int) any {
	return &FakeRecord{
		Row: row,
	}
}

type MockOnCellSelectedSrc struct {
}

func (m *MockOnCellSelectedSrc) OnCellSelectedSubscribe(func(tw.CellParams)) {
}
