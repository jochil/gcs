package tui

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jochil/dlth/pkg/generator"
	"github.com/jochil/dlth/pkg/parser"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	focusedStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69"))

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type sessionState uint

const (
	tableView = iota
	textView
)

type model struct {
	table      table.Model
	viewport   viewport.Model
	state      sessionState
	candidates []*parser.Candidate
}

func NewCandidateModel(candidates []*parser.Candidate, srcPath string) (*model, error) {
	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Function", Width: 30},
		{Title: "Score", Width: 5},
		{Title: "CC", Width: 3},
		{Title: "Lines", Width: 5},
		{Title: "File", Width: 40},
	}

	rows := []table.Row{}
	for i, c := range candidates {
		cc, err := c.CyclomaticComplexity()
		if err != nil {
			slog.Warn("no control flow graph", "func", c.Function.Name)
		}
		rows = append(rows, table.Row{fmt.Sprint(i), c.Function.Name, fmt.Sprintf("%.2f", c.Score), fmt.Sprint(cc), fmt.Sprint(c.Lines), strings.TrimPrefix(c.Path, srcPath)})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
		table.WithWidth(100),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	v := viewport.New(100, 20)
	m := &model{
		state:      tableView,
		table:      t,
		viewport:   v,
		candidates: candidates,
	}
	return m, nil
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "tab":
			if m.state == tableView {
				m.state = textView
			} else {
				m.state = tableView
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "s":
			// TODO doing this in a command?
			i, _ := strconv.ParseInt(m.table.SelectedRow()[0], 10, 0)
			m.viewport.SetContent(m.candidates[i].Code)
		case "t":
			// TODO doing this in a command?
			i, _ := strconv.ParseInt(m.table.SelectedRow()[0], 10, 0)
			testCode := generator.CreateGoTest(m.candidates[i])
			m.viewport.SetContent(testCode)
		}

		switch m.state {
		case tableView:
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
		case textView:
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s string
  // switch between different view elements
	if m.state == tableView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedStyle.Render(m.table.View()), baseStyle.Render(m.viewport.View()))
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, baseStyle.Render(m.table.View()), focusedStyle.Render(m.viewport.View()))
	}
  s += helpStyle.Render("\ntab: focus next • s: view source code • t: generate test • q: exit\n")
	return s
}
