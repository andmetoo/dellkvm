//go:build linux

package main

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func TestParseVCPCode(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "lowercase",
			raw:  "VCP code 0x60 (Input Source): current value = 0x1b, sl=0x1b",
			want: "0x1b",
		},
		{
			name: "uppercase padded",
			raw:  "VCP code 0x60 (Input Source): sl=0x0F",
			want: "0x0f",
		},
		{
			name: "unpadded",
			raw:  "VCP code 0x60 (Input Source): sl=0xf",
			want: "0x0f",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVCPCode(tt.raw)
			if err != nil {
				t.Fatalf("parseVCPCode() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("parseVCPCode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseVCPCodeMissingValue(t *testing.T) {
	raw := "VCP code 0x60 (Input Source): current value unavailable"

	_, err := parseVCPCode(raw)
	if err == nil {
		t.Fatal("parseVCPCode() error = nil, want error")
	}
	if !strings.Contains(err.Error(), raw) {
		t.Fatalf("parseVCPCode() error = %q, want raw output", err)
	}
}

func TestLoadConfigPrefersCurrentDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cwd := t.TempDir()
	t.Chdir(cwd)

	localPath := filepath.Join(cwd, "config.toml")
	userPath := filepath.Join(home, ".config", "dellkvm", "config.toml")
	if err := os.MkdirAll(filepath.Dir(userPath), 0o700); err != nil {
		t.Fatalf("mkdir user config dir: %v", err)
	}

	localConfig := `bus = 6

[[inputs]]
id = "dp"
name = "DisplayPort"
code = "0x0f"
`
	userConfig := `bus = 9

[[inputs]]
id = "hdmi"
name = "HDMI"
code = "0x11"
`
	if err := os.WriteFile(localPath, []byte(localConfig), 0o600); err != nil {
		t.Fatalf("write local config: %v", err)
	}
	if err := os.WriteFile(userPath, []byte(userConfig), 0o600); err != nil {
		t.Fatalf("write user config: %v", err)
	}

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}
	if cfg.Bus != 6 {
		t.Fatalf("loadConfig() bus = %d, want 6", cfg.Bus)
	}
}

func TestLoadConfigDefaultsLanguageToEnglish(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Chdir(t.TempDir())

	config := `bus = 6

[[inputs]]
id = "dp"
name = "DisplayPort"
code = "0x0f"
`
	if err := os.WriteFile("config.toml", []byte(config), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}
	if cfg.Language != "en" {
		t.Fatalf("loadConfig() language = %q, want en", cfg.Language)
	}
}

func TestRunHelp(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "help command", args: []string{"help"}},
		{name: "long flag", args: []string{"--help"}},
		{name: "short flag", args: []string{"-h"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := captureStdout(t, func() error {
				return run(tt.args)
			})
			if err != nil {
				t.Fatalf("run(%v) error = %v, want nil", tt.args, err)
			}
			for _, want := range []string{
				"Usage: dellkvm [command]",
				"detect",
				"current",
				"switch <id>",
				"help",
				"dellkvm init",
			} {
				if !strings.Contains(out, want) {
					t.Fatalf("run(%v) output = %q, want %q", tt.args, out, want)
				}
			}
		})
	}
}

func TestRunLearnIsUsageErrorAndDoesNotTouchConfig(t *testing.T) {
	t.Run("does not create config", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		cwd := t.TempDir()
		t.Chdir(cwd)

		err := run([]string{"learn"})
		if err == nil {
			t.Fatal("run([learn]) error = nil, want usage error")
		}
		var usage usageError
		if !errors.As(err, &usage) {
			t.Fatalf("run([learn]) error = %v, want usageError", err)
		}

		for _, path := range []string{
			filepath.Join(cwd, "config.toml"),
			filepath.Join(home, ".config", "dellkvm", "config.toml"),
		} {
			if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
				t.Fatalf("config path %s exists or stat failed with %v, want not exist", path, statErr)
			}
		}
	})

	t.Run("does not modify existing configs", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		cwd := t.TempDir()
		t.Chdir(cwd)

		localPath := filepath.Join(cwd, "config.toml")
		userPath := filepath.Join(home, ".config", "dellkvm", "config.toml")
		if err := os.MkdirAll(filepath.Dir(userPath), 0o700); err != nil {
			t.Fatalf("mkdir user config dir: %v", err)
		}

		localConfig := "bus = 6\n"
		userConfig := "bus = 9\n"
		if err := os.WriteFile(localPath, []byte(localConfig), 0o600); err != nil {
			t.Fatalf("write local config: %v", err)
		}
		if err := os.WriteFile(userPath, []byte(userConfig), 0o600); err != nil {
			t.Fatalf("write user config: %v", err)
		}

		err := run([]string{"learn"})
		if err == nil {
			t.Fatal("run([learn]) error = nil, want usage error")
		}
		var usage usageError
		if !errors.As(err, &usage) {
			t.Fatalf("run([learn]) error = %v, want usageError", err)
		}

		assertFileContent(t, localPath, localConfig)
		assertFileContent(t, userPath, userConfig)
	})
}

