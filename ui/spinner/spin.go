package spinner

import tea "github.com/charmbracelet/bubbletea"

func Spin(
	header string,
	task func() error,
) (bool, error) {
	var ok bool
	tp := tea.NewProgram(New(
		SpinnerTask{
			Header: header,
			Task:   task,
		},
		&ok,
	))
	if _, err := tp.Run(); err != nil {
		return false, err
	}
	return ok, nil
}

