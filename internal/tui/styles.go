package tui

import "github.com/charmbracelet/lipgloss"

const (
	padding  = 2
	maxWidth = 80
)

var (
	// Private styles (used within the tui package)
	doneStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
	errStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
	labelStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA"))
	timerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0"))
	outputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7DCFFF"))

	// Exported styles (used in main.go)
	ErrStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)
	WarnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAA00")).Bold(true)
	InfoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0"))
)
