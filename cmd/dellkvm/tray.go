//go:build linux || windows

package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"

	"fyne.io/systray"
	"github.com/pelletier/go-toml/v2"
)

// A single worker owns the menu and all monitor operations. Clicks while it
// is busy are discarded, never queued for execution after a switch.
type trayApp struct {
	cfg                      Config
	busy                     atomic.Bool
	done                     chan struct{}
	jobs                     chan func()
	status, monitors, inputs *systray.MenuItem
	configLocation           *systray.MenuItem
	monitorItems, inputItems []*systray.MenuItem
	inputCodes               []string
	inputStop, monitorStop   chan struct{}
	workerDone               chan struct{}
}

func runTray() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ddcContext = ctx
	defer func() { ddcContext = context.Background() }()
	app := &trayApp{done: make(chan struct{}), jobs: make(chan func()), workerDone: make(chan struct{})}
	systray.Run(app.ready, func() {
		close(app.done)
		cancel()
		<-app.workerDone
	})
	return nil
}

func (a *trayApp) ready() {
	systray.SetIcon(trayIcon())
	systray.SetTooltip("dellkvm")
	a.status = systray.AddMenuItem("dellkvm", "")
	a.status.Disable()
	a.monitors = systray.AddMenuItem("Monitor", "")
	a.inputs = systray.AddMenuItem("Switch input", "")
	a.bind(systray.AddMenuItem("Refresh / reconnect", ""), a.refresh)
	a.bind(systray.AddMenuItem("Reload configuration", ""), a.reload)
	a.configLocation = systray.AddMenuItem("Config", "")
	a.configLocation.Disable()
	a.bind(systray.AddMenuItem("Edit configuration…", ""), a.editConfig)
	a.bind(systray.AddMenuItem("Open status log…", ""), a.openLog)
	systray.AddSeparator()
	quit := systray.AddMenuItem("Quit", "")
	go func() {
		select {
		case <-quit.ClickedCh:
			systray.Quit()
		case <-a.done:
		}
	}()
	a.busy.Store(true)
	go func() {
		defer close(a.workerDone)
		a.reload()
		a.busy.Store(false)
		for {
			select {
			case job := <-a.jobs:
				a.inputs.Disable()
				a.monitors.Disable()
				job()
				select {
				case <-a.done:
					return
				default:
				}
				if len(a.cfg.Inputs) > 0 {
					a.inputs.Enable()
				} else {
					a.inputs.Disable()
				}
				a.monitors.Enable()
				a.busy.Store(false)
			case <-a.done:
				return
			}
		}
	}()
}

func (a *trayApp) bind(item *systray.MenuItem, job func()) {
	a.bindUntil(item, job, nil)
}

func (a *trayApp) bindUntil(item *systray.MenuItem, job func(), stop <-chan struct{}) {
	go func() {
		for {
			select {
			case <-stop:
				return
			case _, ok := <-item.ClickedCh:
				if !ok {
					return
				}
				if a.busy.CompareAndSwap(false, true) {
					select {
					case a.jobs <- job:
					case <-a.done:
						return
					}
				}
			case <-a.done:
				return
			}
		}
	}()
}

func (a *trayApp) report(message string) {
	select {
	case <-a.done:
		return
	default:
	}
	short := strings.Join(strings.Fields(message), " ")
	if len([]rune(short)) > 100 {
		short = string([]rune(short)[:97]) + "…"
	}
	a.status.SetTitle(short)
	systray.SetTooltip("dellkvm · " + short)
	// Keep the complete last result available even when a tray host truncates it.
	if path, err := statusLogPath(); err == nil {
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err == nil {
			_ = os.WriteFile(path, []byte(message+"\n"), 0o600)
		}
	}
}

func (a *trayApp) reload() {
	cfg, path, err := startupConfig()
	if path != "" {
		a.configLocation.SetTitle("Config: " + path)
	}
	if err != nil {
		a.cfg = Config{}
		a.markCurrent("")
		a.inputs.Disable()
		a.report(err.Error())
		return
	}
	a.cfg = cfg
	if a.inputStop != nil {
		close(a.inputStop)
	}
	a.inputStop = make(chan struct{})
	for _, item := range a.inputItems {
		item.Remove()
	}
	a.inputItems, a.inputCodes = nil, nil
	for _, input := range cfg.Inputs {
		item := a.inputs.AddSubMenuItemCheckbox(input.Name+" · "+input.Code, "", false)
		a.inputItems = append(a.inputItems, item)
		a.inputCodes = append(a.inputCodes, input.Code)
		a.bindUntil(item, func() {
			a.markCurrent("")
			a.report("Switching to " + input.Name + "…")
			message, err := switchInput(a.cfg, input.ID)
			if err != nil {
				a.report(err.Error())
				return
			}
			a.report(message)
		}, a.inputStop)
	}
	a.refresh()
}

