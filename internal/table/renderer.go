package table

import "github.com/pterm/pterm"

// TableRenderer 表格渲染器
type TableRenderer struct {
	headers []string
	data    [][]string
}

// NewTableRenderer 创建表格渲染器
func NewTableRenderer(headers []string) *TableRenderer {
	return &TableRenderer{
		headers: headers,
		data:    make([][]string, 0),
	}
}

// AddRow 添加一行数据
func (t *TableRenderer) AddRow(row []string) error {
	t.data = append(t.data, row)
	return nil
}

// MustAddRow 添加一行数据，如果出错则 panic
func (t *TableRenderer) MustAddRow(row []string) {
	_ = t.AddRow(row)
}

// Render 渲染表格
func (t *TableRenderer) Render() error {
	tableData := pterm.TableData{t.headers}
	tableData = append(tableData, t.data...)
	return pterm.DefaultTable.
		WithHasHeader().
		WithBoxed(true).
		WithData(tableData).
		Render()
}
