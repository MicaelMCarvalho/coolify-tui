package ui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
)

func filterablePanel(target panel) bool {
	switch target {
	case projectsPanel,
		environmentsPanel,
		resourcesPanel,
		environmentVariablesPanel,
		deploymentsPanel:
		return true

	default:
		return false
	}
}

func (m Model) panelSearchItems(
	target panel,
) []string {
	switch target {
	case projectsPanel:
		items := make([]string, 0, len(m.projects))
		for _, project := range m.projects {
			items = append(items, project.Name)
		}
		return items

	case environmentsPanel:
		if m.project == nil {
			return nil
		}

		items := make(
			[]string,
			0,
			len(m.project.Environments),
		)

		for _, environment := range m.project.Environments {
			items = append(
				items,
				environment.Name,
			)
		}

		return items

	case resourcesPanel:
		items := make(
			[]string,
			0,
			len(m.resources),
		)

		for _, resource := range m.resources {
			items = append(
				items,
				fmt.Sprintf(
					"%s %s %s",
					resource.Type,
					resource.Name,
					resource.Status,
				),
			)
		}

		return items

	case environmentVariablesPanel:
		items := make(
			[]string,
			0,
			len(m.environmentVariables),
		)

		for _, variable := range m.environmentVariables {
			items = append(
				items,
				variable.Key)
		}

		return items

	case deploymentsPanel:
		items := make(
			[]string,
			0,
			len(m.deployments),
		)
		for _, deployment := range m.deployments {
			item := deployment.Status +
				" " + deployment.Commit

			if deployment.CommitMessage != nil {
				item += " " +
					*deployment.CommitMessage
			}
			items = append(items, item)
		}

		return items

	default:
		return nil
	}
}

func (m Model) filteredIndices(
	target panel,
) []int {
	items := m.panelSearchItems(target)
	query := strings.ToLower(
		strings.TrimSpace(m.filters[target]),
	)

	indices := make([]int, 0, len(items))

	for index, item := range items {
		if query == "" ||
			strings.Contains(
				strings.ToLower(item),
				query,
			) {
			indices = append(indices, index)
		}
	}

	return indices
}

func filteredCursorPosition(
	indices []int,
	rawCursor int,
) int {
	for position, index := range indices {
		if index == rawCursor {
			return position
		}
	}
	if len(indices) > 0 {
		return 0
	}
	return -1
}

func filterItemsByIndices(
	items []string,
	indices []int,
) []string {
	filtered := make([]string, 0, len(indices))
	for _, index := range indices {
		if index >= 0 && index < len(items) {
			filtered = append(
				filtered,
				items[index],
			)
		}
	}
	return filtered
}

func nextFilteredIndex(
	indices []int,
	current int,
	change int,
) (int, bool) {
	if len(indices) == 0 {
		return 0, false
	}

	position := filteredCursorPosition(
		indices,
		current,
	)

	next := position + change

	if next < 0 || next >= len(indices) {
		return 0, false
	}

	return indices[next], true
}

func removeLastRune(value string) string {
	if value == "" {
		return ""
	}

	_, size := utf8.DecodeLastRuneInString(value)

	return value[:len(value)-size]
}

func (m Model) handleFilterKey(
	msg tea.KeyMsg,
) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.filters[m.filterPanel] = m.filterOriginal

		m.filtering = false
		m.filterInput = ""

		return m, nil

	case "enter":
		m.filtering = false

		return m, m.moveToBoundary(true)

	case "backspace":
		m.filterInput = removeLastRune(m.filterInput)

	case "ctrl+u":
		m.filterInput = ""

	default:
		if msg.Type == tea.KeyRunes {
			m.filterInput += string(msg.Runes)
		}
	}

	m.filters[m.filterPanel] = m.filterInput

	return m, nil
}
