package ui

import (
	"balance/core"
	"fmt"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TransactionsModel struct {
	transactions    []core.Transaction
	isActive        bool
	showModal       bool
	modalFields     []string
	modalFieldIndex int
	width           int
	height          int
}

func NewTransactionsModel() (TransactionsModel, error) {
	transactions, err := core.LoadTransactions()
	if err != nil {
		return TransactionsModel{}, err
	}
	return TransactionsModel{
		transactions: transactions,
		modalFields:  []string{"", "", ""},
	}, nil
}

func (m TransactionsModel) Init() tea.Cmd {
	return nil
}

func (m TransactionsModel) Update(msg tea.Msg) (TransactionsModel, tea.Cmd) {
	return m, nil
}

func (m TransactionsModel) View() string {
	title := "Transactions"
	if m.isActive {
		title = "▶ " + title + " ◀"
	}

	var content string
	for _, trans := range m.transactions {
		content += fmt.Sprintf("%s → %s: %s\n", trans.From, trans.To, FormatCents(trans.Amount))
	}

	if content == "" {
		content = "[no transactions]"
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

func (m TransactionsModel) SetActive(active bool) TransactionsModel {
	m.isActive = active
	return m
}

func (m TransactionsModel) SetSize(width, height int) TransactionsModel {
	m.width = width
	m.height = height
	return m
}

func (m TransactionsModel) HandleModalInput(msg tea.KeyMsg) (TransactionsModel, bool, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.showModal = false
		m.modalFields = []string{"", "", ""}
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
		if len(m.modalFields) >= 3 && m.modalFields[0] != "" && m.modalFields[1] != "" && m.modalFields[2] != "" {
			if amount, err := strconv.ParseInt(m.modalFields[2], 10, 64); err == nil {
				m.transactions = append(m.transactions, core.Transaction{
					From:      m.modalFields[0],
					To:        m.modalFields[1],
					Amount:    int(amount),
					Timestamp: time.Now(),
				})
				core.SaveTransactions(m.transactions)
			}
		}
		m.showModal = false
		m.modalFields = []string{"", "", ""}
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

func (m TransactionsModel) ShowModal() TransactionsModel {
	m.showModal = true
	return m
}

func (m TransactionsModel) IsShowingModal() bool {
	return m.showModal
}

func (m TransactionsModel) RenderModal() string {
	title := "New Transaction"
	labels := []string{"From", "To", "Amount"}

	var fieldLines []string
	for i := range 10 {
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
