# dellkvm

`dellkvm` switches monitor inputs from the Linux system tray, a terminal UI, or the command line. An optional Omarchy shell widget is included. Monitor operations use `ddcutil`; the application never accesses I2C directly or runs `sudo` internally.

Windows monitor control is not implemented yet. It requires a separate backend because `ddcutil` is Linux-only; the repository currently requires all monitor operations to use `ddcutil`.

## Requirements

- Linux
- `ddcutil` and I2C tools, for example: `sudo dnf install ddcutil i2c-tools`
- Access to `/dev/i2c-*`, usually by adding your user to the `i2c` group
- For the tray: a desktop session with D-Bus and a StatusNotifierItem-compatible tray (including Omarchy shell). No GTK or WebKit build dependencies are required.

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

`dellkvm` reads `./config.toml` first, then `~/.config/dellkvm/config.toml`. The tray and TUI create a starter at the user path on first launch if neither file exists, before contacting a monitor. Existing files are never overwritten. The tray displays the active config path in its menu.

To create the config without launching the UI or contacting a monitor, run `dellkvm init`. It prints the active file's absolute path. You can also copy `config.toml.default` as a starting point.

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

Auto-detection requires exactly one monitor. With multiple monitors, select one in the tray or set `bus` explicitly. A tray selection lasts until configuration reload or application exit; set `bus` in the file to persist it. Bus numbers may change after reconnecting hardware.

The starter config uses the inputs from `config.toml.default`. JSON status and the `switch` command use those defaults in memory when no configuration exists; they do not create a file. Existing invalid files are reported rather than ignored. Input IDs and names must be nonempty; IDs and codes must be unique. Codes are normalized to lowercase hexadecimal bytes.

The tray's **Edit configuration…** action opens the active config, or creates a starter at `~/.config/dellkvm/config.toml` if none exists. After editing, choose **Reload configuration**. Default input codes are examples; use the codes supported by your monitor. New tray and widget labels are currently English; existing command/TUI messages retain their translations.

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

- `dellkvm` or `dellkvm tray` opens the system tray.
- `dellkvm tui` opens the TUI.
- `dellkvm init` creates a missing starter config and prints its path, or validates the existing config without overwriting it.
- `dellkvm detect` prints `ddcutil detect --brief`.
- `dellkvm current` shows the current input.
- `dellkvm switch <id>` switches to a configured input.
- `dellkvm help` shows CLI help.
- `dellkvm status --json` prints configured inputs, current code, bus, and any configuration/hardware error for desktop integrations. Hardware errors are in the JSON `error` field; integrations must inspect it even when the process exits successfully.

## Reliability and status

The tray runs monitor requests serially and discards clicks while busy. Each `ddcutil` process has a 10-second timeout. After a switch, up to three reads allow the monitor time to change inputs; the write is sent only once. A successful write followed by a disconnected monitor is reported as **sent, unverified**, not as a confirmed switch. Use **Refresh / reconnect** after returning to this computer or reconnecting a cable.

The tray clears input checkmarks after switching until a new read confirms the active input. **Open status log…** shows the complete latest result if the desktop truncates the menu text. The log is stored in the user cache directory (`~/.cache/dellkvm/status.log` by default). Quit cancels the active monitor command.

Avoid controlling the same monitor simultaneously from separate tray, CLI, or plugin instances. Requests are serialized within each application instance, not across processes. Monitor/driver behavior still requires testing on the target hardware.

## Linux launcher and autostart

Put the built `dellkvm` binary on `PATH`. The included `integrations/linux/dellkvm.desktop` can be copied to `~/.local/share/applications/` for a launcher or `~/.config/autostart/` for desktop autostart. On Hyprland, an alternative is an `exec-once` entry launching `dellkvm tray`; use only one startup mechanism.

## Omarchy shell widget

The optional plugin targets the Quickshell-based Omarchy shell. It opens a popup with input buttons, the current input, and manual refresh. It invokes `dellkvm` with argument arrays and does not poll or switch automatically.

1. Install the binary on `PATH` and configure the monitor and input codes.
2. Copy `integrations/omarchy/` to `~/.config/omarchy/plugins/dellkvm.inputs/`.
3. Run `omarchy-shell shell rescanPlugins`, then `omarchy plugin enable dellkvm.inputs`.
4. If needed, position it with `omarchy bar move dellkvm.inputs --section right`.

The widget's `executable` setting accepts an absolute path to the binary. It uses the CLI's normal configuration lookup; use the user configuration path rather than relying on the shell's working directory. With multiple monitors, set `bus` in the config. No desktop settings are changed by building this project.
