//go:build linux || windows

package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

const tuiPollInterval = 30 * time.Second

func runTUI() error {
	cfg, err := desktopConfig()
	if err != nil {
		return err
	}

	model := newTUIModel(cfg)
	_, err = tea.NewProgram(model).Run()
	return err
}

type inputItem struct {
	input   Input
	current bool
}

func (i inputItem) Title() string {
	if i.current {
		return "● " + i.input.Name
	}
	return i.input.Name
}

func (i inputItem) Description() string {
	return fmt.Sprintf("%s · %s", i.input.ID, i.input.Code)
}

func (i inputItem) FilterValue() string {
	return i.input.ID + " " + i.input.Name + " " + i.input.Code
}

type tuiModel struct {
	cfg         Config
	localizer   appLocalizer
	list        list.Model
	width       int
	height      int
	currentCode string
	currentName string
	currentBus  int
	verified    bool
	status      string
	busy        bool
	buses       []int
}

type monitorMsg struct {
	buses []int
	err   error
}

type currentMsg struct {
	state currentState
	err   error
}

type switchMsg struct {
	result switchResult
	err    error
}

type pollTickMsg struct{}

func newTUIModel(cfg Config) tuiModel {
	l := localizerForConfig(cfg)
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("205")).
		BorderLeftForeground(lipgloss.Color("205"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("243")).
		BorderLeftForeground(lipgloss.Color("205"))

	items := itemsFromConfig(cfg, "")
	inputs := list.New(items, delegate, 80, 16)
	inputs.Title = l.T("ListTitle", nil)
	inputs.SetFilteringEnabled(false)
	inputs.SetShowStatusBar(false)
	inputs.SetShowHelp(false)
	inputs.SetShowPagination(false)
	inputs.DisableQuitKeybindings()

	model := tuiModel{
		cfg:       cfg,
		localizer: l,
		list:      inputs,
		width:     80,
		height:    24,
		status:    l.T("Refreshing", nil),
		busy:      true,
	}
	model.resizeList()
	return model
}

func (m tuiModel) Init() tea.Cmd {
	if m.cfg.Bus == 0 {
		return tea.Batch(discoverMonitorsCmd(m.localizer), pollTickCmd())
	}
	return tea.Batch(refreshCurrentCmd(m.cfg, m.localizer), pollTickCmd())
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeList()
		return m, nil
	case tea.KeyMsg:
		return m.updateKey(msg)
	case currentMsg:
		return m.updateCurrent(msg)
	case monitorMsg:
		return m.updateMonitors(msg)
	case switchMsg:
		return m.updateSwitch(msg)
	case pollTickMsg:
		if m.busy {
			return m, pollTickCmd()
		}
		m.busy = true
		m.status = m.localizer.T("Refreshing", nil)
		m.resizeList()
		if m.cfg.Bus == 0 {
			return m, tea.Batch(discoverMonitorsCmd(m.localizer), pollTickCmd())
		}
		return m, tea.Batch(refreshCurrentCmd(m.cfg, m.localizer), pollTickCmd())
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m tuiModel) View() string {
	title := m.renderTitle()
	current := m.renderCurrent()
	status := m.renderStatus()
	keys := m.renderKeys()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		current,
		status,
		m.list.View(),
		keys,
	)
}

func (m tuiModel) updateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "r":
		if m.busy {
			return m, nil
		}
		m.busy = true
		m.status = m.localizer.T("Refreshing", nil)
		m.resizeList()
		return m, refreshCurrentCmd(m.cfg, m.localizer)
	case "m":
		if m.busy {
			return m, nil
		}
		if len(m.buses) == 0 {
			m.busy = true
			m.status = "Detecting monitors…"
			return m, discoverMonitorsCmd(m.localizer)
		}
		m.selectNextMonitor()
		return m, refreshCurrentCmd(m.cfg, m.localizer)
	case "enter":
		if m.busy {
			return m, nil
		}

		item, ok := m.list.SelectedItem().(inputItem)
		if !ok {
			return m, nil
		}
		if m.cfg.Bus == 0 {
			m.status = "Select a monitor with m before switching."
			return m, nil
		}
		m.busy = true
		m.status = m.localizer.T("Switching", map[string]any{"Name": item.input.Name})
		m.resizeList()
		return m, switchInputCmd(m.cfg, item.input.ID)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m tuiModel) updateMonitors(msg monitorMsg) (tea.Model, tea.Cmd) {
	m.busy = false
	if msg.err != nil {
		m.status = m.localizer.T("ErrorPrefix", map[string]any{"Error": msg.err.Error()})
		return m, nil
	}
	m.buses = msg.buses
	if len(m.buses) == 0 {
		m.status = "No monitors found. Press m to retry detection."
		return m, nil
	}
	if m.cfg.Bus == 0 && len(m.buses) > 1 {
		m.status = fmt.Sprintf("Monitors %s found. Press m to select one.", formatBuses(m.buses))
		return m, nil
	}
	m.selectNextMonitor()
	return m, refreshCurrentCmd(m.cfg, m.localizer)
}

