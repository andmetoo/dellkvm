param(
    [Parameter(Mandatory = $true)]
    [ValidateSet('amd64', 'arm64')]
    [string]$Arch,
    [Parameter(Mandatory = $true)]
    [string]$Version
)

$ErrorActionPreference = 'Stop'
$projectRoot = (Resolve-Path (Join-Path $PSScriptRoot '../..')).Path
$stage = Join-Path $projectRoot "dist/installer-stage/$Arch"
$output = Join-Path $projectRoot 'dist/installers'
New-Item -ItemType Directory -Force -Path $stage, $output | Out-Null

$env:GOOS = 'windows'
$env:GOARCH = $Arch
$env:CGO_ENABLED = '0'
Push-Location $projectRoot
try {
    & go build -trimpath -ldflags '-s -w' -o (Join-Path $stage 'dellkvm.exe') ./cmd/dellkvm
    if ($LASTEXITCODE -ne 0) { throw 'Building Windows CLI failed' }
    & go build -trimpath -ldflags '-s -w -H=windowsgui' -o (Join-Path $stage 'dellkvm-tray.exe') ./cmd/dellkvm
    if ($LASTEXITCODE -ne 0) { throw 'Building Windows tray failed' }
    Copy-Item README.md, LICENSE, config.toml.default -Destination $stage

    $compiler = Get-Command ISCC.exe -ErrorAction SilentlyContinue
    if ($compiler) {
        $iscc = $compiler.Source
    } else {
        $iscc = @(
            'C:\Program Files (x86)\Inno Setup 6\ISCC.exe',
            'C:\Program Files\Inno Setup 6\ISCC.exe'
        ) | Where-Object { Test-Path $_ } | Select-Object -First 1
    }
    if (-not $iscc) { throw 'Inno Setup ISCC.exe was not found' }

    & $iscc "/DAppVersion=$Version" "/DArch=$Arch" "/DBuildDir=$stage" "/O$output" 'packaging/windows/dellkvm.iss'
    if ($LASTEXITCODE -ne 0) { throw 'Compiling Windows installer failed' }
} finally {
    Pop-Location
}
