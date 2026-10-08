# dellkvm

`dellkvm` is a small way to switch a monitor between connected computers without installing a large Dell management application. It runs from the system tray, a terminal UI, or the command line on Linux and Windows; an optional Omarchy shell widget is included. Linux monitor operations use `ddcutil`; Windows uses the system monitor API in a child process with a timeout. The application never runs `sudo`.

Despite its name, `dellkvm` is not limited to Dell monitors. It sends the standard DDC/CI input-source command (VCP `0x60`), so other monitors that implement this command may work after their input codes are configured. A monitor's built-in KVM may follow the selected video input if its firmware supports that behavior. The app does not control a separate USB or HDMI KVM switch, and monitor/driver compatibility still needs a real hardware check.

## Requirements

- Linux: `ddcutil`, I2C tools, and access to `/dev/i2c-*` (usually through the `i2c` group). The tray requires D-Bus and a StatusNotifierItem-compatible desktop such as Omarchy shell.
- Windows: a monitor with DDC/CI enabled in its on-screen menu. The Windows API is included in the operating system; `ddcutil` is not needed.

## Install

On Windows, download the `windows_amd64_setup.exe` or `windows_arm64_setup.exe` file from the [latest release](https://github.com/andmetoo/dellkvm/releases/latest). The per-user installer places both the tray and CLI executables in `%LOCALAPPDATA%\Programs\dellkvm`, creates a Start Menu shortcut, and enables tray autostart at sign-in by default. You can uncheck autostart during setup. The installer can also launch the tray immediately. Uninstall it from Windows **Installed apps**; your `%APPDATA%\dellkvm\config.toml` is kept. Installers are currently unsigned.

On Linux, download and run the release installer script:

```sh
curl -fsSLO https://github.com/andmetoo/dellkvm/releases/latest/download/install.sh
sh install.sh
```

It verifies the release archive against `checksums.txt`, installs to `~/.local/bin`, adds an application launcher, and enables tray autostart for the next login in `~/.config/autostart`. It does not use `sudo` or change your monitor config. Run `~/.local/bin/dellkvm tray` to start it now. Use `sh install.sh --no-autostart` to disable autostart, `--version vX.Y.Z` to install a specific release, or `--uninstall` to remove the managed files while keeping your config. Install `ddcutil` and grant access to `/dev/i2c-*` through your distribution's normal setup. Portable Windows ZIPs and Linux tarballs remain available on the release page.

## Build

```sh
make setup
make build            # local Linux binary in dist/
make build-all        # Linux and Windows binaries in dist/
```

Useful development commands:

```sh
make test
make lint
make build-all
```

## Configuration

`dellkvm` reads `./config.toml` first. The user fallback is `~/.config/dellkvm/config.toml` on Linux and `%APPDATA%\dellkvm\config.toml` on Windows. The tray and TUI create a starter at the user path on first launch if neither file exists, before contacting a monitor. Existing files are never overwritten. The tray displays the active config path in its menu.

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

Use `bus = 0` for auto-detection. Set a positive bus number to force a specific display. On Linux this is a `ddcutil` I²C bus number; on Windows it is the monitor index printed by `dellkvm detect`. Windows indices can change when displays or docks are reconnected. Supported languages are `en`, `ru`, `fr`, `de`, and `zh`. Raw monitor output is shown unchanged.

Auto-detection requires exactly one monitor. With multiple monitors, select one in the tray or TUI (`m` cycles monitors), or set `bus` explicitly. A tray or TUI selection lasts until configuration reload or application exit; set `bus` in the file to persist it.

The starter config uses the inputs from `config.toml.default`. JSON status and the `switch` command use those defaults in memory when no configuration exists; they do not create a file. Existing invalid files are reported rather than ignored. Input IDs and names must be nonempty; IDs and codes must be unique. Codes are normalized to lowercase hexadecimal bytes.

The tray's **Edit configuration…** action opens the active config, or creates a starter at `~/.config/dellkvm/config.toml` if none exists. After editing, choose **Reload configuration**. Default input codes are examples; use the codes supported by your monitor. New tray and widget labels are currently English; existing command/TUI messages retain their translations.

On Linux, inspect input-source values supported by a display by querying VCP feature `60` for its bus:

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
- `dellkvm tui` opens the TUI. Press `m` to choose a monitor, Enter to switch, `r` to refresh, or `q` to quit.
- `dellkvm init` creates a missing starter config and prints its path, or validates the existing config without overwriting it.
- `dellkvm detect` lists monitors (`ddcutil detect --brief` on Linux).
- `dellkvm current` shows the current input.
- `dellkvm switch <id>` switches to a configured input.
- `dellkvm help` shows CLI help.
- `dellkvm status --json` prints configured inputs, current code, bus, and any configuration/hardware error for desktop integrations. Hardware errors are in the JSON `error` field; integrations must inspect it even when the process exits successfully.

## Reliability and status

The tray runs monitor requests serially and discards clicks while busy. Each monitor request has a 10-second process timeout. After a switch, up to three reads allow the monitor time to change inputs; the write is sent only once. A successful write followed by a disconnected monitor is reported as **sent, unverified**, not as a confirmed switch. Use **Refresh / reconnect** after returning to this computer or reconnecting a cable.

The tray clears input checkmarks after switching until a new read confirms the active input. **Open status log…** shows the complete latest result if the desktop truncates the menu text. The log is stored in the user cache directory (`~/.cache/dellkvm/status.log` by default). Quit cancels the active monitor command.

Avoid controlling the same monitor simultaneously from separate tray, CLI, or plugin instances. Requests are serialized within each application instance, not across processes. Monitor/driver behavior still requires testing on the target hardware.

## Windows

The Windows installer and portable ZIP contain `dellkvm-tray.exe` (no console window) and `dellkvm.exe` (CLI and optional TUI). Run the tray executable, choose a monitor when multiple are listed, and edit the starter config to match the VCP input-source codes of your monitor. `dellkvm.exe detect` shows numbered monitors. Both executables use the same configuration file.

Windows monitor access uses `EnumDisplayMonitors`, `GetPhysicalMonitorsFromHMONITOR`, and VCP code `0x60` through the Windows monitor configuration API. The API runs in a child process so the tray can end a stalled request. The Windows CI tests cover command parsing and building, but cannot prove that a specific monitor, cable, dock, and graphics driver accept the DDC/CI commands. Verify input switching on the target hardware before relying on it.

When exactly one monitor is detected, switching does not require a successful read of its current input. Some monitors return a duplicated input byte (for example, `0x0f0f` for `0x0f`); `dellkvm` accepts that reply. If a write succeeds but the monitor cannot be read afterward, the result says **sent, unverified**. The actual input codes still depend on the monitor model and must match `config.toml`.

## Testing

`make test`, `make lint`, `make build-all`, `make release-check`, and `sh tests/install_test.sh` are the local checks. GitHub Actions runs Go tests on Linux and Windows, cross-builds both operating systems for amd64 and arm64, and compiles both Windows installers on a Windows runner. To run the Linux test job locally with [act](https://nektosact.com/usage/), use `act workflow_dispatch -W .github/workflows/ci.yml -j test --matrix os:ubuntu-latest`. `act` uses Linux containers and does not replace the Windows runner.

## Linux launcher and autostart

The Linux install script sets up the launcher and XDG autostart entry. For a manual build, put the binary on `PATH` and use `integrations/linux/dellkvm.desktop` as a launcher or autostart entry. On Hyprland, an alternative is an `exec-once` entry launching `dellkvm tray`; use only one startup mechanism.

## Omarchy shell widget

The optional plugin targets the Quickshell-based Omarchy shell. It opens a popup with input buttons, the current input, and manual refresh. It invokes `dellkvm` with argument arrays and does not poll or switch automatically.

1. Install the binary on `PATH` and configure the monitor and input codes.
2. Copy `integrations/omarchy/` to `~/.config/omarchy/plugins/dellkvm.inputs/`.
3. Run `omarchy-shell shell rescanPlugins`, then `omarchy plugin enable dellkvm.inputs`.
4. If needed, position it with `omarchy bar move dellkvm.inputs --section right`.

The widget's `executable` setting accepts an absolute path to the binary. It uses the CLI's normal configuration lookup; use the user configuration path rather than relying on the shell's working directory. With multiple monitors, set `bus` in the config. No desktop settings are changed by building this project.
