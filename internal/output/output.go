package output

import (
	"encoding/json"
	"fmt"
	"os"

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

func NewFormatter(format string) (Formatter, error) {
	switch format {
	case "table", "":
		return &TableFormatter{}, nil
	case "json":
		return &JSONFormatter{}, nil
	case "yaml":
		return &YAMLFormatter{}, nil
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
	formatter, err := NewFormatter(format)
	if err != nil {
		return err
	}
	return formatter.FormatList(data, columns)
}

func PrintItem(cmd *cobra.Command, data interface{}, columns []Column) error {
	format := GetFormat(cmd)
	formatter, err := NewFormatter(format)
	if err != nil {
		return err
	}
	return formatter.FormatItem(data, columns)
}

func marshalJSON(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func printJSON(data interface{}) error {
	b, err := marshalJSON(data)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %w", err)
	}
	fmt.Fprintln(os.Stdout, string(b))
	return nil
}

func printYAML(data interface{}) error {
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
	fmt.Fprint(os.Stdout, string(out))
	return nil
}