func (m *tuiModel) selectNextMonitor() {
	current := m.cfg.Bus
	next := m.buses[0]
	for i, bus := range m.buses {
		if bus == current {
			next = m.buses[(i+1)%len(m.buses)]
			break
		}
	}
	m.cfg.Bus = next
	m.currentCode, m.currentName, m.currentBus = "", "", 0
	m.verified = false
	m.status = fmt.Sprintf("Monitor bus %d selected. Reading current input…", next)
	m.busy = true
	m.resizeList()
}

func (m tuiModel) updateCurrent(msg currentMsg) (tea.Model, tea.Cmd) {
	m.busy = false
	if msg.err != nil {
		m.verified = false
		m.status = m.localizer.T("ErrorPrefix", map[string]any{"Error": msg.err.Error()})
		m.resizeList()
		return m, nil
	}

	m.currentCode = msg.state.Code
	m.currentBus = msg.state.Bus
	m.verified = true
	input, ok := findInputByCode(m.cfg, msg.state.Code)
	if ok {
		m.currentName = input.Name
		m.status = m.localizer.T("CurrentRefreshed", map[string]any{"Bus": msg.state.Bus})
	} else {
		m.currentName = ""
		m.status = m.localizer.T("CurrentCodeNotFound", nil)
	}
	if msg.state.AutoDetected {
		m.status = m.localizer.T("AutoDetectedPrefix", map[string]any{"Bus": msg.state.Bus, "Message": m.status})
	}
	m.resizeList()

	return m, m.list.SetItems(itemsFromConfig(m.cfg, msg.state.Code))
}

func (m tuiModel) updateSwitch(msg switchMsg) (tea.Model, tea.Cmd) {
	m.busy = false
	if msg.err != nil {
		m.status = m.localizer.T("ErrorPrefix", map[string]any{"Error": msg.err.Error()})
		m.resizeList()
		return m, nil
	}

	m.status = msg.result.Message
	m.currentCode = msg.result.Code
	m.currentBus = msg.result.Bus
	m.verified = msg.result.Verified
	m.currentName = ""
	if input, ok := findInputByCode(m.cfg, m.currentCode); ok {
		m.currentName = input.Name
	}
	m.resizeList()
	return m, m.list.SetItems(itemsFromConfig(m.cfg, m.currentCode))
}

func (m *tuiModel) resizeList() {
	width := m.contentWidth()
	fixedLines := countLines(m.renderTitle()) +
		countLines(m.renderCurrent()) +
		countLines(m.renderStatus()) +
		countLines(m.renderKeys())
	m.list.SetSize(width, max(m.height-fixedLines, 3))
}

func (m tuiModel) contentWidth() int {
	if m.width > 0 {
		return m.width
	}
	return 80
}

func (m tuiModel) renderTitle() string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Width(m.contentWidth()).
		Render("dellkvm")
}

func (m tuiModel) renderCurrent() string {
	current := m.localizer.T("TuiCurrentUnknown", nil)
	if m.currentCode != "" && !m.verified {
		name := m.currentName
		if name == "" {
			name = m.currentCode
		}
		current = m.localizer.T("TuiLastSelected", map[string]any{
			"Name": name,
			"Bus":  m.currentBus,
		})
	} else if m.currentName != "" {
		current = m.localizer.T("TuiCurrentInput", map[string]any{
			"Name": m.currentName,
			"Code": m.currentCode,
			"Bus":  m.currentBus,
		})
	} else if m.currentCode != "" {
		current = m.localizer.T("TuiCurrentCode", map[string]any{
			"Code": m.currentCode,
			"Bus":  m.currentBus,
		})
	}
	return lipgloss.NewStyle().Width(m.contentWidth()).Render(current)
}

func (m tuiModel) renderStatus() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Width(m.contentWidth()).
		Render(m.status)
}

func (m tuiModel) renderKeys() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Width(m.contentWidth()).
		Render(m.localizer.T("Keys", nil) + fmt.Sprintf("  m: monitor (bus %d)", m.cfg.Bus))
}

func countLines(s string) int {
	if s == "" {
		return 1
	}
	return strings.Count(s, "\n") + 1
}

func refreshCurrentCmd(cfg Config, l appLocalizer) tea.Cmd {
	return func() tea.Msg {
		state, err := getCurrentState(cfg, l)
		return currentMsg{state: state, err: err}
	}
}

func discoverMonitorsCmd(l appLocalizer) tea.Cmd {
	return func() tea.Msg {
		buses, err := detectBuses(l)
		return monitorMsg{buses: buses, err: err}
	}
}

func switchInputCmd(cfg Config, id string) tea.Cmd {
	return func() tea.Msg {
		result, err := switchInputDetailed(cfg, id)
		return switchMsg{result: result, err: err}
	}
}

func pollTickCmd() tea.Cmd {
	return tea.Tick(tuiPollInterval, func(time.Time) tea.Msg { return pollTickMsg{} })
}

func itemsFromConfig(cfg Config, currentCode string) []list.Item {
	items := make([]list.Item, 0, len(cfg.Inputs))
	currentCode = normalizeCode(currentCode)
	for _, input := range cfg.Inputs {
		items = append(items, inputItem{
			input:   input,
			current: currentCode != "" && normalizeCode(input.Code) == currentCode,
		})
	}
	return items
}
