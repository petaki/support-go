package cli

import (
	"fmt"
	"strings"
)

type Table struct {
	Headers []string
	Rows    [][]string
}

func (t *Table) Create() []string {
	width := make([]int, len(t.Headers))

	for i, header := range t.Headers {
		width[i] = len(header)
	}

	for _, row := range t.Rows {
		for i, col := range row {
			if width[i] < len(col) {
				width[i] = len(col)
			}
		}
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

func (t *Table) Print() int {
	for _, tr := range t.Create() {
		fmt.Println(tr)
	}

	return 0
}

func (t *Table) createBorder(width []int) string {
	var border string

	for i, _ := range width {
		if i == 0 {
			border += "+-"
		} else {
			border += "-"
		}

		border += strings.Repeat("-", width[i]+1)

		if i == len(width)-1 {
			border += "+"
		} else {
			border += "-"
		}
	}

	return border
}

func (t *Table) createEmpty(width []int) string {
	var total int

	for _, col := range width {
		total += col
	}

	total += (len(width) - 1) * 3
	text := "No items"

	return "| " + text + strings.Repeat(" ", total-len(text)) + " |"
}

func (t *Table) createRow(row []string, width []int) string {
	var tr string

	for i, col := range row {
		if i == 0 {
			tr += "|"
		}

		tr += " " + col + strings.Repeat(" ", width[i]-len(col)) + " |"
	}

	return tr
}
