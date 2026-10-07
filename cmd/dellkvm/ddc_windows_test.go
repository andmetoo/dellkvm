//go:build windows

package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWindowsInputCode(t *testing.T) {
	for _, tt := range []struct {
		value uint32
		want  uint8
		bad   bool
	}{
		{0x0f, 0x0f, false},
		{0x0f0f, 0x0f, false},
		{0x1111, 0x11, false},
		{0x0f10, 0, true},
		{0x10000, 0, true},
	} {
		got, err := windowsInputCode(tt.value)
		if got != tt.want || (err != nil) != tt.bad {
			t.Errorf("windowsInputCode(0x%x) = 0x%x, %v", tt.value, got, err)
		}
	}
}

func TestWindowsSwitchDoesNotRequireInitialRead(t *testing.T) {
	path := filepath.Join(t.TempDir(), "set-was-called")
	oldCommand, oldDelay := execMonitorCommand, verifyDelay
	execMonitorCommand = func(ctx context.Context, _ string, args ...string) *exec.Cmd {
		cmd := exec.CommandContext(ctx, os.Args[0], append([]string{"-test.run=TestWindowsSwitchHelperProcess", "--"}, args...)...)
		cmd.Env = append(os.Environ(), "DDC_TEST_SWITCH="+path)
		return cmd
	}
	verifyDelay = 0
	t.Cleanup(func() { execMonitorCommand, verifyDelay = oldCommand, oldDelay })

	message, err := switchInput(Config{Inputs: []Input{{ID: "dp", Name: "DisplayPort", Code: "0x0f"}}}, "dp")
	if err != nil {
		t.Fatalf("switchInput: %v", err)
	}
	if !strings.Contains(message, "Auto-detected bus 1") || !strings.Contains(message, "0x0f") {
		t.Fatalf("switchInput message = %q", message)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("set was not called: %v", err)
	}
}

func TestWindowsSwitchHelperProcess(t *testing.T) {
	path := os.Getenv("DDC_TEST_SWITCH")
	if path == "" {
		return
	}
	var args []string
	for i, arg := range os.Args {
		if arg == "--" {
			args = os.Args[i+1:]
			break
		}
	}
	if len(args) > 0 && args[0] == "--monitor-worker" {
		args = args[1:]
	}
	switch {
	case len(args) == 2 && args[0] == "detect":
		_, _ = os.Stdout.WriteString("Display 1\n   Monitor: 1\n")
	case len(args) == 6 && args[3] == "setvcp":
		if err := os.WriteFile(path, []byte("set"), 0o600); err != nil {
			os.Exit(1)
		}
	case len(args) == 4 && args[2] == "getvcp":
		if _, err := os.Stat(path); err != nil {
			_, _ = os.Stderr.WriteString("invalid input code 3855\n")
			os.Exit(1)
		}
		_, _ = os.Stdout.WriteString("VCP code 0x60 (Input Source): sl=0x0f\n")
	default:
		os.Exit(1)
	}
	os.Exit(0)
}

func TestParseMonitorWorkerArgs(t *testing.T) {
	for _, tt := range []struct {
		args      []string
		operation string
		bus       int
		value     uint8
		bad       bool
	}{
		{[]string{"detect", "--brief"}, "detect", 0, 0, false},
		{[]string{"--bus", "2", "getvcp", "60"}, "get", 2, 0, false},
		{[]string{"--bus", "2", "--noverify", "setvcp", "60", "0x0F"}, "set", 2, 15, false},
		{[]string{"--bus", "0", "getvcp", "60"}, "", 0, 0, true},
		{[]string{"--bus", "-1", "getvcp", "60"}, "", 0, 0, true},
		{[]string{"--bus", "2", "--noverify", "setvcp", "60", "0x100"}, "", 0, 0, true},
		{[]string{"--bus", "2", "setvcp", "60", "0x11"}, "", 0, 0, true},
	} {
		operation, bus, value, err := parseMonitorWorkerArgs(tt.args)
		if (err != nil) != tt.bad || (!tt.bad && (operation != tt.operation || bus != tt.bus || value != tt.value)) {
			t.Errorf("args %v: %q, %d, %d, %v", tt.args, operation, bus, value, err)
		}
	}
}

func TestWindowsMonitorSelection(t *testing.T) {
	monitors := []monitorHandle{{id: 1, name: "first"}, {id: 2, name: "second"}}
	got, err := selectPhysicalMonitor(monitors, 2)
	if err != nil || got.name != "second" {
		t.Fatalf("selected %v, %v", got, err)
	}
	if _, err := selectPhysicalMonitor(monitors, 3); err == nil || !strings.Contains(err.Error(), "No monitor") {
		t.Fatalf("missing monitor: %v", err)
	}
}

func TestWindowsDetectParser(t *testing.T) {
	raw := "Display 1\n   Monitor: 1\n   Description: Dell\nDisplay 2\n   Monitor: 2\n   Description: LG\n"
	buses := parseDDCBuses(raw)
	if len(buses) != 2 || buses[0] != 1 || buses[1] != 2 {
		t.Fatalf("parsed monitors: %v", buses)
	}
}

func TestWindowsWorkerRejectsInvalidArgsWithoutHardware(t *testing.T) {
	done, err := runPlatformWorker([]string{"--monitor-worker", "--bus", "0", "getvcp", "60"})
	if !done || err == nil {
		t.Fatalf("unexpected result: %t, %v", done, err)
	}
	if errors.Is(err, errNoMonitor) {
		t.Fatalf("accessed hardware: %v", err)
	}
}

func TestWindowsMonitorRequestTimeout(t *testing.T) {
	oldCommand, oldTimeout := execMonitorCommand, ddcTimeout
	ddcTimeout = 50 * time.Millisecond
	execMonitorCommand = func(ctx context.Context, _ string, _ ...string) *exec.Cmd {
		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestWindowsSleepingWorker")
		cmd.Env = append(os.Environ(), "DDC_TEST_SLEEP=1")
		return cmd
	}
	t.Cleanup(func() { execMonitorCommand, ddcTimeout = oldCommand, oldTimeout })
	_, err := runDDC(defaultLocalizer(), "detect", "--brief")
	if !errors.Is(err, errDDCTimeout) {
		t.Fatalf("timeout error = %v", err)
	}
}

func TestWindowsSleepingWorker(t *testing.T) {
	if os.Getenv("DDC_TEST_SLEEP") != "1" {
		return
	}
	time.Sleep(5 * time.Second)
	os.Exit(0)
}
