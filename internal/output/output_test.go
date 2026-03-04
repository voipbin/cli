package output

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type testItem struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

var testColumns = []Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "STATUS", Field: "status"},
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		format  string
		wantErr bool
	}{
		{"table", false},
		{"json", false},
		{"yaml", false},
		{"", false},
		{"xml", true},
	}

	for _, tt := range tests {
		_, err := NewFormatter(tt.format)
		if (err != nil) != tt.wantErr {
			t.Errorf("NewFormatter(%q) error = %v, wantErr %v", tt.format, err, tt.wantErr)
		}
	}
}

func TestJSONFormatterList(t *testing.T) {
	items := []testItem{
		{ID: "1", Name: "alpha", Status: "active"},
		{ID: "2", Name: "beta", Status: "inactive"},
	}

	out := captureStdout(t, func() {
		f := &JSONFormatter{}
		if err := f.FormatList(items, testColumns); err != nil {
			t.Fatalf("FormatList error: %v", err)
		}
	})

	var parsed []testItem
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, out)
	}
	if len(parsed) != 2 {
		t.Errorf("expected 2 items, got %d", len(parsed))
	}
}

func TestJSONFormatterItem(t *testing.T) {
	item := testItem{ID: "1", Name: "alpha", Status: "active"}

	out := captureStdout(t, func() {
		f := &JSONFormatter{}
		if err := f.FormatItem(item, testColumns); err != nil {
			t.Fatalf("FormatItem error: %v", err)
		}
	})

	var parsed testItem
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, out)
	}
	if parsed.ID != "1" {
		t.Errorf("expected ID=1, got %q", parsed.ID)
	}
}

func TestYAMLFormatterList(t *testing.T) {
	items := []testItem{
		{ID: "1", Name: "alpha", Status: "active"},
	}

	out := captureStdout(t, func() {
		f := &YAMLFormatter{}
		if err := f.FormatList(items, testColumns); err != nil {
			t.Fatalf("FormatList error: %v", err)
		}
	})

	var parsed []map[string]interface{}
	if err := yaml.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid YAML: %v\noutput: %s", err, out)
	}
	if len(parsed) != 1 {
		t.Errorf("expected 1 item, got %d", len(parsed))
	}
}

func TestTableFormatterList(t *testing.T) {
	items := []testItem{
		{ID: "1", Name: "alpha", Status: "active"},
		{ID: "2", Name: "beta", Status: "inactive"},
	}

	out := captureStdout(t, func() {
		f := &TableFormatter{}
		if err := f.FormatList(items, testColumns); err != nil {
			t.Fatalf("FormatList error: %v", err)
		}
	})

	if !strings.Contains(out, "ID") {
		t.Error("table output missing header 'ID'")
	}
	if !strings.Contains(out, "alpha") {
		t.Error("table output missing data 'alpha'")
	}
	if !strings.Contains(out, "beta") {
		t.Error("table output missing data 'beta'")
	}
}

func TestTableFormatterItem(t *testing.T) {
	item := testItem{ID: "1", Name: "alpha", Status: "active"}

	out := captureStdout(t, func() {
		f := &TableFormatter{}
		if err := f.FormatItem(item, testColumns); err != nil {
			t.Fatalf("FormatItem error: %v", err)
		}
	})

	if !strings.Contains(out, "ID:") {
		t.Error("item output missing 'ID:'")
	}
	if !strings.Contains(out, "alpha") {
		t.Error("item output missing data 'alpha'")
	}
}

func TestTableFormatterEmptyList(t *testing.T) {
	out := captureStdout(t, func() {
		f := &TableFormatter{}
		if err := f.FormatList([]testItem{}, testColumns); err != nil {
			t.Fatalf("FormatList error: %v", err)
		}
	})

	if !strings.Contains(out, "No items found") {
		t.Errorf("expected 'No items found' message, got: %s", out)
	}
}

func TestExtractFieldMissing(t *testing.T) {
	item := map[string]interface{}{"id": "1"}
	result := extractField(item, "nonexistent")
	if result != "" {
		t.Errorf("expected empty string for missing field, got %q", result)
	}
}