func TestSupportedLanguagesLoad(t *testing.T) {
	for _, lang := range supportedLanguages {
		t.Run(lang, func(t *testing.T) {
			l, err := newLocalizer(lang)
			if err != nil {
				t.Fatalf("newLocalizer(%q) error = %v", lang, err)
			}
			if got := l.T("ListTitle", nil); got == "" || got == "ListTitle" {
				t.Fatalf("ListTitle for %q = %q", lang, got)
			}
		})
	}
}

func TestLocaleDataMatchesLocaleFiles(t *testing.T) {
	for _, lang := range supportedLanguages {
		t.Run(lang, func(t *testing.T) {
			path := filepath.Join("locales", "active."+lang+".toml")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			if localeData[lang] != string(data) {
				t.Fatalf("localeData[%q] is out of sync with %s", lang, path)
			}
		})
	}
}

func TestUnsupportedLanguageRejected(t *testing.T) {
	_, err := normalizeLanguage("es")
	if err == nil {
		t.Fatal("normalizeLanguage() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "supported: en, ru, fr, de, zh") {
		t.Fatalf("normalizeLanguage() error = %q, want supported list", err)
	}
}

func TestLocalizerTReturnsMessageIDWhenDefaultBundleFails(t *testing.T) {
	oldLocaleData := localeData["en"]
	resetI18nBundleForTest()
	localeData["en"] = "[broken"
	t.Cleanup(func() {
		localeData["en"] = oldLocaleData
		resetI18nBundleForTest()
	})

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("appLocalizer.T() panic = %v", r)
		}
	}()

	if got := (appLocalizer{}).T("MissingMessage", nil); got != "MissingMessage" {
		t.Fatalf("appLocalizer.T() = %q, want MissingMessage", got)
	}
}

func TestDDCCommandErrorNoMonitor(t *testing.T) {
	err := ddcCommandError(
		defaultLocalizer(),
		[]string{"--bus", "6", "getvcp", "60"},
		"No monitor detected on bus /dev/i2c-6",
		errors.New("exit status 1"),
	)

	if err == nil {
		t.Fatal("ddcCommandError() error = nil, want error")
	}
	if !errors.Is(err, errNoMonitor) {
		t.Fatalf("ddcCommandError() error = %v, want errNoMonitor", err)
	}
	message := err.Error()
	for _, want := range []string{"monitor not found on bus 6", "dellkvm detect", "config.toml", "bus = 0"} {
		if !strings.Contains(message, want) {
			t.Fatalf("ddcCommandError() error = %q, want %q", message, want)
		}
	}
}

func TestDDCCommandErrorRetries(t *testing.T) {
	raw := `Error detecting VCP version using VCP feature xDF: Error_Info[DDCRC_RETRIES]
VCP code 0x60 (Input Source): Maximum retries exceeded`

	err := ddcCommandError(
		defaultLocalizer(),
		[]string{"--bus", "10", "getvcp", "60"},
		raw,
		errors.New("exit status 1"),
	)

	if err == nil {
		t.Fatal("ddcCommandError() error = nil, want error")
	}
	if !errors.Is(err, errDDCRetry) {
		t.Fatalf("ddcCommandError() error = %v, want errDDCRetry", err)
	}
	message := err.Error()
	for _, want := range []string{"DDC/VCP", "does not look like a sudo problem", "cable, dock, or port"} {
		if !strings.Contains(message, want) {
			t.Fatalf("ddcCommandError() error = %q, want %q", message, want)
		}
	}
}