func desktopConfig() (Config, error) {
	if _, err := loadConfig(); errors.Is(err, errConfigNotFound) {
		return Config{Language: defaultLanguage, Inputs: defaultInputs()}, nil
	}
	return ensureConfig()
}

// Interactive entry points persist a starter config before touching hardware.
// Existing files keep their normal precedence and are never replaced, even
// when invalid. Read-only status commands continue to use desktopConfig.
func startupConfig() (Config, string, error) {
	path, err := editableConfigPath()
	if err != nil {
		return Config{}, path, fmt.Errorf("prepare configuration: %w", err)
	}
	cfg, err := ensureConfig()
	if err != nil {
		return Config{}, path, fmt.Errorf("%s: %w", path, err)
	}
	return cfg, path, nil
}

func (a *trayApp) refresh() {
	a.markCurrent("")
	a.report("Detecting monitors…")
	l := localizerForConfig(a.cfg)
	buses, err := detectBuses(l)
	if a.monitorStop != nil {
		close(a.monitorStop)
	}
	a.monitorStop = make(chan struct{})
	for _, item := range a.monitorItems {
		item.Remove()
	}
	a.monitorItems = nil
	if err != nil {
		a.report(err.Error())
		return
	}
	for _, bus := range append([]int{0}, buses...) {
		title := fmt.Sprintf("Monitor %d", bus)
		if bus == 0 {
			title = "Auto (single monitor)"
		}
		item := a.monitors.AddSubMenuItemCheckbox(title, "", a.cfg.Bus == bus)
		a.monitorItems = append(a.monitorItems, item)
		a.bindUntil(item, func() { a.cfg.Bus = bus; a.refresh() }, a.monitorStop)
	}
	state, err := getCurrentState(a.cfg, l)
	if err != nil {
		a.report(err.Error())
		return
	}
	a.markCurrent(state.Code)
	name := state.Code
	if input, ok := findInputByCode(a.cfg, state.Code); ok {
		name = input.Name
	}
	a.report(fmt.Sprintf("Monitor %d · %s", state.Bus, name))
}

func (a *trayApp) markCurrent(code string) {
	for i, item := range a.inputItems {
		if a.inputCodes[i] == code {
			item.Check()
		} else {
			item.Uncheck()
		}
	}
}

func (a *trayApp) editConfig() {
	path, err := editableConfigPath()
	if err == nil {
		err = openFile(path)
	}
	if err != nil {
		a.report(err.Error())
	}
}

func editableConfigPath() (string, error) {
	paths, err := configPaths()
	if err != nil {
		return "", err
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return filepath.Abs(path)
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
	}
	path := paths[len(paths)-1]
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", err
	}
	data, err := toml.Marshal(Config{Language: defaultLanguage, Inputs: defaultInputs()})
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if errors.Is(err, os.ErrExist) {
		return path, nil
	}
	if err != nil {
		return "", err
	}
	_, writeErr := f.Write(data)
	return path, errors.Join(writeErr, f.Close())
}

func statusLogPath() (string, error) {
	dir, err := os.UserCacheDir()
	return filepath.Join(dir, "dellkvm", "status.log"), err
}

func (a *trayApp) openLog() {
	path, err := statusLogPath()
	if err == nil {
		err = openFile(path)
	}
	if err != nil {
		a.report(err.Error())
	}
}

// Draw a small monitor icon without shipping an external asset dependency.
func trayIcon() []byte {
	im := image.NewNRGBA(image.Rect(0, 0, 32, 32))
	for y := 5; y < 24; y++ {
		for x := 2; x < 30; x++ {
			c := color.NRGBA{R: 90, G: 190, B: 240, A: 255}
			if x > 4 && x < 27 && y > 7 && y < 20 {
				c = color.NRGBA{R: 24, G: 35, B: 48, A: 255}
			}
			im.SetNRGBA(x, y, c)
		}
	}
	for y := 24; y < 28; y++ {
		for x := 12; x < 20; x++ {
			im.SetNRGBA(x, y, color.NRGBA{R: 90, G: 190, B: 240, A: 255})
		}
	}
	var pngData bytes.Buffer
	_ = png.Encode(&pngData, im)
	if runtime.GOOS == "windows" {
		var ico bytes.Buffer
		for _, value := range []uint16{0, 1, 1} {
			_ = binary.Write(&ico, binary.LittleEndian, value)
		}
		ico.Write([]byte{32, 32, 0, 0})
		for _, value := range []uint16{1, 32} {
			_ = binary.Write(&ico, binary.LittleEndian, value)
		}
		for _, value := range []uint32{uint32(pngData.Len()), 22} {
			_ = binary.Write(&ico, binary.LittleEndian, value)
		}
		ico.Write(pngData.Bytes())
		return ico.Bytes()
	}
	return pngData.Bytes()
}
