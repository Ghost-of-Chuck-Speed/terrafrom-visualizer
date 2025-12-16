package ui

import (
	"fmt"
	"os"
	"sort"

	"tfviz/internal/group"
	"tfviz/internal/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

func NewModel(groups map[string]*group.Group) *Model {
	// Create a sorted list of group names for consistent ordering.
	groupNames := make([]string, 0, len(groups))
	for k := range groups {
		groupNames = append(groupNames, k)
	}
	sort.Strings(groupNames)

	return &Model{
		groups:     groups,
		groupNames: groupNames,
		groupKeys:  groupNames,
		focus:      groupPane,
	}
}

type pane int

const (
	groupPane pane = iota
	resourcePane
)

type Model struct {
	groups     map[string]*group.Group
	groupNames []string
	cursor     int
	groups         map[string]*group.Group
	groupKeys      []string
	groupCursor    int
	resourceCursor int
	focus          pane
	width, height  int
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			if m.focus == groupPane && m.groupCursor > 0 {
				m.groupCursor--
				m.resourceCursor = 0 // Reset resource cursor when changing groups
			} else if m.focus == resourcePane && m.resourceCursor > 0 {
				m.resourceCursor--
			}
		case "down", "j":
			if m.cursor < len(m.groupNames)-1 {
				m.cursor++
			selectedGroup := m.groups[m.groupKeys[m.groupCursor]]
			if m.focus == groupPane && m.groupCursor < len(m.groupKeys)-1 {
				m.groupCursor++
				m.resourceCursor = 0 // Reset resource cursor when changing groups
			} else if m.focus == resourcePane && m.resourceCursor < len(selectedGroup.Resources)-1 {
				m.resourceCursor++
			}
		case "tab":
			m.focus = resourcePane
		case "shift+tab":
			m.focus = groupPane
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *Model) View() string {
	s := "Resource Groups:\n\n"
	for i, groupKey := range m.groupNames {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
	if m.width == 0 {
		return "Initializing..."
	}

	groupView := m.renderPane("Resource Groups", m.groupKeys, m.renderGroup, m.groupCursor, m.focus == groupPane)

	selectedGroup := m.groups[m.groupKeys[m.groupCursor]]
	resourceView := m.renderPane("Resources", selectedGroup.Resources, m.renderResource, m.resourceCursor, m.focus == resourcePane)

	help := "q: quit | ←/→ or tab: switch panes | ↑/↓: navigate"

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, groupView, resourceView)
	return lipgloss.JoinVertical(lipgloss.Left, mainContent, help)
}

func (m *Model) renderPane(title string, items interface{}, renderer func(i int, item interface{}) string, cursor int, hasFocus bool) string {
	var s string
	var listItems []string

	switch v := items.(type) {
	case []string:
		for i, item := range v {
			listItems = append(listItems, renderer(i, item))
		}
		g := m.groups[groupKey]
		s += fmt.Sprintf("%s %s (%d resources)\n", cursor, g.Name, len(g.Resources))
	case []*model.Resource:
		for i, item := range v {
			listItems = append(listItems, renderer(i, item))
		}
	}
	s += "\nPress q to quit.\n"
	return s

	s = lipgloss.JoinVertical(lipgloss.Left, listItems...)

	paneStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Height(m.height - 3).
		Width(m.width/2 - 4)

	if hasFocus {
		paneStyle = paneStyle.Border(lipgloss.NormalBorder(), false, false, false, false, title).
			BorderForeground(lipgloss.Color("63"))
	} else {
		paneStyle = paneStyle.Border(lipgloss.HiddenBorder(), false, false, false, false, title)
	}

	return paneStyle.Render(s)
}

func (m *Model) renderGroup(i int, item interface{}) string {
	groupKey := item.(string)
	g := m.groups[groupKey]
	cursor := " "
	if m.groupCursor == i {
		cursor = ">"
	}
	return fmt.Sprintf("%s %s (%d)", cursor, g.Name, len(g.Resources))
}

func (m *Model) renderResource(i int, item interface{}) string {
	res := item.(*model.Resource)
	cursor := " "
	if m.resourceCursor == i {
		cursor = ">"
	}
	return fmt.Sprintf("%s %s", cursor, res.Address)
}

// Run starts the TUI.
func Run(m *Model) error {
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
