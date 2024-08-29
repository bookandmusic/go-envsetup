package utils

import (
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
)

// TableConfig 结构体包含生成表格所需的配置信息
type TableConfig struct {
	Header table.Row
	Data   []table.Row
}

func RenderTable(config *TableConfig, mirror io.Writer) error {
	t := table.NewWriter()
	t.SetOutputMirror(mirror)
	t.AppendHeader(config.Header)
	t.AppendRows(config.Data)
	t.SetStyle(table.StyleColoredBright)
	t.Render()
	return nil
}
