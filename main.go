package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Tab int

const (
	AccountsTab Tab = iota
	TransactionsTab
)

type ModalMode int

const (
	NoModal ModalMode = iota
	AccountModal
	TransactionModal
)

type model struct {
	accounts        []Account
	transactions    []Transaction
	activeTab       Tab
	modalMode       ModalMode
	modalInput      string
	modalFields     []string
	modalFieldIndex int
	lastSpaceTime   time.Time
	spacePressed    bool
	width           int
	height          int
}

func newModel() (model, error) {
	accounts, err := loadAccounts()
	if err != nil {
		return model{}, err
	}

	transactions, err := loadTransactions()
	if err != nil {
		return model{}, err
	}

	return model{
		accounts:      accounts,
		transactions:  transactions,
		activeTab:     AccountsTab,
		modalMode:     NoModal,
		modalFields:   []string{"", ""},
		lastSpaceTime: time.Now().Add(-time.Second),
	}, nil
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if m.modalMode != NoModal {
			return m.handleModalInput(msg)
		}

		keyStr := msg.String()

		switch keyStr {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.activeTab == AccountsTab {
				m.activeTab = TransactionsTab
			} else {
				m.activeTab = AccountsTab
			}
		case " ":
			if m.spacePressed && time.Since(m.lastSpaceTime) < 500*time.Millisecond {
				m.modalMode = TransactionModal
				m.modalFields = []string{"", "", ""}
				m.modalFieldIndex = 0
				m.spacePressed = false
			} else {
				m.spacePressed = true
				m.lastSpaceTime = time.Now()
			}
		case "a":
			if m.spacePressed && time.Since(m.lastSpaceTime) < 500*time.Millisecond {
				m.modalMode = AccountModal
				m.modalFields = []string{"", ""}
				m.modalFieldIndex = 0
				m.spacePressed = false
			}
		default:
			m.spacePressed = false
		}
	}

	return m, nil
}

func (m model) handleModalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.modalMode = NoModal
		m.modalFields = []string{}
		m.modalFieldIndex = 0
	case "tab":
		m.modalFieldIndex = (m.modalFieldIndex + 1) % len(m.modalFields)
	case "shift+tab":
		m.modalFieldIndex--
		if m.modalFieldIndex < 0 {
			m.modalFieldIndex = len(m.modalFields) - 1
		}
	case "enter":
		m.submitModal()
	case "backspace":
		if len(m.modalFields[m.modalFieldIndex]) > 0 {
			m.modalFields[m.modalFieldIndex] = m.modalFields[m.modalFieldIndex][:len(m.modalFields[m.modalFieldIndex])-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.modalFields[m.modalFieldIndex] += msg.String()
		}
	}
	return m, nil
}

func (m *model) submitModal() {
	switch m.modalMode {
	case AccountModal:
		if len(m.modalFields) >= 2 && m.modalFields[0] != "" {
			balance := 0
			if m.modalFields[1] != "" {
				if b, err := strconv.ParseInt(m.modalFields[1], 10, 64); err == nil {
					balance = int(b)
				}
			}
			m.accounts = append(m.accounts, Account{
				Name:    m.modalFields[0],
				Balance: balance,
			})
			saveAccounts(m.accounts)
		}
	case TransactionModal:
		if len(m.modalFields) >= 3 && m.modalFields[0] != "" && m.modalFields[1] != "" && m.modalFields[2] != "" {
			if amount, err := strconv.ParseInt(m.modalFields[2], 10, 64); err == nil {
				m.transactions = append(m.transactions, Transaction{
					From:      m.modalFields[0],
					To:        m.modalFields[1],
					Amount:    int(amount),
					Timestamp: time.Now(),
				})
				saveTransactions(m.transactions)
			}
		}
	}
	m.modalMode = NoModal
	m.modalFields = []string{}
	m.modalFieldIndex = 0
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	var cardWidth = m.width / 2
	var cardHeight = m.height - 3

	accountsCard := m.renderAccountsCard(cardWidth, cardHeight)
	transactionsCard := m.renderTransactionsCard(cardWidth, cardHeight)

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, accountsCard, transactionsCard)

	footer := m.renderFooter()

	content := lipgloss.JoinVertical(lipgloss.Left, mainContent, footer)

	if m.modalMode != NoModal {
		modal := m.renderModal()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
	}

	return content
}

func (m model) renderAccountsCard(width, height int) string {
	title := "Accounts"
	if m.activeTab == AccountsTab {
		title = "▶ " + title + " ◀"
	}

	var content string
	for _, acc := range m.accounts {
		content += fmt.Sprintf("%s: %s\n", acc.Name, formatCents(acc.Balance))
	}

	if content == "" {
		content = "[no accounts]"
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Width(width - 2).
		Height(height)

	if m.activeTab == AccountsTab {
		style = style.BorderForeground(lipgloss.Color("33"))
	}

	return style.Render(lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(title),
		content,
	))
}

func (m model) renderTransactionsCard(width, height int) string {
	title := "Transactions"
	if m.activeTab == TransactionsTab {
		title = "▶ " + title + " ◀"
	}

	var content string
	for _, trans := range m.transactions {
		content += fmt.Sprintf("%s → %s: %s\n", trans.From, trans.To, formatCents(trans.Amount))
	}

	if content == "" {
		content = "[no transactions]"
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Width(width - 2).
		Height(height)

	if m.activeTab == TransactionsTab {
		style = style.BorderForeground(lipgloss.Color("33"))
	}

	return style.Render(lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(title),
		content,
	))
}

func (m model) renderModal() string {
	var title string
	var fieldCount int

	switch m.modalMode {
	case AccountModal:
		title = "New Account"
		fieldCount = 2
	case TransactionModal:
		title = "New Transaction"
		fieldCount = 3
	default:
		return ""
	}

	var fieldLines []string
	labels := []string{"Name", "Balance"}
	if m.modalMode == TransactionModal {
		labels = []string{"From", "To", "Amount"}
	}

	for i := range fieldCount {
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

func (m model) renderFooter() string {
	footer := "Tab: switch | Space+A: new account | Space+Space: new transaction | Q: quit"
	return lipgloss.NewStyle().
		Faint(true).
		Render(footer)
}

func main() {
	m, err := newModel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
