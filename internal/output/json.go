package output

import "io"

type JSONFormatter struct {
	Writer io.Writer
}

func (f *JSONFormatter) FormatList(data interface{}, columns []Column) error {
	return printJSON(f.Writer, data)
}

func (f *JSONFormatter) FormatItem(data interface{}, columns []Column) error {
	return printJSON(f.Writer, data)
}
