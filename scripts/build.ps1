# Sugu CLI ビルドスクリプト (PowerShell)

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir
$OutputDir = Join-Path $ProjectRoot "dist"

Write-Host "Building Sugu CLI..."

# 出力ディレクトリを作成
if (-not (Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir | Out-Null
}

Push-Location $ProjectRoot
try {
    go build -o "$OutputDir/sugu.exe" .

    Write-Host "Created: $OutputDir/sugu.exe"
}
finally {
    Pop-Location
}
