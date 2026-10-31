package cli

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const emptyText = "No items"

// Table type.
type Table struct {
	Headers []string
	Rows    [][]string
}

// Create function.
func (t *Table) Create() []string {
	width := make([]int, len(t.Headers))

	for i, header := range t.Headers {
		width[i] = utf8.RuneCountInString(header)
	}

	for _, row := range t.Rows {
		for i, col := range row[:min(len(row), len(width))] {
			width[i] = max(width[i], utf8.RuneCountInString(col))
		}
	}

	if len(t.Rows) == 0 {
		if len(width) == 0 {
			width = append(width, 0)
		}

		total := (len(width) - 1) * 3

		for _, col := range width {
			total += col
		}

		width[len(width)-1] += max(len(emptyText)-total, 0)
	}

	var height int

	if len(t.Rows) == 0 {
		height = 5
	} else {
		height = 4 + len(t.Rows)
	}

	table := make([]string, height)
	border := t.createBorder(width)

	table[0] = border
	table[1] = t.createRow(t.Headers, width)
	table[2] = border

	if len(t.Rows) == 0 {
		table[3] = t.createEmpty(width)
	} else {
		for i, row := range t.Rows {
			table[3+i] = t.createRow(row, width)
		}
	}

	table[height-1] = border

	return table
}

// Print function.
func (t *Table) Print() int {
	for _, tr := range t.Create() {
		fmt.Println(tr)
	}

	return Success
}

func (t *Table) createBorder(width []int) string {
	var border strings.Builder

	for i := range width {
		if i == 0 {
			border.WriteString("+-")
		} else {
			border.WriteString("-")
		}

		border.WriteString(strings.Repeat("-", width[i]+1))

		if i == len(width)-1 {
			border.WriteString("+")
		} else {
			border.WriteString("-")
		}
	}

	return border.String()
}

func (t *Table) createEmpty(width []int) string {
	var total int

	for _, col := range width {
		total += col
	}

	total += (len(width) - 1) * 3

	return "| " + emptyText + strings.Repeat(" ", total-len(emptyText)) + " |"
}

func (t *Table) createRow(row []string, width []int) string {
	var tr strings.Builder

	tr.WriteString("|")

	for i := range width {
		var col string

		if i < len(row) {
			col = row[i]
		}

		tr.WriteString(" " + col + strings.Repeat(" ", width[i]-utf8.RuneCountInString(col)) + " |")
	}

	return tr.String()
}
