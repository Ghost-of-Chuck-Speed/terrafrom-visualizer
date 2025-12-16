package ui

import (
	"encoding/json"
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
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func NewModel(groups map[string]*group.Group) *Model {
	groupNames := make([]string, 0, len(groups))
	for k := range groups {
		groupNames = append(groupNames, k)
	}
	sort.Strings(groupNames)

	return &Model{
		groups:    groups,
		groupKeys: groupNames,
		focus:     groupPane,
	}
}

type pane int

const (
	groupPane pane = iota
	resourcePane
	detailsPane // We won't focus this, but it's good for clarity
)

type Model struct {
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
			if m.focus == groupPane && m.groupCursor > 0 {
				m.groupCursor--
				m.resourceCursor = 0 // Reset resource cursor when changing groups
				m.resourceCursor = 0
			} else if m.focus == resourcePane && m.resourceCursor > 0 {
				m.resourceCursor--
			}
		case "down", "j":
			selectedGroup := m.groups[m.groupKeys[m.groupCursor]]
			if m.focus == groupPane && m.groupCursor < len(m.groupKeys)-1 {
				m.groupCursor++
				m.resourceCursor = 0 // Reset resource cursor when changing groups
			} else if m.focus == resourcePane && m.resourceCursor < len(selectedGroup.Resources)-1 {
				m.resourceCursor++
				m.resourceCursor = 0
			} else if m.focus == resourcePane {
				selectedGroup := m.groups[m.groupKeys[m.groupCursor]]
				if m.resourceCursor < len(selectedGroup.Resources)-1 {
					m.resourceCursor++
				}
			}
		case "tab":
			m.focus = resourcePane
		case "shift+tab":
			m.focus = groupPane
		case "right", "l", "tab":
			if m.focus == groupPane {
				m.focus = resourcePane
			}
		case "left", "h", "shift+tab":
			if m.focus == resourcePane {
				m.focus = groupPane
			}
		case "enter":
			// Future use: expand/collapse resource instances
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	groupView := m.renderPane("Resource Groups", m.groupKeys, m.renderGroup, m.groupCursor, m.focus == groupPane)
	groupView := m.renderPane("Groups", m.groupKeys, m.renderGroup, m.groupCursor, m.focus == groupPane, m.width/4)

	selectedGroup := m.groups[m.groupKeys[m.groupCursor]]
	resourceView := m.renderPane("Resources", selectedGroup.Resources, m.renderResource, m.resourceCursor, m.focus == resourcePane)
	resourceView := m.renderPane("Resources", selectedGroup.Resources, m.renderResource, m.resourceCursor, m.focus == resourcePane, m.width/4)

	help := "q: quit | ←/→ or tab: switch panes | ↑/↓: navigate"
	detailsView := m.renderDetailsPane()

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, groupView, resourceView)
	help := helpStyle.Render("q: quit | ←/→/tab: switch panes | ↑/↓: navigate")

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, groupView, resourceView, detailsView)
	return lipgloss.JoinVertical(lipgloss.Left, mainContent, help)
}

func (m *Model) renderPane(title string, items interface{}, renderer func(i int, item interface{}) string, cursor int, hasFocus bool) string {
func (m *Model) renderPane(title string, items interface{}, renderer func(i int, item interface{}) string, cursor int, hasFocus bool, width int) string {
	var s string
	var listItems []string

	switch v := items.(type) {
	case []string:
		for i, item := range v {
			listItems = append(listItems, renderer(i, item))
		}
	case []*model.Resource:
		for i, item := range v {
			listItems = append(listItems, renderer(i, item))
		}
	}

	s = lipgloss.JoinVertical(lipgloss.Left, listItems...)

	paneStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Height(m.height - 3).
		Width(m.width/2 - 4)
		Width(width)

	if hasFocus {
		paneStyle = paneStyle.BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63"))
		s = lipgloss.NewStyle().Padding(0, 1).Render(title) + "\n" + s
		s = lipgloss.NewStyle().Padding(0, 1).SetString(title).String() + "\n" + s
	} else {
		paneStyle = paneStyle.Border(lipgloss.HiddenBorder())
		paneStyle = paneStyle.BorderStyle(lipgloss.HiddenBorder())
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

func (m *Model) renderDetailsPane() string {
	width := m.width - (m.width/4)*2 - 6 // Remaining width
	style := lipgloss.NewStyle().Padding(1, 2).Height(m.height - 3).Width(width).BorderStyle(lipgloss.NormalBorder())

	selectedGroup := m.groups[m.groupKeys[m.groupCursor]]
	if len(selectedGroup.Resources) == 0 {
		return style.Render("No resources in this group.")
	}

	// Ensure resourceCursor is valid
	if m.resourceCursor >= len(selectedGroup.Resources) {
		m.resourceCursor = 0
	}
	selectedResource := selectedGroup.Resources[m.resourceCursor]

	// Pretty-print the resource attributes as JSON
	attrBytes, err := json.MarshalIndent(selectedResource.Attributes, "", "  ")
	if err != nil {
		return style.Render(fmt.Sprintf("Error rendering details: %v", err))
	}

	content := fmt.Sprintf("Details for: %s\n\n%s", selectedResource.Address, string(attrBytes))
	return style.Render(content)
}

// Run starts the TUI.
func Run(m *Model) error {
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
