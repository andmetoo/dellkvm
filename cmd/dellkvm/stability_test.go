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

func TestFailedRefreshClearsStaleInput(t *testing.T) {
	m := newTUIModel(Config{Inputs: defaultInputs()})
	m.currentCode, m.currentName, m.currentBus = "0x0f", "DisplayPort", 13
	updated, _ := m.updateCurrent(currentMsg{err: errors.New("disconnected")})
	got := updated.(tuiModel)
	if got.currentCode != "" || got.currentName != "" || got.currentBus != 0 {
		t.Fatal("stale current input retained")
	}
	for _, item := range got.list.Items() {
		if item.(inputItem).current {
			t.Fatal("stale checkmark retained")
		}
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
	_, err := switchInput(Config{Bus: 6, Inputs: defaultInputs()}, "dp")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "write\nread\nread\n" {
		t.Fatalf("unexpected monitor operations: %s", data)
	}
}
