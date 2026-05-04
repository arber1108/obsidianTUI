# Obsidian TUI

A terminal interface for browsing and reading your Obsidian vault.

## Requirements

- [Obsidian](https://obsidian.md/) with a vault open
- [Local REST API](https://github.com/coddingtonbear/obsidian-local-rest-api) community plugin enabled
- Your API key from the plugin settings (Settings → Community Plugins → Local REST API)

## Installation

Requires [Go](https://go.dev/dl/) installed.

```bash
git clone https://github.com/arber1108/obsidianTUI.git
cd obsidianTUI
```

```bash
# macOS / Linux
go build -o obsidianTUI .

# Windows
go build -o obsidianTUI.exe .
```

## Usage

```bash
# macOS / Linux
./obsidianTUI

# Windows
obsidianTUI.exe
```

Enter your Obsidian Local REST API key when prompted.

## Keybindings

| Key         | Action                  |
|-------------|-------------------------|
| `↑↓`        | Navigate files          |
| `Enter`     | Open file or folder     |
| `Backspace` | Go back                 |
| `o`         | Open file in Obsidian   |
| `q`         | Quit                    |
