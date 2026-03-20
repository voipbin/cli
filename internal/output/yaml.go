package output

import "io"

type YAMLFormatter struct {
	Writer io.Writer
}

func (f *YAMLFormatter) FormatList(data interface{}, columns []Column) error {
	return printYAML(f.Writer, data)
}

func (f *YAMLFormatter) FormatItem(data interface{}, columns []Column) error {
	return printYAML(f.Writer, data)
}
