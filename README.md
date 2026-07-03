# dellkvm

`dellkvm` is a small Linux-only terminal app for switching a monitor input through `ddcutil`. It never talks to I2C directly and never runs `sudo` internally.

## Requirements

- Linux
- `ddcutil` and I2C tools, for example: `sudo dnf install ddcutil i2c-tools`
- Access to `/dev/i2c-*`, usually by adding your user to the `i2c` group

## Build

```sh
make setup
make build
```

Useful development commands:

```sh
make test
make lint
make build-all
```

## Configuration

`dellkvm` reads `./config.toml` first, then `~/.config/dellkvm/config.toml`. Copy `config.toml.default` as a starting point.

```toml
bus = 0
language = "en"

[[inputs]]
id = "tb"
name = "Thunderbolt / USB-C"
code = "0x19"

[[inputs]]
id = "dp"
name = "DisplayPort"
code = "0x0f"
```

Use `bus = 0` for auto-detection. Set a positive bus number to force a specific display. Supported languages are `en`, `ru`, `fr`, `de`, and `zh`. Only text produced by `dellkvm` is translated; raw `ddcutil` output is shown unchanged.

To inspect input-source values supported by a display, query VCP feature `60` for its bus:

```sh
ddcutil -b 13 capabilities | grep -A 4 "Feature: 60"
```

Example output:

```text
Feature: 60 (Input Source)
   Values:
      19: Unrecognized value
      0f: DisplayPort-1
      11: HDMI-1
```

Use the detected values as `code` entries in `config.toml`, adding the `0x` prefix.

For example, on a Dell U2725QE, `19: Unrecognized value` may be Thunderbolt / USB-C, so the config value is `0x19`.

## Commands

- `dellkvm` opens the TUI.
- `dellkvm detect` prints `ddcutil detect --brief`.
- `dellkvm current` shows the current input.
- `dellkvm switch <id>` switches to a configured input.
- `dellkvm learn` creates or updates the config by reading input codes.
