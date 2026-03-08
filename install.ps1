# VibeDiff Windows Installation Script
# Usage: irm https://raw.githubusercontent.com/vibediff/vibediff/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

function Write-Success { Write-ColorOutput Green $args }
function Write-Info { Write-ColorOutput Cyan $args }
function Write-Error { Write-ColorOutput Red $args }

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "i386" }

# Get latest version
Write-Info "Detecting latest version..."
try {
    $latest = Invoke-RestMethod "https://api.github.com/repos/vibediff/vibediff/releases/latest"
    $version = $latest.tag_name
    Write-Success "Latest version: $version"
} catch {
    Write-Info "Could not detect version, using latest"
    $version = "latest"
}

# Determine install directory
$binDir = Join-Path $env:USERPROFILE "bin"
New-Item -ItemType Directory -Force -Path $binDir | Out-Null

# Download binary
$downloadUrl = "https://github.com/shubhamgurunglama007-oss/vibediff/releases/download/$version/vibediff-windows-$arch.exe"
$outputPath = Join-Path $binDir "vibediff.exe"

Write-Info "Downloading from $downloadUrl"
Invoke-WebRequest -Uri $downloadUrl -OutFile $outputPath

Write-Success ""
Write-Success "Successfully installed VibeDiff!"
Write-Success ""
Write-Info "Add to PATH (if not already):"
Write-Info "  [Environment]::SetEnvironmentVariable('Path', "`$env:Path;$binDir", 'User')"
Write-Success ""
Write-Info "Then restart your terminal and run:"
Write-Info "  vibediff version"
Write-Success ""
Write-Info "Uninstall with:"
Write-Info "  Remove-Item $outputPath"
