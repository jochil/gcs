package tui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/generator"
	"github.com/jochil/dlth/pkg/types"
)

var (
	detailsStyle = lipgloss.NewStyle()
	baseStyle    = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	pagerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderTop(true).
			BorderBottom(true)
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type sessionState uint

const (
	tableView = iota
	codeView
)

type model struct {
	table      table.Model
	code       viewport.Model
	details    viewport.Model
	state      sessionState
	candidates []*candidate.Candidate
	ready      bool
}

func NewCandidateModel(candidates []*candidate.Candidate, srcPath string) (*model, error) {
	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Function", Width: 40},
		{Title: "Score", Width: 5},
	}

	rows := []table.Row{}
	for i, c := range candidates {
		rows = append(rows, table.Row{
			fmt.Sprint(i),
			c.Function.Name,
			fmt.Sprintf("%.2f", c.Score),
		})
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

	m := &model{
		state:      tableView,
		table:      t,
		details:    viewport.New(100, 30),
		candidates: candidates,
	}
	return m, nil
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.code = viewport.New(msg.Width, msg.Height-5)
			m.ready = true
		} else {
			m.code.Width = msg.Width
			m.code.Height = msg.Height - 5
		}
	case tea.KeyMsg:

		switch msg.String() {
		case "esc":
			if m.state == codeView {
				m.state = tableView
			} else {
				return m, tea.Quit
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "s":
			if m.state == tableView {
				m.state = codeView
				// TODO doing this in a command?
				i, _ := strconv.ParseInt(m.table.SelectedRow()[0], 10, 0)
				m.code.SetContent(string(m.candidates[i].Code))
			}
		case "t":
			if m.state == tableView {
				// TODO doing this in a command?
				m.state = codeView
				i, _ := strconv.ParseInt(m.table.SelectedRow()[0], 10, 0)
				testCode := generator.Render(m.candidates[i])
				m.code.SetContent(testCode)
			}
		}

		switch m.state {
		case tableView:
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
		case codeView:
			m.code, cmd = m.code.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.state == tableView {
		return m.listView()
	} else {
		return m.codeView()
	}
}

func (m model) listView() string {
	if m.table.SelectedRow() != nil {
		current, _ := strconv.ParseInt(m.table.SelectedRow()[0], 10, 0)
		m.details.SetContent(detailsContent(m.candidates[current]))
	}

	s := lipgloss.JoinHorizontal(lipgloss.Top, baseStyle.Render(m.table.View()), detailsStyle.Render(m.details.View()))
	s += helpStyle.Render("\ns: view source code • t: generate test • esc: exit\n")
	return s
}

func (m model) codeView() string {
	s := pagerStyle.Render(m.code.View())
	s += helpStyle.Render("\nesc: back to list\n")
	return s
}

func detailsContent(c *candidate.Candidate) string {
	return fmt.Sprintf(
		`
  Name:     %s
  Package:  %s
  Class:    %s 
  Params:   %s
  Return:   %s
  Static:   %t
  Public:   %t
  Language: %s
  Path:     %s

  # Metrics
  Cyclomatic Complexity:  %d 
  Lines of Code:          %d

`,
		c.Function.Name,
		c.Package,
		c.Class,
		c.Function.Parameters,
		c.Function.ReturnValues,
		c.Function.Static,
		c.Function.Visibility == types.VisibilityPublic,
		c.Language,
		c.Path,
		c.Metrics.CyclomaticComplexity,
		c.Metrics.LinesOfCode,
	)
}
