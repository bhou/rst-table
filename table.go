package rsttable

import (
	"fmt"
	"sort"
	"strings"
)

// Row is a map of column name to value.
type Row any
type ColRender func(Row, string) string
type lessFunc func(Row, Row) bool

type Table struct {
	Rows             []Row
	Cols             []string
	ColRenders       map[string]ColRender
	DefaultColRender ColRender
}

// default col render
func DefaultColRender(row Row, col string) string {
	rowMap, ok := row.(map[string]any)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", rowMap[col])
}

func (t *Table) AddRow(row Row) {
	t.Rows = append(t.Rows, row)
}

func NewTable() *Table {
	return &Table{
		Rows:             []Row{},
		Cols:             []string{},
		ColRenders:       make(map[string]ColRender),
		DefaultColRender: DefaultColRender,
	}
}

func (t *Table) AddCol(col string, render ColRender) {
	t.Cols = append(t.Cols, col)
	t.ColRenders[col] = render
}

// GenerateRstTable generates a reStructuredText table.
func (t Table) GenerateRstTable(groupBy []string) string {
	// reorder the rows according to the groupBy
	sort.Slice(t.Rows, func(i, j int) bool {
		for _, groupCol := range groupBy {
			if t.ColRenders[groupCol](t.Rows[i], groupCol) < t.ColRenders[groupCol](t.Rows[j], groupCol) {
				return true
			}
		}
		return false
	})

	// first set the table value to string using the col render
	maxColLen := 0
	strRows := make([]map[string]string, 0)
	for _, col := range t.Cols {
		if len(col) > maxColLen {
			maxColLen = len(col)
		}
	}
	for _, row := range t.Rows {
		r := make(map[string]string)
		for _, col := range t.Cols {
			r[col] = t.ColRenders[col](row, col)
			if len(r[col]) > maxColLen {
				maxColLen = len(r[col])
			}
		}
		strRows = append(strRows, r)
	}

	// reorder the cols according to the groupBy
	displayCols := []string{}
	for _, groupCol := range groupBy {
		displayCols = append(displayCols, groupCol)
	}
	for _, col := range t.Cols {
		found := false
		for _, groupCol := range groupBy {
			if col == groupCol {
				found = true
				break
			}
		}
		if !found {
			displayCols = append(displayCols, col)
		}
	}

	// render the header
	buf := ""
	buf += tableSplitLine("-", len(t.Cols), maxColLen, []int{})
	buf += "|"
	for _, colName := range displayCols {
		buf += fmt.Sprintf(colFormat(maxColLen), colName)
	}
	buf += "\n"

	lastRow := make([]string, len(displayCols))

	// render the row
	for index, row := range strRows {
		fillChar := "-"
		if index == 0 {
			fillChar = "="
		}

		topLine := "+"
		valueLine := "|"
		// now render the content
		changedColIndex := 0
		for i, col := range displayCols {
			if lastRow[i] != row[col] || changedColIndex < i {
				changedColIndex = i
				topLine += splitCell(fillChar, maxColLen)
				valueLine += fmt.Sprintf(colFormat(maxColLen), row[col])
				lastRow[i] = row[col]
			} else {
				changedColIndex += 1
				topLine += splitCell(" ", maxColLen)
				valueLine += fmt.Sprintf(colFormat(maxColLen), " ")
			}
		}
		topLine += "\n"
		valueLine += "\n"
		buf += topLine + valueLine
	}

	// render the last line
	buf += tableSplitLine("-", len(t.Cols), maxColLen, []int{})

	return buf
}

// Draw a table split line.
func tableSplitLine(c string, nCols int, maxColLen int, skips []int) string {
	buf := "+"
	if len(skips) > nCols {
		return ""
	}
	for i := 0; i < nCols; i++ {
		doSkip := false
		for _, skip := range skips {
			if skip == i {
				doSkip = true
				break
			}
		}
		if doSkip {
			buf += strings.Repeat(" ", maxColLen+2) + "+"
		} else {
			buf += strings.Repeat(c, maxColLen+2) + "+"
		}
	}
	buf += "\n"
	return buf
}

func colFormat(n int) string {
	return fmt.Sprintf(" %%-%ds |", n)
}
func splitCell(c string, n int) string {
	return strings.Repeat(c, n+2) + "+"
}
