# Lambda 用ビルドスクリプト (PowerShell)
# Linux x86_64 向けにクロスコンパイルし、ZIP パッケージを作成する

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir
$OutputDir = Join-Path $ProjectRoot "dist"

Write-Host "Building Sugu Lambda..."

# 出力ディレクトリを作成
if (-not (Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir | Out-Null
}

# Linux x86_64 向けにビルド
Push-Location $ProjectRoot
try {
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    $env:CGO_ENABLED = "0"

    go build -o "$OutputDir/bootstrap" ./lambda

    Write-Host "Created: $OutputDir/bootstrap"

    # ZIP パッケージを作成
    $zipPath = Join-Path $OutputDir "sugu-lambda.zip"
    if (Test-Path $zipPath) {
        Remove-Item $zipPath
    }
    Compress-Archive -Path "$OutputDir/bootstrap" -DestinationPath $zipPath

    Write-Host "Created: $zipPath"
    Write-Host ""
    Write-Host "Deploy to AWS Lambda:"
    Write-Host "  Runtime: provided.al2023"
    Write-Host "  Architecture: x86_64"
    Write-Host "  Handler: bootstrap"
}
finally {
    # 環境変数をクリア
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
    Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
    Pop-Location
}
