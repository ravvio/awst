package tlog

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	DefaultNameStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).PaddingRight(1)
	DefaultTimestampStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).PaddingRight(1)
	DefaultMessageStyle   = lipgloss.NewStyle().PaddingRight(1)
)

type Log struct {
	GroupName *string
	Timestamp *int64
	Message   *string
}

type LogRenderer struct {
	NameStyle      lipgloss.Style
	TimestampStyle lipgloss.Style
	MessageStyle   lipgloss.Style
	DateFormat     string
}

func DefaultRenderer() LogRenderer {
	return LogRenderer{
		NameStyle:      DefaultNameStyle,
		TimestampStyle: DefaultTimestampStyle,
		MessageStyle:   DefaultMessageStyle,
		DateFormat:     time.RFC3339,
	}
}

func (l *LogRenderer) Render(log *Log) error {
	_, err := fmt.Printf(
		"%s%s%s\n",
		l.NameStyle.Render(*log.GroupName),
		l.TimestampStyle.Render(time.UnixMilli(*log.Timestamp).Format(l.DateFormat)),
		l.MessageStyle.Render(strings.Trim(*log.Message, " \n")),
	)
	return err
}
