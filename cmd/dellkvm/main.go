//go:build linux

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pelletier/go-toml/v2"
)

const (
	vcpInputSource = "60"
)

var (
	errConfigNotFound = errors.New("config not found")
	errDDCNotFound    = errors.New("ddcutil not found")
	errI2CAccess      = errors.New("i2c access denied")
	errNoMonitor      = errors.New("monitor not found")
	errNoI2CDevices   = errors.New("i2c devices not found")
	errDDCRetry       = errors.New("ddc retries exceeded")
	errDDCTimeout     = errors.New("ddcutil timeout")

	vcpCodePattern   = regexp.MustCompile(`sl=0x([0-9a-fA-F]+)`)
	i2cBusPattern    = regexp.MustCompile(`/dev/i2c-([0-9]+)`)
	inputCodePattern = regexp.MustCompile(`^0x[0-9a-f]{2}$`)
	ddcTimeout       = 10 * time.Second
	verifyDelay      = 500 * time.Millisecond
	execDDCCommand   = exec.CommandContext
	ddcContext       = context.Background()
)

type Config struct {
	Bus      int     `toml:"bus"`
	Language string  `toml:"language"`
	Inputs   []Input `toml:"inputs"`
}

type Input struct {
	ID   string `toml:"id" json:"id"`
	Name string `toml:"name" json:"name"`
	Code string `toml:"code" json:"code"`
}

type currentState struct {
	Code         string
	Raw          string
	Bus          int
	AutoDetected bool
}

type usageError string

func (e usageError) Error() string {
	return string(e)
}

type userError struct {
	err     error
	message string
}

func (e userError) Error() string {
	return e.message
}

func (e userError) Unwrap() error {
	return e.err
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCode(err))
	}
}

func run(args []string) error {
	l := defaultLocalizer()
	if len(args) == 0 {
		return runTray()
	}

	switch args[0] {
	case "init":
		if len(args) != 1 {
			return usageError("Usage: dellkvm init")
		}
		_, path, err := startupConfig()
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	case "status":
		if len(args) != 2 || args[1] != "--json" {
			return usageError("Usage: dellkvm status --json")
		}
		return printStatusJSON()
	case "tray", "tui":
		if len(args) != 1 {
			return usageError("Usage: dellkvm [tray|tui]")
		}
		if args[0] == "tui" {
			return runTUI()
		}
		return runTray()
	case "detect":
		if len(args) != 1 {
			return usageError(l.T("UsageDetect", nil))
		}
		out, err := detectDisplays(l)
		if err != nil {
			return err
		}
		printOutput(out)
		return nil
	case "current":
		if len(args) != 1 {
			return usageError(l.T("UsageCurrent", nil))
		}
		return printCurrent()
	case "switch":
		if len(args) != 2 {
			return usageError(l.T("UsageSwitch", nil))
		}
		return printSwitch(args[1])
	case "help", "--help", "-h":
		if len(args) != 1 {
			return usageError(l.T("UsageRoot", nil))
		}
		printHelp(l)
		return nil
	default:
		return usageError(l.T("UsageRoot", nil))
	}
}

// Machine-readable status for desktop integrations. Hardware errors remain
// in the payload so the UI can still present configured input buttons.
func printStatusJSON() error {
	result := struct {
		Inputs []Input `json:"inputs"`
		Code   string  `json:"code"`
		Bus    int     `json:"bus"`
		Error  string  `json:"error,omitempty"`
	}{}
	cfg, err := desktopConfig()
	if err == nil {
		result.Inputs = cfg.Inputs
		var state currentState
		state, err = getCurrentState(cfg, localizerForConfig(cfg))
		result.Code, result.Bus = state.Code, state.Bus
	}
	if err != nil {
		result.Error = err.Error()
	}
	return json.NewEncoder(os.Stdout).Encode(result)
}

func exitCode(err error) int {
	var usage usageError
	switch {
	case errors.As(err, &usage):
		return 2
	case errors.Is(err, errDDCNotFound):
		return 127
	default:
		return 1
	}
}

func loadConfig() (Config, error) {
	l := defaultLocalizer()
	paths, err := configPaths()
	if err != nil {
		return Config{}, err
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return Config{}, err
		}

		var cfg Config
		if err := toml.Unmarshal(data, &cfg); err != nil {
			return Config{}, fmt.Errorf("%s: %w", l.T("ConfigReadFailed", map[string]any{"Path": path}), err)
		}
		if cfg.Inputs == nil {
			cfg.Inputs = []Input{}
		}
		applyConfigDefaults(&cfg)
		return cfg, nil
	}

	return Config{}, fmt.Errorf("%w: %s", errConfigNotFound, strings.Join(paths, ", "))
}

