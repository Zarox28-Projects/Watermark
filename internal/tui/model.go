package tui

import (
	"image/color"
	"time"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/timer"
)

// DoneMsg is sent by the main goroutine when processing is complete.
type DoneMsg struct {
	Err        error
	OutputPath string
}

// ProgressModel is the bubbletea model for the progress bar + elapsed timer.
type ProgressModel struct {
	bar       progress.Model
	timer     timer.Model
	elapsed   time.Duration
	startTime time.Time
	finishing bool // DoneMsg received, waiting for bar animation to complete
	done      bool // Bar animation complete, ready to show result
	err       error
	output    string
}

func NewProgressModel() ProgressModel {
	return ProgressModel{
		bar: progress.New(progress.WithColorFunc(func(total, current float64) color.Color {
			// Dark orange → bright red
			r := uint8(0xCC + current*(0xFF-0xCC))
			g := uint8(0x44 - current*0x44)
			return color.RGBA{R: r, G: g, B: 0x00, A: 0xFF}
		})),
		timer:     timer.New(24*time.Hour, timer.WithInterval(time.Millisecond*100)),
		startTime: time.Now(),
	}
}
