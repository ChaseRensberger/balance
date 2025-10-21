package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Tab int

const (
	AccountsTab Tab = iota
	TransactionsTab
)

type PrimaryModel struct {
	accounts      AccountsModel
	transactions  TransactionsModel
	activeTab     Tab
	lastSpaceTime time.Time
	spacePressed  bool
	width         int
	height        int
}

func NewPrimaryModel() (PrimaryModel, error) {
	accounts, err := NewAccountsModel()
	if err != nil {
		return PrimaryModel{}, err
	}

	transactions, err := NewTransactionsModel()
	if err != nil {
		return PrimaryModel{}, err
	}

	return PrimaryModel{
		accounts:      accounts,
		transactions:  transactions,
		activeTab:     AccountsTab,
		lastSpaceTime: time.Now().Add(-time.Second),
	}, nil
}

func (m PrimaryModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m PrimaryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		cardWidth := m.width / 2
		cardHeight := m.height - 3

		m.accounts = m.accounts.SetSize(cardWidth, cardHeight)
		m.transactions = m.transactions.SetSize(cardWidth, cardHeight)

	case tea.KeyMsg:
		if m.accounts.IsShowingModal() {
			var handled bool
			m.accounts, handled, _ = m.accounts.HandleModalInput(msg)
			if handled {
				return m, nil
			}
			return m, nil
		}

		if m.transactions.IsShowingModal() {
			var handled bool
			m.transactions, handled, _ = m.transactions.HandleModalInput(msg)
			if handled {
				return m, nil
			}
			return m, nil
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
			m.accounts = m.accounts.SetActive(m.activeTab == AccountsTab)
			m.transactions = m.transactions.SetActive(m.activeTab == TransactionsTab)
		case " ":
			if m.spacePressed && time.Since(m.lastSpaceTime) < 500*time.Millisecond {
				m.transactions = m.transactions.ShowModal()
				m.spacePressed = false
			} else {
				m.spacePressed = true
				m.lastSpaceTime = time.Now()
			}
		case "a":
			if m.spacePressed && time.Since(m.lastSpaceTime) < 500*time.Millisecond {
				m.accounts = m.accounts.ShowModal()
				m.spacePressed = false
			}
		default:
			m.spacePressed = false
		}
	}

	return m, nil
}

func (m PrimaryModel) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	m.accounts = m.accounts.SetActive(m.activeTab == AccountsTab)
	m.transactions = m.transactions.SetActive(m.activeTab == TransactionsTab)

	accountsCard := m.accounts.View()
	transactionsCard := m.transactions.View()

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, accountsCard, transactionsCard)

	footer := m.renderFooter()

	content := lipgloss.JoinVertical(lipgloss.Left, mainContent, footer)

	if m.accounts.IsShowingModal() {
		modal := m.accounts.RenderModal()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
	}

	if m.transactions.IsShowingModal() {
		modal := m.transactions.RenderModal()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
	}

	return content
}

func (m PrimaryModel) renderFooter() string {
	footer := "Tab: switch | Space+A: new account | Space+Space: new transaction | Q: quit"
	return lipgloss.NewStyle().
		Faint(true).
		Render(footer)
}
