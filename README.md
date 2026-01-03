<p align="center">
  <img src="docs/logo.png" alt="Sugu Logo" width="200">
</p>

# Sugu

> A JavaScript-like interpreted language that runs "sugu" (immediately)

## Documentation

📖 **[Language Documentation](https://n-yokomachi.github.io/sugu/)**

## Overview

Sugu is a simple interpreted language implemented in Go. It features JavaScript-like syntax and is ideal for learning and experimentation.

## Features

- JavaScript-like syntax
- Dynamic typing
- Simple variable declarations (`mut` / `const`)
- Functions and closures
- Control flow (if, switch, while, for)

## Installation

Download the latest release from [Releases](https://github.com/n-yokomachi/sugu/releases) and extract the binary.

### macOS (Apple Silicon)
```bash
unzip sugu-darwin-arm64.zip
chmod +x sugu
sudo mv sugu /usr/local/bin/
```

### Windows
Extract `sugu-windows-amd64.zip` and add the directory to your PATH.

## Usage

```bash
sugu              # Start REPL
sugu script.sugu  # Run a file
```

## License

MIT License
