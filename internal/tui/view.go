package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

// View returns the current view for the progress model
func (m ProgressModel) View() tea.View {
	pad := strings.Repeat(" ", padding)
	var s strings.Builder

	s.WriteString("\n")

	// Handle done state
	if m.done {
		elapsed := fmt.Sprintf("%.2fs", m.elapsed.Seconds())
		if m.err != nil {
			s.WriteString(pad + errStyle.Render("✗ Error: "+m.err.Error()) + "\n")
		} else {
			s.WriteString(pad + doneStyle.Render("✓ Done!"))
			s.WriteString("  " + timerStyle.Render(elapsed) + "\n")
			s.WriteString("\n")
			s.WriteString(pad + labelStyle.Render("Output: ") + outputStyle.Render(m.output) + "\n")
		}
	} else {
		elapsed := fmt.Sprintf("%.2fs", time.Since(m.startTime).Seconds())
		s.WriteString(pad + labelStyle.Render("Processing...") + "  " + timerStyle.Render(elapsed) + "\n\n")
		s.WriteString(pad + m.bar.View() + "\n\n")
		s.WriteString(pad + helpStyle.Render("ctrl+c to cancel"))
	}

	s.WriteString("\n")
	return tea.NewView(s.String())
}
