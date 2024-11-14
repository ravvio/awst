package tables

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ravvio/awst/ui/style"
)

type InteractiveTable struct {
	table table.Model
	quit *bool
}

func NewInteractive(cols []table.Column, quit *bool) InteractiveTable {
	s := table.DefaultStyles()
	s.Header = s.Header.
		Foreground(style.Primary).
		Bold(true).
		BorderBottom(true)
	s.Selected = s.Selected.
		Foreground(style.Accent).
		Bold(true)

	return InteractiveTable{
		table: table.New(
			table.WithColumns(cols),
			table.WithHeight(20),
			table.WithFocused(true),
			table.WithStyles(s),
		),
		quit: quit,
	}
}

func (m InteractiveTable) WithRows(rows []table.Row) InteractiveTable {
	m.table.SetRows(rows)
	if len(rows) < m.table.Height() {
		m.table.SetHeight(len(rows) + 2)
	}
	return m;
}

func (m *InteractiveTable) SelectedRow() table.Row {
	return m.table.SelectedRow()
}

func (m InteractiveTable) Init() tea.Cmd {
	return nil
}

func (m InteractiveTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			*m.quit = false
			cmds = append(cmds, tea.Quit)
		case "ctrl+c", "esc", "q":
			*m.quit = true
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m InteractiveTable) View() string {
	body := strings.Builder{}

	body.WriteString(m.table.View())

	return body.String()
}
