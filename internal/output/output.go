package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Column struct {
	Name  string
	Field string
	Width int
}

type Formatter interface {
	FormatList(data interface{}, columns []Column) error
	FormatItem(data interface{}, columns []Column) error
}

func NewFormatter(format string, w io.Writer) (Formatter, error) {
	switch format {
	case "table", "":
		return &TableFormatter{Writer: w}, nil
	case "json":
		return &JSONFormatter{Writer: w}, nil
	case "yaml":
		return &YAMLFormatter{Writer: w}, nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}

func GetFormat(cmd *cobra.Command) string {
	f, _ := cmd.Flags().GetString("output")
	return f
}

func PrintList(cmd *cobra.Command, data interface{}, columns []Column) error {
	format := GetFormat(cmd)
	w := cmd.OutOrStdout()
	formatter, err := NewFormatter(format, w)
	if err != nil {
		return err
	}
	return formatter.FormatList(data, columns)
}

func PrintItem(cmd *cobra.Command, data interface{}, columns []Column) error {
	format := GetFormat(cmd)
	w := cmd.OutOrStdout()
	formatter, err := NewFormatter(format, w)
	if err != nil {
		return err
	}
	return formatter.FormatItem(data, columns)
}

func marshalJSON(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func printJSON(w io.Writer, data interface{}) error {
	b, err := marshalJSON(data)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %w", err)
	}
	fmt.Fprintln(w, string(b))
	return nil
}

func printYAML(w io.Writer, data interface{}) error {
	// Convert through JSON to handle pointer fields properly
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("could not marshal data: %w", err)
	}
	var generic interface{}
	if err := json.Unmarshal(b, &generic); err != nil {
		return fmt.Errorf("could not process data: %w", err)
	}
	out, err := yaml.Marshal(generic)
	if err != nil {
		return fmt.Errorf("could not marshal YAML: %w", err)
	}
	fmt.Fprint(w, string(out))
	return nil
}
