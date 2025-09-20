package find

type TableView struct {
	table TextTable
	ref   []int
}

func NewTableView(table TextTable) *TableView {
	ref := everyElement(table.GetRowCount())
	return &TableView{
		table: table,
		ref:   ref,
	}
}

func (v *TableView) UpdateFind(rawPattern string) {
	v.ref = Find(rawPattern, v.table)
}

func (v *TableView) GetColumnCount() int {
	return v.table.GetColumnCount()
}

func (v *TableView) GetColumn(col int) string {
	return v.table.GetColumn(col)
}

func (v *TableView) GetRowCount() int {
	return len(v.ref)
}
func (v *TableView) GetCell(row int, col int) string {
	return v.table.GetCell(v.ref[row], col)
}
