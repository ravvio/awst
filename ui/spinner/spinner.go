package spinner

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ravvio/awst/ui/style"
)

type StopSpinnerMsg struct {
	err error
}

func (s StopSpinnerMsg) Error() string {
	return s.err.Error()
}

type SpinnerTask struct {
	Header string
	Task   func() error
}

func (s SpinnerTask) Do() tea.Msg {
	err := s.Task()
	return StopSpinnerMsg{err: err}
}

type model struct {
	spinner spinner.Model
	task    SpinnerTask
	err     error
	done    bool
	ok      *bool
}

func New(task SpinnerTask, ok *bool) model {
	s := spinner.New()
	s.Spinner = spinner.Line

	return model{
		spinner: s,
		task:    task,
		err:     nil,
		done:    false,
		ok:      ok,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.task.Do)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case StopSpinnerMsg:
		m.done = true
		*m.ok = true
		if msg.err != nil {
			*m.ok = false
			m.err = msg
		}
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := ""
	if !m.done {
		s += style.ProgressStyle.Render(fmt.Sprintf("%s %s", m.spinner.View(), m.task.Header))
	} else {
		if m.err != nil {
			s += style.ErrorStyle.Render(fmt.Sprintf("* %s ... Failed: %v", m.task.Header, m.err))
		} else {
			s += style.SuccessStyle.Render(fmt.Sprintf("* %s ... Done", m.task.Header))
		}
	}

	s += "\n"
	return s
}

func (m model) Err() error {
	return m.err
}
