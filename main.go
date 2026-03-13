package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"watermark/internal/engine"
	"watermark/internal/tui"

	tea "charm.land/bubbletea/v2"
	"github.com/alexflint/go-arg"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func logError(msg string) {
	fmt.Fprintln(os.Stderr, tui.ErrStyle.Render("✗ "+msg))
}

type args struct {
	InputPath  string `arg:"-i,--input,required" help:"Path to the input file"`
	Text       string `arg:"-t,--text,required" help:"Text to add as a watermark"`
	OutputPath string `arg:"-o,--output" help:"Path to the output file"`
}

func (args) Description() string {
	return "Simple program to add a watermark to a PDF file or image"
}

func parseArgs() args {
	var a args
	parser, err := arg.NewParser(arg.Config{}, &a)
	if err != nil {
		logError(err.Error())
		os.Exit(0)
	}
	if err := parser.Parse(os.Args[1:]); err != nil {
		if err == arg.ErrHelp {
			parser.WriteHelp(os.Stdout)
			os.Exit(0)
		}
		if err == arg.ErrVersion {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, tui.ErrStyle.Render("✗ "+err.Error()))
		parser.WriteUsage(os.Stderr)
		os.Exit(0)
	}
	return a
}

func main() {
	// Parse command line arguments
	a := parseArgs()

	// Set default output path if not provided
	if a.OutputPath == "" {
		ext := strings.ToLower(filepath.Ext(a.InputPath))
		a.OutputPath = "output" + ext
	}

	// Check if input file exists
	if _, err := os.Stat(a.InputPath); os.IsNotExist(err) {
		logError("File not found: " + a.InputPath)
		os.Exit(0)
	}

	// Check if output file exists and prompt to overwrite
	if _, err := os.Stat(a.OutputPath); !os.IsNotExist(err) {
		theme := huh.ThemeBase()
		theme.Focused.FocusedButton = theme.Focused.FocusedButton.
			Background(lipgloss.Color("#FF2200")).
			Foreground(lipgloss.Color("#FFFFFF"))

		var overwrite bool
		fmt.Println()
		err := huh.NewConfirm().
			Title(lipgloss.NewStyle().Bold(true).Render(filepath.Base(a.OutputPath)) + " already exists. Overwrite?").
			Affirmative("Yes").
			Negative("No").
			Value(&overwrite).
			WithTheme(theme).
			Run()
		if err != nil || !overwrite {
			os.Exit(0)
		}
	}

	// Check if file type is supported
	ext := strings.ToLower(filepath.Ext(a.InputPath))
	allowed := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".pdf": true}
	if !allowed[ext] {
		logError("Unsupported file type: " + ext)
		os.Exit(0)
	}

	// Start processing
	m := tui.NewProgressModel()
	prog := tea.NewProgram(m)

	go func() {
		var err error
		switch strings.ToLower(filepath.Ext(a.InputPath)) {
		case ".jpg", ".jpeg", ".png": // Image files
			_, errStr := engine.ProcessImage(a.InputPath, a.Text, a.OutputPath)
			if errStr != nil {
				err = fmt.Errorf("%s", *errStr)
			}
		case ".pdf": // PDF files
			_, errStr := engine.ProcessPDF(a.InputPath, a.Text, a.OutputPath)
			if errStr != nil {
				err = fmt.Errorf("%s", *errStr)
			}
		}
		prog.Send(tui.DoneMsg{Err: err, OutputPath: a.OutputPath})
	}()

	if _, err := prog.Run(); err != nil {
		logError("Unexpected error: " + err.Error())
		os.Exit(1)
	}
}
