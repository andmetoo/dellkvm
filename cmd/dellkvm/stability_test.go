//go:build linux

package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestValidateInputs(t *testing.T) {
	for _, tt := range []struct {
		name   string
		inputs []Input
		bad    bool
	}{
		{"normalize", []Input{{ID: " dp ", Name: "DisplayPort", Code: " 0xF "}}, false},
		{"empty id", []Input{{Name: "DP", Code: "0x0f"}}, true},
		{"empty name", []Input{{ID: "dp", Code: "0x0f"}}, true},
		{"missing prefix", []Input{{ID: "dp", Name: "DP", Code: "15"}}, true},
		{"overflow", []Input{{ID: "dp", Name: "DP", Code: "0x100"}}, true},
		{"invalid", []Input{{ID: "dp", Name: "DP", Code: "0xgg"}}, true},
		{"duplicate id", []Input{{ID: "dp", Name: "DP", Code: "0x0f"}, {ID: "dp", Name: "HDMI", Code: "0x11"}}, true},
		{"duplicate code", []Input{{ID: "dp", Name: "DP", Code: "0x0f"}, {ID: "dp2", Name: "DP2", Code: "0xF"}}, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{Inputs: tt.inputs}
			err := validateInputs(&cfg)
			if (err != nil) != tt.bad {
				t.Fatalf("error = %v, bad = %v", err, tt.bad)
			}
			if !tt.bad && (cfg.Inputs[0].ID != "dp" || cfg.Inputs[0].Code != "0x0f") {
				t.Fatalf("not normalized: %+v", cfg)
			}
		})
	}
}

func TestBusyTUIDoesNotStartAnotherOperation(t *testing.T) {
	m := newTUIModel(Config{Inputs: defaultInputs()})
	for _, key := range []tea.KeyMsg{{Type: tea.KeyRunes, Runes: []rune{'r'}}, {Type: tea.KeyEnter}} {
		_, cmd := m.updateKey(key)
		if cmd != nil {
			t.Fatal("started a command while busy")
		}
	}
}

func TestFailedRefreshRetainsLastInputAsUnverified(t *testing.T) {
	m := newTUIModel(Config{Inputs: defaultInputs()})
	m.currentCode, m.currentName, m.currentBus = "0x0f", "DisplayPort", 13
	m.verified = true
	m.list.SetItems(itemsFromConfig(m.cfg, m.currentCode))
	updated, _ := m.updateCurrent(currentMsg{err: errors.New("disconnected")})
	got := updated.(tuiModel)
	if got.currentCode != "0x0f" || got.currentName != "DisplayPort" || got.currentBus != 13 || got.verified {
		t.Fatalf("last selection lost or incorrectly verified: %+v", got)
	}
	if !strings.Contains(got.renderCurrent(), "unverified") {
		t.Fatalf("uncertainty hidden: %q", got.renderCurrent())
	}
	found := false
	for _, item := range got.list.Items() {
		if item.(inputItem).current {
			found = true
		}
	}
	if !found {
		t.Fatal("last selected input was unmarked")
	}
}

func TestTUISwitchKeepsSelectionUntilNextRead(t *testing.T) {
	m := newTUIModel(Config{Bus: 13, Inputs: defaultInputs()})
	updated, _ := m.updateSwitch(switchMsg{result: switchResult{Code: "0x0f", Bus: 13, Message: "sent"}})
	m = updated.(tuiModel)
	if m.currentCode != "0x0f" || m.verified || !strings.Contains(m.renderCurrent(), "unverified") {
		t.Fatalf("switch state = %+v", m)
	}
	updated, _ = m.updateCurrent(currentMsg{state: currentState{Code: "0x11", Bus: 13}})
	m = updated.(tuiModel)
	if m.currentCode != "0x11" || !m.verified {
		t.Fatalf("read did not correct selection: %+v", m)
	}
}

func TestTUIPeriodicRefreshSkipsBusyMonitor(t *testing.T) {
	m := newTUIModel(Config{Bus: 13, Inputs: defaultInputs()})
	updated, cmd := m.Update(pollTickMsg{})
	if !updated.(tuiModel).busy || cmd == nil {
		t.Fatal("busy TUI did not schedule next poll")
	}
	m.busy = false
	updated, cmd = m.Update(pollTickMsg{})
	if !updated.(tuiModel).busy || cmd == nil {
		t.Fatal("idle TUI did not schedule a refresh")
	}
}

func TestTraySelectionTracksVerificationAndBus(t *testing.T) {
	selection := inputSelection{}.onSwitch(switchResult{Code: "0x0f", Bus: 13})
	selection = selection.onRead(currentState{}, errors.New("disconnected"))
	if selection.code != "0x0f" || selection.bus != 13 || selection.verified {
		t.Fatalf("unverified selection lost: %+v", selection)
	}
	selection = selection.onRead(currentState{Code: "0x11", Bus: 13}, nil)
	if selection.code != "0x11" || !selection.verified {
		t.Fatalf("read did not correct selection: %+v", selection)
	}
	if got := selection.forBus(6); got.code != "" {
		t.Fatalf("selection leaked to another monitor: %+v", got)
	}
}

func TestEditableConfigDoesNotOverwrite(t *testing.T) {
	t.Chdir(t.TempDir())
	data := []byte("invalid config deliberately preserved")
	if err := os.WriteFile("config.toml", data, 0o600); err != nil {
		t.Fatal(err)
	}
	path, err := editableConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil || string(got) != string(data) {
		t.Fatalf("config changed: %q, %v", got, err)
	}
}

func TestVCPRejectsOverflow(t *testing.T) {
	if _, err := parseVCPCode("sl=0x100"); err == nil {
		t.Fatal("accepted value outside input source byte")
	}
}

func TestAutoSwitchRefusesMultipleMonitors(t *testing.T) {
	withFakeDDC(t, "multiple-monitors")
	_, err := switchInput(Config{Inputs: defaultInputs()}, "dp")
	if err == nil || !strings.Contains(err.Error(), "multiple monitors") {
		t.Fatalf("unexpected result: %v", err)
	}
}

func TestSwitchWaitsForMonitorWithoutRepeatingWrite(t *testing.T) {
	withFakeDDC(t, "delayed-verification")
	withoutSwitchDelay(t)
	path := filepath.Join(t.TempDir(), "events")
	t.Setenv("DDC_TEST_EVENTS", path)
	result, err := switchInputDetailed(Config{Bus: 6, Inputs: defaultInputs()}, "dp")
	if err != nil {
		t.Fatal(err)
	}
	if !result.Verified || result.Bus != 6 || result.Code != "0x0f" {
		t.Fatalf("verified switch state = %+v", result)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "write\nread\nread\n" {
		t.Fatalf("unexpected monitor operations: %s", data)
	}
}