func TestDDCCommandErrorRussian(t *testing.T) {
	l, err := newLocalizer("ru")
	if err != nil {
		t.Fatalf("newLocalizer() error = %v", err)
	}
	err = ddcCommandError(
		l,
		[]string{"--bus", "6", "getvcp", "60"},
		"No monitor detected on bus /dev/i2c-6",
		errors.New("exit status 1"),
	)
	if err == nil {
		t.Fatal("ddcCommandError() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "монитор не найден на bus 6") {
		t.Fatalf("ddcCommandError() error = %q, want Russian text", err)
	}
}

func TestDDCCommandErrorPreservesRawOutput(t *testing.T) {
	raw := "RAW DDC OUTPUT: do not translate"
	err := ddcCommandError(
		defaultLocalizer(),
		[]string{"--bus", "6", "capabilities"},
		raw,
		errors.New("exit status 1"),
	)
	if err == nil {
		t.Fatal("ddcCommandError() error = nil, want error")
	}
	if !strings.Contains(err.Error(), raw) {
		t.Fatalf("ddcCommandError() error = %q, want raw output %q", err, raw)
	}
}

func TestParseDDCBuses(t *testing.T) {
	raw := `Display 1
   I2C bus:             /dev/i2c-9
Display 2
   I2C bus:             /dev/i2c-6
Display 3
   I2C bus:             /dev/i2c-9
`

	got := parseDDCBuses(raw)
	want := []int{6, 9}
	if len(got) != len(want) {
		t.Fatalf("parseDDCBuses() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("parseDDCBuses() = %v, want %v", got, want)
		}
	}
}

func TestParseDDCBusesSkipsInvalidLaptopDisplay(t *testing.T) {
	raw := `Invalid display
   I2C bus:  /dev/i2c-10
   DRM_connector:           card1-eDP-1
   This is a laptop display.  Laptop displays do not support DDC/CI.

Display 1
   I2C bus:  /dev/i2c-13
   DRM_connector:           card1-DP-3
   EDID synopsis:
      Mfg id:               DEL - Dell Inc.
      Model:                DELL U2725QE
   VCP version:         2.1
`

	got := parseDDCBuses(raw)
	want := []int{13}
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("parseDDCBuses() = %v, want %v", got, want)
	}
}

func TestTUIStatusWrapsToTerminalWidth(t *testing.T) {
	model := newTUIModel(Config{
		Bus:      0,
		Language: "en",
		Inputs:   defaultInputs(),
	})
	model.width = 42
	model.height = 16
	model.status = "Error: ddc retries exceeded: monitor was found, but the DDC/VCP request did not get a stable response. This does not look like a sudo problem."
	model.resizeList()

	for _, line := range strings.Split(model.renderStatus(), "\n") {
		if width := lipgloss.Width(line); width > model.width {
			t.Fatalf("status line width = %d, want <= %d: %q", width, model.width, line)
		}
	}
}

func TestTUIDefaultLabelsEnglish(t *testing.T) {
	model := newTUIModel(Config{
		Bus:    0,
		Inputs: defaultInputs(),
	})

	if model.list.Title != "Inputs" {
		t.Fatalf("list title = %q, want Inputs", model.list.Title)
	}
	if got := model.renderCurrent(); !strings.Contains(got, "Current input: unknown") {
		t.Fatalf("renderCurrent() = %q, want English current label", got)
	}
	if got := model.renderKeys(); !strings.Contains(got, "Enter: switch") {
		t.Fatalf("renderKeys() = %q, want English keys", got)
	}
}

func TestSwitchVerificationMismatchReturnsError(t *testing.T) {
	withFakeDDC(t, "switch-mismatch")
	withoutSwitchDelay(t)

	cfg := Config{
		Bus:      6,
		Language: "en",
		Inputs: []Input{
			{ID: "dp", Name: "DisplayPort", Code: "0x0f"},
			{ID: "hdmi", Name: "HDMI", Code: "0x11"},
		},
	}

	_, err := switchInput(cfg, "dp")
	if err == nil {
		t.Fatal("switchInput() error = nil, want verification mismatch")
	}
	for _, want := range []string{"verification reported 0x11", "0x0f"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("switchInput() error = %q, want %q", err, want)
		}
	}
}