func ensureConfig() (Config, error) {
	l := defaultLocalizer()
	cfg, err := loadConfig()
	if errors.Is(err, errConfigNotFound) {
		paths, pathsErr := configPaths()
		if pathsErr != nil {
			return Config{}, pathsErr
		}
		return Config{}, errors.New(l.T("ConfigNotFound", map[string]any{"Paths": strings.Join(paths, ", ")}))
	}
	if err != nil {
		return Config{}, err
	}
	if _, err := normalizeLanguage(cfg.Language); err != nil {
		return Config{}, err
	}
	l, err = newLocalizer(cfg.Language)
	if err != nil {
		return Config{}, err
	}
	if cfg.Bus < 0 {
		return Config{}, errors.New(l.T("InvalidBus", nil))
	}
	if len(cfg.Inputs) == 0 {
		return Config{}, errors.New(l.T("MissingInputs", nil))
	}
	if err := validateInputs(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func validateInputs(cfg *Config) error {
	ids := make(map[string]bool)
	codes := make(map[string]bool)
	for i := range cfg.Inputs {
		input := &cfg.Inputs[i]
		input.ID = strings.TrimSpace(input.ID)
		input.Name = strings.TrimSpace(input.Name)
		input.Code = normalizeCode(input.Code)
		if input.ID == "" || input.Name == "" {
			return fmt.Errorf("input %d: id and name must not be empty", i+1)
		}
		if ids[input.ID] {
			return fmt.Errorf("duplicate input id %q", input.ID)
		}
		if !inputCodePattern.MatchString(input.Code) {
			return fmt.Errorf("input %q: code must be a hexadecimal byte (0x00–0xff)", input.ID)
		}
		if codes[input.Code] {
			return fmt.Errorf("duplicate input code %q", input.Code)
		}
		ids[input.ID], codes[input.Code] = true, true
	}
	return nil
}

func applyConfigDefaults(cfg *Config) {
	if strings.TrimSpace(cfg.Language) == "" {
		cfg.Language = defaultLanguage
	}
}

func localizerForConfig(cfg Config) appLocalizer {
	l, err := newLocalizer(cfg.Language)
	if err != nil {
		return defaultLocalizer()
	}
	return l
}

func runDDC(l appLocalizer, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ddcContext, ddcTimeout)
	defer cancel()

	cmd := execDDCCommand(ctx, "ddcutil", args...)
	// Bound waiting for pipes even if a child process keeps them open.
	cmd.WaitDelay = time.Second
	out, err := cmd.CombinedOutput()
	raw := string(out)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return "", ddcTimeoutError(l, args, raw)
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return "", context.Canceled
	}
	if err == nil {
		return raw, nil
	}

	var execErr *exec.Error
	if errors.As(err, &execErr) && errors.Is(execErr.Err, exec.ErrNotFound) {
		return "", userError{err: errDDCNotFound, message: l.T("MissingDDC", nil)}
	}
	return "", ddcCommandError(l, args, raw, err)
}

func ddcTimeoutError(l appLocalizer, args []string, raw string) error {
	data := map[string]any{
		"Seconds": int(ddcTimeout.Seconds()),
		"Command": strings.Join(args, " "),
		"Raw":     strings.TrimSpace(raw),
	}
	if strings.TrimSpace(raw) == "" {
		return userError{err: errDDCTimeout, message: l.T("DDCTimeout", data)}
	}
	return userError{err: errDDCTimeout, message: l.T("DDCTimeoutWithOutput", data)}
}

func detectDisplays(l appLocalizer) (string, error) {
	return runDDC(l, "detect", "--brief")
}

func getCurrentState(cfg Config, l appLocalizer) (currentState, error) {
	if cfg.Bus > 0 {
		return getCurrentStateOnBus(cfg.Bus, l)
	}

	return autoDetectCurrentState(l)
}

func getCurrentStateOnBus(bus int, l appLocalizer) (currentState, error) {
	raw, err := runDDC(l, "--bus", strconv.Itoa(bus), "getvcp", vcpInputSource)
	if err != nil {
		return currentState{}, err
	}

	code, err := parseVCPCodeWithLocalizer(raw, l)
	if err != nil {
		return currentState{}, err
	}
	return currentState{Code: code, Raw: raw, Bus: bus}, nil
}

