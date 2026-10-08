param(
    [Parameter(Mandatory = $true)]
    [string]$Version
)

$ErrorActionPreference = 'Stop'
$projectRoot = (Resolve-Path (Join-Path $PSScriptRoot '../..')).Path
$setup = Join-Path $projectRoot "dist/installers/dellkvm_${Version}_windows_amd64_setup.exe"
$target = Join-Path $env:RUNNER_TEMP "dellkvm-smoke-$([Guid]::NewGuid().ToString('N'))"
$startup = Join-Path ([Environment]::GetFolderPath('Startup')) 'dellkvm.lnk'

if (-not (Test-Path $setup)) { throw "Installer is missing: $setup" }
if (Test-Path $startup) { throw "Unexpected existing startup shortcut: $startup" }

$arguments = "/VERYSILENT /SUPPRESSMSGBOXES /NORESTART /TASKS=autostart /DIR=`"$target`""
$process = Start-Process -FilePath $setup -ArgumentList $arguments -Wait -PassThru
if ($process.ExitCode -ne 0) { throw "Installer exited with $($process.ExitCode)" }
foreach ($name in @('dellkvm.exe', 'dellkvm-tray.exe', 'unins000.exe')) {
    if (-not (Test-Path (Join-Path $target $name))) { throw "Missing installed file: $name" }
}
if (-not (Test-Path $startup)) { throw 'Tray autostart shortcut was not created' }

$uninstaller = Join-Path $target 'unins000.exe'
$process = Start-Process -FilePath $uninstaller -ArgumentList '/VERYSILENT /SUPPRESSMSGBOXES /NORESTART' -Wait -PassThru
if ($process.ExitCode -ne 0) { throw "Uninstaller exited with $($process.ExitCode)" }
if (Test-Path $startup) { throw 'Tray autostart shortcut remained after uninstall' }
if (Test-Path (Join-Path $target 'dellkvm-tray.exe')) { throw 'Tray executable remained after uninstall' }
Write-Output 'Windows installer and autostart checks passed'