func TestExplicitBusDoesNotAutoDetectForSwitch(t *testing.T) {
	withFakeDDC(t, "explicit-bus-no-monitor")
	withoutSwitchDelay(t)

	cfg := Config{
		Bus:      6,
		Language: "en",
		Inputs: []Input{
			{ID: "dp", Name: "DisplayPort", Code: "0x0f"},
		},
	}

	_, err := switchInput(cfg, "dp")
	if err == nil {
		t.Fatal("switchInput() error = nil, want no monitor on explicit bus")
	}
	if !errors.Is(err, errNoMonitor) {
		t.Fatalf("switchInput() error = %v, want errNoMonitor", err)
	}
	if !strings.Contains(err.Error(), "bus 6") {
		t.Fatalf("switchInput() error = %q, want bus 6", err)
	}
}

func TestSwitchAutoDetectedPrefixOnNoVerify(t *testing.T) {
	withFakeDDC(t, "auto-detect-switch-no-verify")
	withoutSwitchDelay(t)

	cfg := Config{
		Bus:      0,
		Language: "en",
		Inputs: []Input{
			{ID: "dp", Name: "DisplayPort", Code: "0x0f"},
		},
	}

	result, err := switchInputDetailed(cfg, "dp")
	if err != nil {
		t.Fatalf("switchInputDetailed() error = %v, want nil", err)
	}
	if result.Code != "0x0f" || result.Bus != 13 || result.Verified {
		t.Fatalf("unverified write state = %+v", result)
	}
	for _, want := range []string{"Auto-detected bus 13", "Switch command sent"} {
		if !strings.Contains(result.Message, want) {
			t.Fatalf("switchInputDetailed() message = %q, want %q", result.Message, want)
		}
	}
}

func TestCurrentWorksWithoutConfig(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Chdir(t.TempDir())
	withFakeDDC(t, "single-monitor-current")
	out, err := captureStdout(t, printCurrent)
	if err != nil || !strings.Contains(out, "DisplayPort") {
		t.Fatalf("current without config = %q, %v", out, err)
	}
	if _, err := loadConfig(); !errors.Is(err, errConfigNotFound) {
		t.Fatalf("current created config: %v", err)
	}
}

func TestExplicitBusDoesNotAutoDetectForCurrent(t *testing.T) {
	withFakeDDC(t, "explicit-bus-no-monitor")

	cfg := Config{
		Bus:      6,
		Language: "en",
		Inputs: []Input{
			{ID: "dp", Name: "DisplayPort", Code: "0x0f"},
		},
	}

	_, err := getCurrentState(cfg, defaultLocalizer())
	if err == nil {
		t.Fatal("getCurrentState() error = nil, want no monitor on explicit bus")
	}
	if !errors.Is(err, errNoMonitor) {
		t.Fatalf("getCurrentState() error = %v, want errNoMonitor", err)
	}
	if !strings.Contains(err.Error(), "bus 6") {
		t.Fatalf("getCurrentState() error = %q, want bus 6", err)
	}
}

func TestRunDDCTimeout(t *testing.T) {
	withFakeDDC(t, "timeout")

	oldTimeout := ddcTimeout
	ddcTimeout = 50 * time.Millisecond
	t.Cleanup(func() { ddcTimeout = oldTimeout })

	_, err := runDDC(defaultLocalizer(), "detect", "--brief")
	if err == nil {
		t.Fatal("runDDC() error = nil, want timeout")
	}
	if !errors.Is(err, errDDCTimeout) {
		t.Fatalf("runDDC() error = %v, want errDDCTimeout", err)
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("runDDC() error = %q, want timeout message", err)
	}
}

func withFakeDDC(t *testing.T, scenario string) {
	t.Helper()

	oldExec := execDDCCommand
	execDDCCommand = func(ctx context.Context, command string, args ...string) *exec.Cmd {
		helperArgs := []string{"-test.run=TestHelperProcess", "--", command}
		helperArgs = append(helperArgs, args...)
		cmd := exec.CommandContext(ctx, os.Args[0], helperArgs...)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1", "DDC_SCENARIO="+scenario)
		return cmd
	}
	t.Cleanup(func() { execDDCCommand = oldExec })
}

func withoutSwitchDelay(t *testing.T) {
	t.Helper()

	oldDelay := verifyDelay
	verifyDelay = 0
	t.Cleanup(func() { verifyDelay = oldDelay })
}

