package style

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

const (
	PrimaryFg = lipgloss.Color("0")
	Primary   = lipgloss.Color("4")

	SecondaryFg = lipgloss.Color("0")
	Secondary   = lipgloss.Color("3")

	Fg    = lipgloss.Color("15")
	DimFg = lipgloss.Color("8")

	Accent     = lipgloss.Color("3")
	ErrorFg    = lipgloss.Color("1")
	SuccessFg  = lipgloss.Color("2")
	ProgressFg = lipgloss.Color("11")
)

var (
	AccentStyle  = lipgloss.NewStyle().Foreground(Accent)
	SuccessStyle = lipgloss.NewStyle().PaddingLeft(1).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(SuccessFg).
			Width(100)
	ProgressStyle = lipgloss.NewStyle().PaddingLeft(1).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(ProgressFg)
	InfoStyle = lipgloss.NewStyle().PaddingLeft(1).
			Faint(true).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(Primary).
			Width(100)
	DescStyle = lipgloss.NewStyle().Faint(true).PaddingLeft(1).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(Primary).
			Width(100)
	HintStyle  = lipgloss.NewStyle().Faint(true).Width(100)
	ErrorStyle = lipgloss.NewStyle().Faint(true).PaddingLeft(1).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(ErrorFg).
			Width(100)

	TitleStyle = lipgloss.NewStyle().Foreground(Primary).Bold(true)

	TableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).BorderForeground(DimFg).
			BorderTop(true).BorderBottom(true)
	HeaderStyle = lipgloss.NewStyle().Foreground(Primary).Bold(true).Padding(0, 1)
	RowStyle    = lipgloss.NewStyle().Padding(0, 1)

	LogTitle   = lipgloss.NewStyle().Foreground(Primary).PaddingRight(1)
	LogDate    = lipgloss.NewStyle().Foreground(Secondary).PaddingRight(1)
	LogContent = lipgloss.NewStyle().PaddingRight(1)
)

func StyleError(err string, args ...any) string {
	return ErrorStyle.Render(fmt.Sprintf(err, args...))
}

func PrintError(err string, args ...any) {
	fmt.Fprintln(os.Stderr, StyleError(err, args...))
}

func StyleHint(hint string, args ...any) string {
	return HintStyle.Render(fmt.Sprintf(hint, args...))
}

func PrintHint(hint string, args ...any) {
	fmt.Fprintln(os.Stderr, StyleHint(hint, args...))
}

func StyleInfo(info string, args ...any) string {
	return InfoStyle.Render(fmt.Sprintf(info, args...))
}

func PrintInfo(info string, args ...any) {
	fmt.Fprintln(os.Stderr, StyleInfo(info, args...))
}