func switchInput(cfg Config, id string) (string, error) {
	l := localizerForConfig(cfg)
	input, ok := findInputByID(cfg, id)
	if !ok {
		return "", errors.New(l.T("InputNotFound", map[string]any{"ID": id}))
	}

	bus, autoDetected, err := resolveBusForSwitch(cfg, l)
	if err != nil {
		return "", err
	}

	_, err = runDDC(l, "--bus", strconv.Itoa(bus), "--noverify", "setvcp", vcpInputSource, input.Code)
	if err != nil {
		return "", err
	}

	codeState, err := verifyInput(bus, input.Code, l)
	if err != nil {
		message := l.T("SwitchSentNoVerify", map[string]any{
			"Name":  input.Name,
			"ID":    input.ID,
			"Code":  input.Code,
			"Bus":   bus,
			"VCP":   vcpInputSource,
			"Error": err,
		})
		return withAutoDetectedPrefix(l, bus, autoDetected, message), nil
	}

	if normalizeCode(codeState.Code) != normalizeCode(input.Code) {
		return "", errors.New(l.T("SwitchVerificationMismatch", map[string]any{
			"Name":      input.Name,
			"ID":        input.ID,
			"Requested": input.Code,
			"Bus":       bus,
			"Observed":  codeState.Code,
		}))
	}

	current, ok := findInputByCode(cfg, codeState.Code)
	if !ok {
		message := l.T("SwitchVerifiedCode", map[string]any{
			"Name":        input.Name,
			"ID":          input.ID,
			"Code":        input.Code,
			"Bus":         bus,
			"CurrentCode": codeState.Code,
		})
		return withAutoDetectedPrefix(l, bus, autoDetected, message), nil
	}

	message := l.T("SwitchVerifiedInput", map[string]any{
		"Name":        input.Name,
		"ID":          input.ID,
		"Code":        input.Code,
		"Bus":         bus,
		"CurrentName": current.Name,
		"CurrentCode": codeState.Code,
	})
	return withAutoDetectedPrefix(l, bus, autoDetected, message), nil
}

func verifyInput(bus int, requested string, l appLocalizer) (currentState, error) {
	var state currentState
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		time.Sleep(verifyDelay)
		state, err = getCurrentStateOnBus(bus, l)
		if err == nil && normalizeCode(state.Code) == normalizeCode(requested) {
			return state, nil
		}
		// Retry reads only: a write may have succeeded even if the host loses
		// its connection. Never send another switch after an ambiguous result.
		if err != nil && !errors.Is(err, errDDCRetry) {
			return state, err
		}
	}
	return state, err
}

func withAutoDetectedPrefix(l appLocalizer, bus int, autoDetected bool, message string) string {
	if autoDetected {
		return l.T("AutoDetectedPrefix", map[string]any{"Bus": bus, "Message": message})
	}
	return message
}

func parseVCPCode(raw string) (string, error) {
	return parseVCPCodeWithLocalizer(raw, defaultLocalizer())
}

func parseVCPCodeWithLocalizer(raw string, l appLocalizer) (string, error) {
	match := vcpCodePattern.FindStringSubmatch(raw)
	if len(match) != 2 {
		return "", errors.New(l.T("VCPCodeNotFound", map[string]any{"Raw": strings.TrimSpace(raw)}))
	}

	value, err := strconv.ParseUint(match[1], 16, 8)
	if err != nil {
		return "", fmt.Errorf("%s: %w", l.T("VCPCodeParseFailed", map[string]any{"Value": match[1]}), err)
	}
	return fmt.Sprintf("0x%02x", value), nil
}

func runTUI() error {
	cfg, _, err := startupConfig()
	if err != nil {
		return err
	}

	model := newTUIModel(cfg)
	_, err = tea.NewProgram(model).Run()
	return err
}

func printCurrent() error {
	cfg, err := ensureConfig()
	if err != nil {
		return err
	}
	l := localizerForConfig(cfg)

	state, err := getCurrentState(cfg, l)
	if err != nil {
		return err
	}

	input, ok := findInputByCode(cfg, state.Code)
	if !ok {
		fmt.Println(l.T("CurrentCode", map[string]any{"Code": state.Code, "Bus": state.Bus}))
		return nil
	}

	fmt.Println(l.T("CurrentInput", map[string]any{
		"Name": input.Name,
		"ID":   input.ID,
		"Code": state.Code,
		"Bus":  state.Bus,
	}))
	return nil
}

func printSwitch(id string) error {
	cfg, err := desktopConfig()
	if err != nil {
		return err
	}

	message, err := switchInput(cfg, id)
	if err != nil {
		return err
	}
	fmt.Println(message)
	return nil
}

