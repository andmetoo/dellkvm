package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func openFile(path string) error {
	cmd := exec.Command("xdg-open", path)
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() { _ = cmd.Wait() }()
	return nil
}

func userConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "dellkvm", "config.toml"), nil
}

func runPlatformWorker(args []string) (bool, error) { return false, nil }
