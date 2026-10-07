# Repository Guidelines

## Project Structure & Module Organization

This repository contains `dellkvm`, a Linux and Windows Go tray/CLI/TUI for switching monitor inputs.

- `cmd/dellkvm/main.go` contains the CLI commands, TUI model, config handling, and `ddcutil` integration.
- `cmd/dellkvm/main_test.go` contains Go unit tests for parsing, config, and command behavior.
- `config.toml.default` is the copyable default configuration. Local `config.toml` is ignored and takes priority over `~/.config/dellkvm/config.toml`.
- `mise.toml`, `Makefile`, and `.goreleaser.yaml` define local tooling, developer commands, and release builds.

Do not add direct I2C access on Linux. Linux monitor operations go through `exec.Command` and `ddcutil`. Windows uses the system monitor configuration API in a child process so calls have a process timeout.

## Build, Test, and Development Commands

- `make` shows available targets.
- `make setup` trusts and installs tools from `mise.toml`.
- `make build` builds the local `dellkvm` binary through GoReleaser snapshot mode.
- `make build-all` builds all configured Linux and Windows release targets into `dist/`.
- `make test` runs `go test ./...`.
- `make lint` runs the current lint task, currently `go vet ./...`.
- `make release-check` validates the GoReleaser configuration.

If the Go build cache is not writable in a sandbox, run tests with `GOCACHE=/tmp/dellkvm-go-build-cache make test`.

## Coding Style & Naming Conventions

Use standard Go formatting with `gofmt`. Keep functions small and focused; existing helper names such as `loadConfig`, `runDDC`, `getCurrentCode`, and `switchInput` are the preferred style. Keep symbols unexported unless they are needed outside the package. Return user-facing errors instead of panicking in normal CLI scenarios.

Configuration is TOML. Preserve the `bus = 0` auto-detect behavior, and keep input codes normalized as lowercase `0xNN` strings.

## Testing Guidelines

Use Go's built-in `testing` package. Name tests `Test...` and prefer table-driven tests for parsers and command behavior. Add tests when changing config loading, VCP parsing, bus detection, or error classification. Hardware-dependent behavior should be isolated behind command execution helpers so tests remain deterministic.

## Commit & Pull Request Guidelines

The history currently has only an initial commit, so no strict project convention is established. Use concise imperative commit subjects, for example `Add TOML config loading`; Conventional Commit prefixes are acceptable when useful.

Pull requests should describe behavior changes, list commands run (`make test`, `make lint`), and note any DDC/CI hardware assumptions or manual monitor checks.

## Security & Configuration Tips

Never invoke `sudo` from the application. Keep local machine settings in ignored `config.toml` or `~/.config/dellkvm/config.toml`; do not commit personal bus numbers unless they belong in `config.toml.default`.