func configPaths() ([]string, error) {
	userPath, err := userConfigPath()
	if err != nil {
		return nil, err
	}

	return []string{localConfigPath(), userPath}, nil
}

func localConfigPath() string {
	return "config.toml"
}

func userConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "dellkvm", "config.toml"), nil
}

func printOutput(out string) {
	fmt.Print(out)
	if out != "" && !strings.HasSuffix(out, "\n") {
		fmt.Println()
	}
}

func printHelp(l appLocalizer) {
	fmt.Println(l.T("Help", nil))
}

func defaultInputs() []Input {
	return []Input{
		{ID: "tb", Name: "Thunderbolt / USB-C", Code: "0x19"},
		{ID: "dp", Name: "DisplayPort", Code: "0x0f"},
		{ID: "hdmi", Name: "HDMI", Code: "0x11"},
	}
}

func findInputByID(cfg Config, id string) (Input, bool) {
	for _, input := range cfg.Inputs {
		if input.ID == id {
			return input, true
		}
	}
	return Input{}, false
}

func findInputByCode(cfg Config, code string) (Input, bool) {
	code = normalizeCode(code)
	for _, input := range cfg.Inputs {
		if normalizeCode(input.Code) == code {
			return input, true
		}
	}
	return Input{}, false
}

func normalizeCode(code string) string {
	code = strings.TrimSpace(strings.ToLower(code))
	if !strings.HasPrefix(code, "0x") {
		return code
	}

	value, err := strconv.ParseUint(strings.TrimPrefix(code, "0x"), 16, 16)
	if err != nil {
		return code
	}
	return fmt.Sprintf("0x%02x", value)
}

func isI2CAccessError(raw string) bool {
	raw = strings.ToLower(raw)
	mentionsI2C := strings.Contains(raw, "/dev/i2c") || strings.Contains(raw, "i2c")
	permissionDenied := strings.Contains(raw, "permission denied") ||
		strings.Contains(raw, "operation not permitted") ||
		strings.Contains(raw, "access denied")
	return mentionsI2C && permissionDenied
}

func ddcCommandError(l appLocalizer, args []string, raw string, err error) error {
	if isI2CAccessError(raw) {
		return userError{err: errI2CAccess, message: l.T("I2CAccess", nil)}
	}
	if isNoI2CDevicesError(raw) {
		return userError{err: errNoI2CDevices, message: l.T("NoI2CDevices", nil)}
	}
	if isNoMonitorError(raw) {
		return userError{err: errNoMonitor, message: noMonitorMessage(l, args, raw)}
	}
	if isDDCRetryError(raw) {
		return userError{err: errDDCRetry, message: l.T("DDCRetry", nil)}
	}
	if strings.TrimSpace(raw) == "" {
		return errors.New(l.T("DDCCommandFailedNoOutput", map[string]any{
			"Command": strings.Join(args, " "),
			"Error":   err,
		}))
	}
	return errors.New(l.T("DDCCommandFailedWithOutput", map[string]any{
		"Command": strings.Join(args, " "),
		"Error":   err,
		"Raw":     strings.TrimSpace(raw),
	}))
}

func isNoMonitorError(raw string) bool {
	raw = strings.ToLower(raw)
	return strings.Contains(raw, "no monitor detected on bus")
}

func isNoI2CDevicesError(raw string) bool {
	raw = strings.ToLower(raw)
	return strings.Contains(raw, "no /dev/i2c devices exist") ||
		strings.Contains(raw, "requires module i2c-dev")
}

func isDDCRetryError(raw string) bool {
	raw = strings.ToLower(raw)
	return strings.Contains(raw, "maximum retries exceeded") ||
		strings.Contains(raw, "ddcrc_retries")
}

func resolveBusForSwitch(cfg Config, l appLocalizer) (int, bool, error) {
	if cfg.Bus > 0 {
		return cfg.Bus, false, nil
	}
	return autoDetectBusForSwitch(l)
}

func autoDetectBusForSwitch(l appLocalizer) (int, bool, error) {
	buses, err := detectBuses(l)
	if err != nil {
		return 0, false, fmt.Errorf("%s: %w", l.T("AutoDetectFailed", nil), err)
	}
	if len(buses) == 0 {
		return 0, false, errors.New(l.T("AutoDetectNoDisplays", nil))
	}
	if len(buses) > 1 {
		return 0, false, fmt.Errorf("multiple monitors detected (%s): select a monitor in the tray or set bus in config.toml", formatBuses(buses))
	}

	for _, bus := range buses {
		_, err := getCurrentStateOnBus(bus, l)
		if err == nil || errors.Is(err, errDDCRetry) {
			return bus, true, nil
		}
		if !errors.Is(err, errNoMonitor) {
			return 0, false, err
		}
	}

	return 0, false, errors.New(l.T("AutoDetectNoResponsive", map[string]any{
		"Buses": formatBuses(buses),
		"VCP":   vcpInputSource,
	}))
}

