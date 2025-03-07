package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/david-krentzlin/graphscope/internal/analyzer"
	"github.com/fatih/color"
)

type Model struct {
	analyzer    *analyzer.Analyzer
	ProgressBar progress.Model
	TotalFiles  int
	Processed   int
	CurrentFile string
	Done        bool
}

func NewModel(analyzer *analyzer.Analyzer) Model {
	return Model{
		ProgressBar: progress.New(progress.WithDefaultGradient()),
		analyzer:    analyzer,
		TotalFiles:  0,
		Done:        false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch update := msg.(type) {
	case analyzer.ParseSchemaComplete:
		m.Done = true

	case analyzer.ParseSchemaFile:
		m.TotalFiles = update.TotalFiles
		m.CurrentFile = update.CurrentFile
		m.Processed++

		progressValue := float64(m.Processed) / float64(m.TotalFiles)
		if m.Processed == m.TotalFiles {
			m.Done = true
		}

		return m, tea.Batch(
			tea.Println(m.View()),
			m.ProgressBar.SetPercent(progressValue),
		)
	}

	return m, nil
}

// View renders UI output
func (m Model) View() string {
	header := color.CyanString("📜 Parsing GraphQL Schema Files...")
	progressBar := m.ProgressBar.ViewAs(float64(m.Processed) / float64(m.TotalFiles))

	if m.Done {
		return fmt.Sprintf("\n%s\n✅ Parsing Complete!\n", header)
	}

	return fmt.Sprintf(
		"\n%s\n%s\n📂 Processing: %s (%d/%d files)\n",
		header, progressBar, m.CurrentFile, m.Processed, m.TotalFiles,
	)
}
