package ui

import (
	"balance/core"
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AccountsModel struct {
	accounts        []core.Account
	isActive        bool
	showModal       bool
	modalFields     []string
	modalFieldIndex int
	width           int
	height          int
}

func NewAccountsModel() (AccountsModel, error) {
	accounts, err := core.LoadAccounts()
	if err != nil {
		return AccountsModel{}, err
	}
	return AccountsModel{
		accounts:    accounts,
		modalFields: []string{"", ""},
	}, nil
}

func (m AccountsModel) Init() tea.Cmd {
	return nil
}

func (m AccountsModel) Update(msg tea.Msg) (AccountsModel, tea.Cmd) {
	return m, nil
}

func (m AccountsModel) View() string {
	title := "Accounts"
	if m.isActive {
		title = "▶ " + title + " ◀"
	}

	var content string
	for _, acc := range m.accounts {
		content += fmt.Sprintf("%s: %s\n", acc.Name, FormatCents(acc.Balance))
	}

	if content == "" {
		content = "[no accounts]"
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Width(m.width - 2).
		Height(m.height)

	if m.isActive {
		style = style.BorderForeground(lipgloss.Color("33"))
	}

	return style.Render(lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(title),
		content,
	))
}

func (m AccountsModel) SetActive(active bool) AccountsModel {
	m.isActive = active
	return m
}

func (m AccountsModel) SetSize(width, height int) AccountsModel {
	m.width = width
	m.height = height
	return m
}

func (m AccountsModel) HandleModalInput(msg tea.KeyMsg) (AccountsModel, bool, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.showModal = false
		m.modalFields = []string{"", ""}
		m.modalFieldIndex = 0
		return m, true, nil
	case "tab":
		m.modalFieldIndex = (m.modalFieldIndex + 1) % len(m.modalFields)
	case "shift+tab":
		m.modalFieldIndex--
		if m.modalFieldIndex < 0 {
			m.modalFieldIndex = len(m.modalFields) - 1
		}
	case "enter":
		if len(m.modalFields) >= 2 && m.modalFields[0] != "" {
			balance := 0
			if m.modalFields[1] != "" {
				if b, err := strconv.ParseInt(m.modalFields[1], 10, 64); err == nil {
					balance = int(b)
				}
			}
			m.accounts = append(m.accounts, core.Account{
				Name:    m.modalFields[0],
				Balance: balance,
			})
			core.SaveAccounts(m.accounts)
		}
		m.showModal = false
		m.modalFields = []string{"", ""}
		m.modalFieldIndex = 0
		return m, true, nil
	case "backspace":
		if len(m.modalFields[m.modalFieldIndex]) > 0 {
			m.modalFields[m.modalFieldIndex] = m.modalFields[m.modalFieldIndex][:len(m.modalFields[m.modalFieldIndex])-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.modalFields[m.modalFieldIndex] += msg.String()
		}
	}
	return m, false, nil
}

func (m AccountsModel) ShowModal() AccountsModel {
	m.showModal = true
	return m
}

func (m AccountsModel) IsShowingModal() bool {
	return m.showModal
}

func (m AccountsModel) RenderModal() string {
	title := "New Account"
	labels := []string{"Name", "Balance"}

	var fieldLines []string
	for i := 0; i < 2; i++ {
		prefix := " "
		if i == m.modalFieldIndex {
			prefix = ">"
		}
		input := m.modalFields[i]
		if input == "" {
			input = "[empty]"
		}
		fieldLines = append(fieldLines, fmt.Sprintf("%s %s: %s", prefix, labels[i], input))
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(title),
		"",
		lipgloss.JoinVertical(lipgloss.Left, fieldLines...),
		"",
		lipgloss.NewStyle().Faint(true).Render("Tab: navigate | Enter: submit | Esc: cancel"),
	)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("69")).
		Padding(1).
		Width(50).
		Render(content)
}

func FormatCents(cents int) string {
	dollars := cents / 100
	remaining := cents % 100
	if remaining < 0 {
		remaining = -remaining
	}
	return fmt.Sprintf("$%d.%02d", dollars, remaining)
}