func autoDetectCurrentState(l appLocalizer) (currentState, error) {
	buses, err := detectBuses(l)
	if err != nil {
		return currentState{}, fmt.Errorf("%s: %w", l.T("AutoDetectFailed", nil), err)
	}
	if len(buses) == 0 {
		return currentState{}, errors.New(l.T("AutoDetectNoDisplays", nil))
	}
	if len(buses) > 1 {
		return currentState{}, fmt.Errorf("multiple monitors detected (%s): select a monitor in the tray or set bus in config.toml", formatBuses(buses))
	}

	for _, bus := range buses {
		state, err := getCurrentStateOnBus(bus, l)
		if err == nil {
			state.AutoDetected = true
			return state, nil
		}
		if !errors.Is(err, errNoMonitor) {
			return currentState{}, err
		}
	}

	return currentState{}, errors.New(l.T("AutoDetectNoResponsive", map[string]any{
		"Buses": formatBuses(buses),
		"VCP":   vcpInputSource,
	}))
}

func detectBuses(l appLocalizer) ([]int, error) {
	raw, err := detectDisplays(l)
	if err != nil {
		return nil, err
	}
	return parseDDCBuses(raw), nil
}

func parseDDCBuses(raw string) []int {
	seen := map[int]bool{}
	buses := []int{}
	inValidDisplay := false

	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Invalid display") {
			inValidDisplay = false
			continue
		}
		if strings.HasPrefix(line, "Display ") {
			inValidDisplay = true
			continue
		}
		if !inValidDisplay {
			continue
		}

		bus := busFromOutput(line)
		if bus == "" {
			continue
		}

		busNumber, err := strconv.Atoi(bus)
		if err != nil || seen[busNumber] {
			continue
		}
		seen[busNumber] = true
		buses = append(buses, busNumber)
	}

	sort.Ints(buses)
	return buses
}

func formatBuses(buses []int) string {
	parts := make([]string, 0, len(buses))
	for _, bus := range buses {
		parts = append(parts, strconv.Itoa(bus))
	}
	return strings.Join(parts, ", ")
}

func noMonitorMessage(l appLocalizer, args []string, raw string) string {
	bus := busFromArgs(args)
	if bus == "" {
		bus = busFromOutput(raw)
	}
	if bus == "" {
		return l.T("NoMonitorUnknownBus", nil)
	}
	return l.T("NoMonitorBus", map[string]any{"Bus": bus})
}

func busFromArgs(args []string) string {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "--bus" {
			return args[i+1]
		}
	}
	return ""
}

func busFromOutput(raw string) string {
	match := i2cBusPattern.FindStringSubmatch(raw)
	if len(match) != 2 {
		return ""
	}
	return match[1]
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
	status      string
	busy        bool
}

type currentMsg struct {
	state currentState
	err   error
}

type switchMsg struct {
	message string
	err     error
}

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
	return refreshCurrentCmd(m.cfg, m.localizer)
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
	case switchMsg:
		return m.updateSwitch(msg)
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
	case "enter":
		if m.busy {
			return m, nil
		}

		item, ok := m.list.SelectedItem().(inputItem)
		if !ok {
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

func (m tuiModel) updateCurrent(msg currentMsg) (tea.Model, tea.Cmd) {
	m.busy = false
	if msg.err != nil {
		m.currentCode, m.currentName, m.currentBus = "", "", 0
		m.status = m.localizer.T("ErrorPrefix", map[string]any{"Error": msg.err.Error()})
		m.resizeList()
		return m, m.list.SetItems(itemsFromConfig(m.cfg, ""))
	}

	m.currentCode = msg.state.Code
	m.currentBus = msg.state.Bus
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

	m.status = msg.message
	// Switching can disconnect this host. Preserve the verification result
	// instead of replacing it with an immediate failing refresh.
	m.currentCode, m.currentName, m.currentBus = "", "", 0
	m.resizeList()
	return m, m.list.SetItems(itemsFromConfig(m.cfg, ""))
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
	if m.currentName != "" {
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
		Render(m.localizer.T("Keys", nil))
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

func switchInputCmd(cfg Config, id string) tea.Cmd {
	return func() tea.Msg {
		message, err := switchInput(cfg, id)
		return switchMsg{message: message, err: err}
	}
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
