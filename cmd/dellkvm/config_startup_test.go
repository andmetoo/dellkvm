//go:build linux

package main

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestStartupConfigCreatesPersistentDefaults(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Chdir(t.TempDir())

	cfg, path, err := startupConfig()
	if err != nil {
		t.Fatal(err)
	}
	wantPath := filepath.Join(home, ".config", "dellkvm", "config.toml")
	if path != wantPath {
		t.Fatalf("path = %q, want %q", path, wantPath)
	}
	if cfg.Bus != 0 || cfg.Language != "en" || !reflect.DeepEqual(cfg.Inputs, defaultInputs()) {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	loaded, err := ensureConfig()
	if err != nil || !reflect.DeepEqual(loaded, cfg) {
		t.Fatalf("persisted config = %+v, %v", loaded, err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("permissions = %o", info.Mode().Perm())
	}
	if _, err := os.Stat("config.toml"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("unexpected local config: %v", err)
	}
}

func TestStartupConfigPreservesExistingFiles(t *testing.T) {
	valid := "bus = 7\nlanguage = \"ru\"\n[[inputs]]\nid = \"custom\"\nname = \"My input\"\ncode = \"0x12\"\n"
	for _, tt := range []struct {
		name, local, user string
		bad               bool
	}{
		{"local takes precedence", valid, "invalid [", false},
		{"user config", "", valid, false},
		{"invalid local", "invalid [", valid, true},
		{"invalid user", "", "invalid [", true},
		{"invalid values", "bus = -1\n", "", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			cwd := t.TempDir()
			t.Chdir(cwd)
			localPath := filepath.Join(cwd, "config.toml")
			userPath := filepath.Join(home, ".config", "dellkvm", "config.toml")
			for path, data := range map[string]string{localPath: tt.local, userPath: tt.user} {
				if data == "" {
					continue
				}
				if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			for range 2 {
				cfg, path, err := startupConfig()
				wantPath := userPath
				if tt.local != "" {
					wantPath = localPath
				}
				if path != wantPath || (err != nil) != tt.bad {
					t.Fatalf("path = %q, error = %v", path, err)
				}
				if !tt.bad && (cfg.Bus != 7 || cfg.Inputs[0].ID != "custom") {
					t.Fatalf("config replaced: %+v", cfg)
				}
				if tt.bad && !strings.Contains(err.Error(), wantPath) {
					t.Fatalf("error lacks config path: %v", err)
				}
			}
			for path, data := range map[string]string{localPath: tt.local, userPath: tt.user} {
				got, err := os.ReadFile(path)
				if data == "" {
					if !errors.Is(err, os.ErrNotExist) {
						t.Fatalf("unexpected config at %q: %v", path, err)
					}
				} else if err != nil || string(got) != data {
					t.Fatalf("config changed at %q: %q, %v", path, got, err)
				}
			}
		})
	}
}

func TestStartupConfigReportsCreationFailure(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Chdir(t.TempDir())
	if err := os.WriteFile(filepath.Join(home, ".config"), []byte("not a directory"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := startupConfig(); err == nil {
		t.Fatal("creation failure was ignored")
	}
}

func TestRunInit(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Chdir(t.TempDir())
	if err := run([]string{"init", "extra"}); exitCode(err) != 2 {
		t.Fatalf("invalid arguments: %v", err)
	}
	if _, err := loadConfig(); !errors.Is(err, errConfigNotFound) {
		t.Fatalf("invalid arguments created config: %v", err)
	}
	out, err := captureStdout(t, func() error { return run([]string{"init"}) })
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".config", "dellkvm", "config.toml")
	if strings.TrimSpace(out) != want {
		t.Fatalf("init output = %q, want %q", out, want)
	}
	if _, err := ensureConfig(); err != nil {
		t.Fatal(err)
	}
}

func TestDesktopConfigDoesNotCreateFiles(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Chdir(t.TempDir())
	if _, err := desktopConfig(); err != nil {
		t.Fatal(err)
	}
	if got := configLocationTitle(); got != "Config: built-in defaults (optional)" {
		t.Fatalf("tray config label = %q", got)
	}
	if _, err := loadConfig(); !errors.Is(err, errConfigNotFound) {
		t.Fatalf("read-only loader created config: %v", err)
	}
}
