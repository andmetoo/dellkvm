//go:build linux

package main

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

var execDDCCommand = exec.CommandContext

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
