package find

type TableView struct {
	table TextTable
}

func NewTableView(table TextTable) *TableView {
	return &TableView{
		table: table,
	}
}

