//go:build linux

package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTUIMultipleMonitorsRequiresSelection(t *testing.T) {
	m := newTUIModel(Config{Inputs: defaultInputs()})
	updated, cmd := m.updateMonitors(monitorMsg{buses: []int{6, 13}})
	m = updated.(tuiModel)
	if cmd != nil || m.cfg.Bus != 0 || m.busy {
		t.Fatalf("selected a bus implicitly: %+v", m)
	}
	updated, cmd = m.updateKey(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(tuiModel)
	if cmd != nil || !strings.Contains(m.status, "Select a monitor") {
		t.Fatalf("switch not blocked: %q", m.status)
	}
	updated, cmd = m.updateKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	m = updated.(tuiModel)
	if cmd == nil || m.cfg.Bus != 6 {
		t.Fatalf("first monitor not selected: %+v", m)
	}
	m.busy = false
	updated, cmd = m.updateKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	m = updated.(tuiModel)
	if cmd == nil || m.cfg.Bus != 13 {
		t.Fatalf("second monitor not selected: %+v", m)
	}
	m.busy = false
	updated, _ = m.updateKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	if updated.(tuiModel).cfg.Bus != 6 {
		t.Fatal("monitor selection did not wrap")
	}
}

func TestTUISingleMonitorPinsBus(t *testing.T) {
	m := newTUIModel(Config{Inputs: defaultInputs()})
	updated, cmd := m.updateMonitors(monitorMsg{buses: []int{13}})
	m = updated.(tuiModel)
	if cmd == nil || m.cfg.Bus != 13 || !m.busy {
		t.Fatalf("single monitor not pinned: %+v", m)
	}
	updated, _ = m.updateCurrent(currentMsg{state: currentState{Bus: 13, Code: "0x0f"}})
	m = updated.(tuiModel)
	if m.busy || m.cfg.Bus != 13 {
		t.Fatalf("lost selected bus after refresh: %+v", m)
	}
	updated, cmd = m.updateKey(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil || updated.(tuiModel).cfg.Bus != 13 {
		t.Fatal("switch command missing pinned bus")
	}
}
