package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
)

func ParseFormat(input string) (Format, error) {
	normalized := strings.TrimSpace(strings.ToLower(input))
	if normalized == "" || normalized == string(FormatTable) {
		return FormatTable, nil
	}
	if normalized == string(FormatJSON) {
		return FormatJSON, nil
	}
	return "", fmt.Errorf("unsupported format %q", input)
}

func PrintJSON(w io.Writer, value any) error {
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json output: %w", err)
	}
	_, err = fmt.Fprintln(w, string(raw))
	return err
}

func PrintTable(w io.Writer, headers []string, rows [][]string) error {
	if len(headers) == 0 {
		return nil
	}

	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i := range headers {
			if i < len(row) && len(row[i]) > widths[i] {
				widths[i] = len(row[i])
			}
		}
	}

	printRow := func(cols []string) error {
		for i := range headers {
			cell := ""
			if i < len(cols) {
				cell = cols[i]
			}
			if _, err := fmt.Fprintf(w, "%-*s", widths[i], cell); err != nil {
				return err
			}
			if i != len(headers)-1 {
				if _, err := io.WriteString(w, "  "); err != nil {
					return err
				}
			}
		}
		_, err := io.WriteString(w, "\n")
		return err
	}

	if err := printRow(headers); err != nil {
		return err
	}
	separator := make([]string, len(headers))
	for i, w := range widths {
		separator[i] = strings.Repeat("-", w)
	}
	if err := printRow(separator); err != nil {
		return err
	}
	for _, row := range rows {
		if err := printRow(row); err != nil {
			return err
		}
	}

	return nil
}
