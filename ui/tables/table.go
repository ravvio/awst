package tables

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/ravvio/awst/ui/style"
)

type Row = map[string]string

type Alignment int

const (
	Left Alignment = iota
	Right
	Center
)

type Column struct {
	Key       string
	Title     string
	Active    bool
	Alignment Alignment
}

func NewColumn(key string, title string, active bool) Column {
	return Column{
		Key:       key,
		Title:     title,
		Active:    active,
		Alignment: Left,
	}
}

func (c Column) WithAlignment(a Alignment) Column {
	c.Alignment = a
	return c
}

type Table struct {
	columns []Column
	rows    []Row
}

func New(columns []Column) Table {
	return Table{
		columns: columns,
		rows:    []Row{},
	}
}

func (t Table) WithRows(rows []Row) Table {
	t.rows = rows
	return t
}

func (t Table) Render() string {
	aligments := []Alignment{}
	headers := []string{}
	for _, col := range t.columns {
		if !col.Active {
			continue
		}
		headers = append(headers, col.Title)
		aligments = append(aligments, col.Alignment)
	}

	rows := [][]string{}
	for _, rowEntry := range t.rows {
		row := []string{}
		for _, col := range t.columns {
			if !col.Active {
				continue
			}
			row = append(row, rowEntry[col.Key])
		}
		rows = append(rows, row)
	}

	lt := table.New().
		Headers(headers...).
		Rows(rows...).
		Border(lipgloss.NormalBorder()).
		BorderLeft(false).BorderRight(false).BorderTop(false).BorderBottom(false).
		BorderColumn(false).BorderHeader(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			var sty lipgloss.Style

			switch {
			case row == table.HeaderRow:
				sty = style.HeaderStyle
			default:
				sty = style.RowStyle
			}

			switch aligments[col] {
			case Left:
				sty = sty.Align(lipgloss.Left)
			case Center:
				sty = sty.Align(lipgloss.Center)
			case Right:
				sty = sty.Align(lipgloss.Right)
			}

			return sty
		})

	return lt.Render()
}