func resetI18nBundleForTest() {
	bundleOnce = sync.Once{}
	bundle = nil
	bundleErr = nil
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	args := helperArgs()
	if len(args) > 0 && args[0] == "ddcutil" {
		args = args[1:]
	}

	switch os.Getenv("DDC_SCENARIO") {
	case "switch-mismatch":
		handleSwitchMismatch(args)
	case "explicit-bus-no-monitor":
		handleExplicitBusNoMonitor(args)
	case "auto-detect-switch-no-verify":
		handleAutoDetectSwitchNoVerify(args)
	case "timeout":
		time.Sleep(5 * time.Second)
		os.Exit(0)
	case "multiple-monitors":
		if containsArg(args, "detect") {
			_, _ = os.Stdout.WriteString("Display 1\n I2C bus: /dev/i2c-6\nDisplay 2\n I2C bus: /dev/i2c-13\n")
			os.Exit(0)
		}
		os.Exit(99)
	case "single-monitor-current":
		if containsArg(args, "detect") {
			_, _ = os.Stdout.WriteString("Display 1\n I2C bus: /dev/i2c-13\n")
		} else if containsArg(args, "getvcp") {
			_, _ = os.Stdout.WriteString("sl=0x0f\n")
		} else {
			os.Exit(99)
		}
		os.Exit(0)
	case "delayed-verification":
		path := os.Getenv("DDC_TEST_EVENTS")
		data, _ := os.ReadFile(path)
		event := "read\n"
		if containsArg(args, "setvcp") {
			event = "write\n"
		}
		if err := os.WriteFile(path, append(data, []byte(event)...), 0o600); err != nil {
			os.Exit(99)
		}
		if event == "write\n" {
			os.Exit(0)
		}
		code := "0x11"
		if strings.Count(string(data), "read\n") >= 1 {
			code = "0x0f"
		}
		_, _ = os.Stdout.WriteString("sl=" + code)
		os.Exit(0)
	default:
		os.Exit(1)
	}
}

func helperArgs() []string {
	for i, arg := range os.Args {
		if arg == "--" {
			return os.Args[i+1:]
		}
	}
	return nil
}

func handleSwitchMismatch(args []string) {
	if containsArg(args, "setvcp") {
		os.Exit(0)
	}
	if containsArg(args, "getvcp") {
		_, _ = os.Stdout.WriteString("VCP code 0x60 (Input Source): sl=0x11\n")
		os.Exit(0)
	}
	os.Exit(1)
}

func handleExplicitBusNoMonitor(args []string) {
	if containsArg(args, "detect") {
		_, _ = os.Stdout.WriteString("Display 1\n   I2C bus:  /dev/i2c-13\n")
		os.Exit(0)
	}
	if containsArg(args, "setvcp") && containsArg(args, "6") {
		_, _ = os.Stdout.WriteString("No monitor detected on bus /dev/i2c-6\n")
		os.Exit(1)
	}
	if containsArg(args, "getvcp") && containsArg(args, "6") {
		_, _ = os.Stdout.WriteString("No monitor detected on bus /dev/i2c-6\n")
		os.Exit(1)
	}
	if containsArg(args, "setvcp") && containsArg(args, "13") {
		os.Exit(0)
	}
	if containsArg(args, "getvcp") && containsArg(args, "13") {
		_, _ = os.Stdout.WriteString("VCP code 0x60 (Input Source): sl=0x0f\n")
		os.Exit(0)
	}
	os.Exit(1)
}

func handleAutoDetectSwitchNoVerify(args []string) {
	if containsArg(args, "detect") {
		_, _ = os.Stdout.WriteString("Display 1\n   I2C bus:  /dev/i2c-13\n")
		os.Exit(0)
	}
	if containsArg(args, "setvcp") && containsArg(args, "13") {
		os.Exit(0)
	}
	if containsArg(args, "getvcp") && containsArg(args, "13") {
		_, _ = os.Stdout.WriteString("VCP code 0x60 (Input Source): Maximum retries exceeded\n")
		os.Exit(1)
	}
	os.Exit(1)
}

func containsArg(args []string, want string) bool {
	for _, arg := range args {
		if arg == want {
			return true
		}
	}
	return false
}

func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()

	oldStdout := os.Stdout
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdout: %v", err)
	}
	os.Stdout = write
	defer func() {
		os.Stdout = oldStdout
	}()

	runErr := fn()

	os.Stdout = oldStdout
	if err := write.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}
	out, err := io.ReadAll(read)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	if err := read.Close(); err != nil {
		t.Fatalf("close stdout reader: %v", err)
	}
	return string(out), runErr
}

func assertFileContent(t *testing.T, path string, want string) {
	t.Helper()

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(got) != want {
		t.Fatalf("%s = %q, want %q", path, got, want)
	}
}
