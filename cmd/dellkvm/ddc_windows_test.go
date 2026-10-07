//go:build windows

package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

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
