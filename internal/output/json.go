package output

type JSONFormatter struct{}

func (f *JSONFormatter) FormatList(data interface{}, columns []Column) error {
	return printJSON(data)
}

func (f *JSONFormatter) FormatItem(data interface{}, columns []Column) error {
	return printJSON(data)
}
