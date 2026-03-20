package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
)

type TableFormatter struct {
	Writer io.Writer
}

func (f *TableFormatter) FormatList(data interface{}, columns []Column) error {
	rows, err := toRows(data, columns)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		fmt.Fprintln(f.Writer, "No items found.")
		return nil
	}

	headers := make([]string, len(columns))
	for i, c := range columns {
		headers[i] = c.Name
	}

	table := tablewriter.NewWriter(f.Writer)
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("  ")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(rows)
	table.Render()
	return nil
}

func (f *TableFormatter) FormatItem(data interface{}, columns []Column) error {
	row, err := toRow(data, columns)
	if err != nil {
		return err
	}

	maxWidth := 0
	for _, c := range columns {
		if len(c.Name) > maxWidth {
			maxWidth = len(c.Name)
		}
	}

	for i, c := range columns {
		fmt.Fprintf(f.Writer, "%-*s  %s\n", maxWidth, c.Name+":", row[i])
	}
	return nil
}

func toRows(data interface{}, columns []Column) ([][]string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var items []map[string]interface{}
	if err := json.Unmarshal(b, &items); err != nil {
		// Try as single item
		var item map[string]interface{}
		if err2 := json.Unmarshal(b, &item); err2 != nil {
			return nil, fmt.Errorf("could not convert data to table: %w", err)
		}
		items = []map[string]interface{}{item}
	}

	rows := make([][]string, 0, len(items))
	for _, item := range items {
		row := make([]string, len(columns))
		for i, c := range columns {
			row[i] = extractField(item, c.Field)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func toRow(data interface{}, columns []Column) ([]string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var item map[string]interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return nil, fmt.Errorf("could not convert data to row: %w", err)
	}

	row := make([]string, len(columns))
	for i, c := range columns {
		row[i] = extractField(item, c.Field)
	}
	return row, nil
}

func extractField(item map[string]interface{}, field string) string {
	val, ok := item[field]
	if !ok || val == nil {
		return ""
	}
	return fmt.Sprintf("%v", val)
}
