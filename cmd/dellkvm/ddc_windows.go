//go:build windows

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32                      = syscall.NewLazyDLL("user32.dll")
	dxva2                       = syscall.NewLazyDLL("dxva2.dll")
	enumDisplayMonitors         = user32.NewProc("EnumDisplayMonitors")
	getNumberOfPhysicalMonitors = dxva2.NewProc("GetNumberOfPhysicalMonitorsFromHMONITOR")
	getPhysicalMonitors         = dxva2.NewProc("GetPhysicalMonitorsFromHMONITOR")
	destroyPhysicalMonitors     = dxva2.NewProc("DestroyPhysicalMonitors")
	getVCPFeature               = dxva2.NewProc("GetVCPFeatureAndVCPFeatureReply")
	setVCPFeature               = dxva2.NewProc("SetVCPFeature")
	execMonitorCommand          = exec.CommandContext
)

type physicalMonitor struct {
	handle      uintptr
	description [128]uint16
}

type monitorHandle struct {
	id     int
	handle uintptr
	name   string
}

func runDDC(l appLocalizer, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ddcContext, ddcTimeout)
	defer cancel()
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	cmd := execMonitorCommand(ctx, executable, append([]string{"--monitor-worker"}, args...)...)
	cmd.WaitDelay = time.Second
	out, err := cmd.CombinedOutput()
	raw := strings.TrimSpace(string(out))
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return "", userError{err: errDDCTimeout, message: fmt.Sprintf("monitor request timed out after %s", ddcTimeout)}
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return "", context.Canceled
	}
	if err != nil {
		if strings.Contains(strings.ToLower(raw), "no monitor") {
			return "", userError{err: errNoMonitor, message: noMonitorMessage(l, args, raw)}
		}
		if raw == "" {
			raw = err.Error()
		}
		return "", fmt.Errorf("Windows monitor request: %s", raw)
	}
	return string(out), nil
}

func runPlatformWorker(args []string) (bool, error) {
	if len(args) == 0 || args[0] != "--monitor-worker" {
		return false, nil
	}
	operation, bus, value, err := parseMonitorWorkerArgs(args[1:])
	if err != nil {
		return true, err
	}
	monitors, release, err := enumeratePhysicalMonitors()
	if err != nil {
		return true, err
	}
	defer release()
	switch operation {
	case "detect":
		for _, monitor := range monitors {
			fmt.Printf("Display %d\n   Monitor: %d\n   Description: %s\n", monitor.id, monitor.id, monitor.name)
		}
		return true, nil
	case "get":
		monitor, err := selectPhysicalMonitor(monitors, bus)
		if err != nil {
			return true, err
		}
		var kind, current, maximum uint32
		ok, _, callErr := getVCPFeature.Call(monitor.handle, 0x60, uintptr(unsafe.Pointer(&kind)), uintptr(unsafe.Pointer(&current)), uintptr(unsafe.Pointer(&maximum)))
		if ok == 0 {
			return true, winCallError("GetVCPFeatureAndVCPFeatureReply", callErr)
		}
		if current > 0xff {
			return true, fmt.Errorf("monitor %d returned invalid input code %d", bus, current)
		}
		fmt.Printf("VCP code 0x60 (Input Source): sl=0x%02x\n", current)
		return true, nil
	case "set":
		monitor, err := selectPhysicalMonitor(monitors, bus)
		if err != nil {
			return true, err
		}
		ok, _, callErr := setVCPFeature.Call(monitor.handle, 0x60, uintptr(value))
		if ok == 0 {
			return true, winCallError("SetVCPFeature", callErr)
		}
		return true, nil
	}
	return true, errors.New("unsupported monitor operation")
}

func parseMonitorWorkerArgs(args []string) (string, int, uint8, error) {
	if len(args) == 2 && args[0] == "detect" && args[1] == "--brief" {
		return "detect", 0, 0, nil
	}
	if len(args) < 4 || args[0] != "--bus" {
		return "", 0, 0, errors.New("invalid monitor worker arguments")
	}
	bus, err := strconv.Atoi(args[1])
	if err != nil || bus <= 0 {
		return "", 0, 0, errors.New("invalid monitor index")
	}
	if len(args) == 4 && args[2] == "getvcp" && args[3] == vcpInputSource {
		return "get", bus, 0, nil
	}
	if len(args) == 6 && args[2] == "--noverify" && args[3] == "setvcp" && args[4] == vcpInputSource {
		code := normalizeCode(args[5])
		if !inputCodePattern.MatchString(code) {
			return "", 0, 0, errors.New("invalid input code")
		}
		parsed, _ := strconv.ParseUint(strings.TrimPrefix(code, "0x"), 16, 8)
		return "set", bus, uint8(parsed), nil
	}
	return "", 0, 0, errors.New("invalid monitor worker arguments")
}

func enumeratePhysicalMonitors() ([]monitorHandle, func(), error) {
	var logical []uintptr
	callback := syscall.NewCallback(func(hmonitor, hdc, rectangle, data uintptr) uintptr {
		logical = append(logical, hmonitor)
		return 1
	})
	ok, _, callErr := enumDisplayMonitors.Call(0, 0, callback, 0)
	if ok == 0 {
		return nil, func() {}, winCallError("EnumDisplayMonitors", callErr)
	}
	var arrays [][]physicalMonitor
	release := func() {
		for _, array := range arrays {
			_, _, _ = destroyPhysicalMonitors.Call(uintptr(len(array)), uintptr(unsafe.Pointer(&array[0])))
		}
	}
	var monitors []monitorHandle
	for _, logicalHandle := range logical {
		var count uint32
		ok, _, _ := getNumberOfPhysicalMonitors.Call(logicalHandle, uintptr(unsafe.Pointer(&count)))
		if ok == 0 || count == 0 {
			continue
		}
		if count > 64 {
			release()
			return nil, func() {}, fmt.Errorf("unexpected physical monitor count %d", count)
		}
		array := make([]physicalMonitor, count)
		ok, _, _ = getPhysicalMonitors.Call(logicalHandle, uintptr(count), uintptr(unsafe.Pointer(&array[0])))
		if ok == 0 {
			release()
			return nil, func() {}, fmt.Errorf("cannot enumerate physical monitors")
		}
		arrays = append(arrays, array)
		for _, physical := range array {
			monitors = append(monitors, monitorHandle{id: len(monitors) + 1, handle: physical.handle, name: syscall.UTF16ToString(physical.description[:])})
		}
	}
	return monitors, release, nil
}

func selectPhysicalMonitor(monitors []monitorHandle, bus int) (monitorHandle, error) {
	for _, monitor := range monitors {
		if monitor.id == bus {
			return monitor, nil
		}
	}
	return monitorHandle{}, fmt.Errorf("No monitor detected on bus %d", bus)
}

func winCallError(operation string, err error) error {
	if err == syscall.Errno(0) {
		return fmt.Errorf("%s failed; check DDC/CI in the monitor OSD", operation)
	}
	return fmt.Errorf("%s failed: %w", operation, err)
}
