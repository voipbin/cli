package output

type YAMLFormatter struct{}

func (f *YAMLFormatter) FormatList(data interface{}, columns []Column) error {
	return printYAML(data)
}

func (f *YAMLFormatter) FormatItem(data interface{}, columns []Column) error {
	return printYAML(data)
}
