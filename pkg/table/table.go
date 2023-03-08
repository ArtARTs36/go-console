package table

type TableInterface interface {
	SetHeaders(data *TableData) *Table
	GetHeaders() *TableData

	SetRows(data *TableData) *Table
	GetRows() *TableData

	SetHeaderTitle(title string) *Table
	GetHeaderTitle() string

	SetFooterTitle(title string) *Table
	GetFooterTitle() string

	GetLinesAsList() []TableRowInterface
	GetColumnsAsList() []TableColumnInterface
	GetCellsAsList() []TableCellInterface
}

type Table struct {
	headers *TableData
	rows    *TableData

	headerTitle string
	footerTitle string

	columnsPadding map[int]PaddingType
}

// Table constructors

func NewTable() *Table {
	return &Table{
		headers:        NewTableData(),
		rows:           NewTableData(),
		columnsPadding: map[int]PaddingType{},
	}
}

// Implements TableInterface.

var _ TableInterface = (*Table)(nil)

func (t *Table) SetHeaders(data *TableData) *Table {
	t.headers = data
	return t
}

func (t *Table) GetHeaders() *TableData {
	return t.headers
}

func (t *Table) SetRows(data *TableData) *Table {
	t.rows = data
	return t
}

func (t *Table) GetRows() *TableData {
	return t.rows
}

func (t *Table) GetLinesAsList() []TableRowInterface {
	lines := []TableRowInterface{}

	for _, line := range t.headers.GetRows() {
		lines = append(lines, line)
	}

	for _, line := range t.rows.GetRows() {
		lines = append(lines, line)
	}

	return lines
}

func (t *Table) SetHeaderTitle(title string) *Table {
	t.headerTitle = title
	return t
}

func (t *Table) GetHeaderTitle() string {
	return t.headerTitle
}

func (t *Table) SetFooterTitle(title string) *Table {
	t.footerTitle = title
	return t
}

func (t *Table) GetFooterTitle() string {
	return t.footerTitle
}

// Computations Helpers

func (t *Table) GetColumnsAsList() []TableColumnInterface {
	columns := []TableColumnInterface{}

	columns = append(columns, t.headers.GetColumnsAsList()...)
	columns = append(columns, t.rows.GetColumnsAsList()...)

	return columns
}

func (t *Table) GetCellsAsList() []TableCellInterface {
	cells := []TableCellInterface{}

	cells = append(cells, t.headers.GetCellsAsList()...)
	cells = append(cells, t.rows.GetCellsAsList()...)

	return cells
}

// Data injections Helpers for Headers

func (t *Table) SetHeadersFromString(rows [][]string) *Table {
	data := NewTableData()
	t.SetHeaders(data.setDataFromString(rows))
	return t
}

func (t *Table) AddHeaders(rows []TableRowInterface) *Table {
	for _, row := range rows {
		t.AddHeader(row)
	}

	return t
}

func (t *Table) AddHeadersFromString(rows [][]string) *Table {
	for _, row := range rows {
		t.AddHeaderFromString(row)
	}

	return t
}

func (t *Table) AddHeader(row TableRowInterface) *Table {
	t.headers.AddRow(row)
	return t
}

func (t *Table) AddHeaderFromString(row []string) *Table {
	t.headers.AddRowFromString(row)
	return t
}

func (t *Table) setHeader(column int, row TableRowInterface) *Table {
	t.headers.SetRow(column, row)
	return t
}

func (t *Table) setHeaderFromString(column int, rowData []string) *Table {
	row := MakeRowFromStrings(rowData)
	t.headers.SetRow(column, row)
	return t
}

// Data injections Helpers for Rows

func (t *Table) SetRowsFromString(rows [][]string) *Table {
	data := NewTableData()
	t.SetRows(data.setDataFromString(rows))
	return t
}

func (t *Table) AddRows(rows []TableRowInterface) *Table {
	for _, row := range rows {
		t.AddRow(row)
	}

	return t
}

func (t *Table) AddRowsFromString(rows [][]string) *Table {
	for _, row := range rows {
		t.AddRowFromString(row)
	}

	return t
}

func (t *Table) AddRow(row TableRowInterface) *Table {
	t.rows.AddRow(row)
	return t
}

func (t *Table) AddRowFromString(row []string) *Table {
	t.rows.AddRowFromString(row)
	return t
}

func (t *Table) setRow(column int, row TableRowInterface) *Table {
	t.rows.SetRow(column, row)
	return t
}

func (t *Table) setRowFromString(column int, rowData []string) *Table {
	row := MakeRowFromStrings(rowData)
	t.rows.SetRow(column, row)
	return t
}

func (t *Table) AddTableSeparator() *Table {
	row := NewTableRow().
		AddColumn(
			NewTableColumn().
				SetCell(NewTableSeparator()),
		)

	t.AddRow(row)

	return t
}

// Columns Padding

func (t *Table) SetColumnPadding(column int, padding PaddingType) *Table {
	t.columnsPadding[column-1] = padding
	return t
}

func (t *Table) GetColumnPadding(column int) PaddingType {
	if _, ok := t.columnsPadding[column-1]; ok {
		return PadDefault
	}

	return t.columnsPadding[column]
}