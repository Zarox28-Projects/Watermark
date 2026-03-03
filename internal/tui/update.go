package tui

import (
	"time"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"
)

// Init initializes the progress model
func (m ProgressModel) Init() tea.Cmd {
	return tea.Batch(m.timer.Init(), m.bar.SetPercent(0))
}

// Update handles messages for the progress model
func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle key press events
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil

	// Handle window resize events
	case tea.WindowSizeMsg:
		w := msg.Width - padding*2 - 4
		w = min(w, maxWidth)
		m.bar.SetWidth(w)
		return m, nil

	// Handle timer tick events
	case timer.TickMsg:
		if m.finishing {
			return m, nil
		}
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.elapsed = time.Since(m.startTime)

		// Animate bar up to 90% while processing
		if m.bar.Percent() < 0.9 {
			incrCmd := m.bar.IncrPercent(0.02)
			return m, tea.Batch(cmd, incrCmd)
		}
		return m, cmd

	// Handle done message
	case DoneMsg:
		m.finishing = true
		m.err = msg.Err
		m.output = msg.OutputPath
		m.elapsed = time.Since(m.startTime)
		// Fill bar to 100% and wait for the animation to finish
		cmd := m.bar.SetPercent(1.0)
		return m, cmd

	// Handle progress frame events
	case progress.FrameMsg:
		var cmd tea.Cmd
		m.bar, cmd = m.bar.Update(msg)

		// Bar has finished animating to 100%, now show the result
		if m.finishing && !m.bar.IsAnimating() {
			m.done = true
			return m, tea.Quit
		}
		return m, cmd

	default:
		return m, nil
	}
}
