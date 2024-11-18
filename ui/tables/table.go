package tables

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/ravvio/awst/ui/style"
)

type Row = map[string]string

type Column struct {
    Key string
    Title string
    Active bool
}

func NewColumn(key string, title string, active bool) Column {
    return Column{
        Key: key,
        Title: title,
        Active: active,
    }
}

type Table struct {
    columns []Column
    rows []Row
}

func New(columns []Column) Table {
    return Table{
        columns: columns,
        rows: []Row{},
    }
}

func (t Table) WithRows(rows []Row) Table {
    t.rows = rows
    return t
}

func (t Table) Render() string {
    headers := []string{}
    for _, col := range t.columns {
        if !col.Active { continue }
        headers = append(headers, col.Title)
    }

    rows := [][]string{}
    for _, rowEntry := range t.rows {
        row := []string{}
        for _, col := range t.columns {
            if !col.Active { continue }
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
            switch {
            case row == table.HeaderRow:
                return style.HeaderStyle
            case row%2 == 0:
                return style.RowStyle
            default:
                return style.RowStyle
            }
        })

    return lt.Render()
}
