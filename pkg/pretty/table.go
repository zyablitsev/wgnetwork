package pretty

import (
	"fmt"
	"strings"
)

// Table helps print pretty-table.
type Table struct {
	columns     int
	columnSize  []int
	columnAlign []int
	header      []string
	rows        [][]string
	footer      []string
}

// NewTable constructor.
func NewTable(columns int) *Table {
	t := &Table{
		columns:     columns,
		columnSize:  make([]int, columns),
		columnAlign: make([]int, columns),
		header:      nil,
		rows:        [][]string{},
		footer:      make([]string, columns),
	}

	return t
}

// SetColumnMinSize method.
func (t *Table) SetColumnMinSize(idx, size int) error {
	if idx > t.columns-1 || idx < 0 {
		return fmt.Errorf("column index should be in range 0…%d, got %d",
			t.columns-1, idx)
	}

	if size > t.columnSize[idx] {
		t.columnSize[idx] = size
	}

	return nil
}

// SetColumnAlign method.
func (t *Table) SetColumnAlign(idx, align int) error {
	if idx > t.columns-1 || idx < 0 {
		return fmt.Errorf("column index should be in range 0…%d, got %d",
			t.columns-1, idx)
	}

	if align > TextAlignRight || align < TextAlignLeft {
		return fmt.Errorf("align value should be in range 0…%d, got %d",
			TextAlignRight, align)
	}

	t.columnAlign[idx] = align

	return nil
}

// SetHeader sets the header row columns content to render.
func (t *Table) SetHeader(columns []string) error {
	if len(columns) != t.columns {
		return fmt.Errorf("columns expected: %d, got %d",
			t.columns, len(columns))
	}

	t.header = columns

	for idx, s := range columns {
		if len(s) > t.columnSize[idx] {
			t.columnSize[idx] = len(s)
		}
	}

	return nil
}

// AddRow sets the footer row columns content to render.
func (t *Table) AddRow(columns []string) error {
	if len(columns) != t.columns {
		return fmt.Errorf("columns expected: %d, got %d",
			t.columns, len(columns))
	}

	rows := [][]string{}
	for i := 0; i < len(columns); i++ {
		p := strings.Split(columns[i], "\n")

		for j := 0; j < len(p); j++ {
			if len(rows) < j+1 {
				row := make([]string, len(columns))
				rows = append(rows, row)
			}

			rows[j][i] = p[j]
		}
	}

	for _, row := range rows {
		t.rows = append(t.rows, row)

		for idx, s := range row {
			if len(s) > t.columnSize[idx] {
				t.columnSize[idx] = len(s)
			}
		}
	}

	return nil
}

// SetFooter sets the footer row columns content to render.
func (t *Table) SetFooter(columns []string) error {
	if len(columns) != t.columns {
		return fmt.Errorf("columns expected: %d, got %d",
			t.columns, len(columns))
	}

	t.footer = columns

	for idx, s := range columns {
		if len(s) > t.columnSize[idx] {
			t.columnSize[idx] = len(s)
		}
	}

	return nil
}

// Render table.
func (t *Table) Render() string {
	var w strings.Builder

	if len(t.header) > 0 {
		for idx, s := range t.header {
			s = AlignText(s, t.columnSize[idx], t.columnAlign[idx])
			w.WriteString(s)
			w.WriteString(" ")
		}
		w.WriteString("\n")
	}

	for i := 0; i < t.columns; i++ {
		w.WriteString(strings.Repeat("-", t.columnSize[i]))
		w.WriteString(" ")
	}
	w.WriteString("\n")

	for _, r := range t.rows {
		for idx, s := range r {
			s = AlignText(s, t.columnSize[idx], t.columnAlign[idx])
			w.WriteString(s)
			w.WriteString(" ")
		}
		w.WriteString("\n")
	}

	return w.String()
}
