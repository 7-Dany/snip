# 📦 SNIP - Code Snippet Manager

> A beautiful terminal-based code snippet manager built with Go and Bubble Tea

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/7-Dany/snip)](https://github.com/7-Dany/snip/releases)

## ✨ Features

- 🎨 **Beautiful TUI** - Interactive terminal interface with tabbed navigation
- 📁 **Category Management** - Organize snippets into logical categories
- 🏷️ **Tag System** - Multi-tag support for flexible organization
- 🔍 **Full-Text Search** - Quickly find snippets by title, description, or code
- ⌨️ **Syntax Highlighting** - Code editor with line numbers
- 💾 **SQLite Storage** - Fast, reliable local database
- 🚀 **Dual Interface** - Use interactive TUI or traditional CLI commands

## 📸 Screenshots

### Interactive TUI Mode

```
  ╭────────────────────╮╭────────────────────╮╭────────────────────╮╭────────────────────╮
  │       ⌨ Home       ││     Categories     ││        Tags        ││      Snippets      │
  │                    └┴────────────────────┴┴────────────────────┴┴────────────────────┤
  │                                                                                      │
  │    ████████                                                                          │
  │  ██        ██                                                                        │
  │  ██  ██  ██  ██                                                                      │
  │  ██  ██  ██  ██                                                                      │
  │  ██          ██                                                                      │
  │  ██  ████████ ██                                                                     │
  │  ████      ████                                                                      │
  │                                                                                      │
  │  SNIP - Code Snippet Manager                                                         │
  │                                                                                      │
  │  Welcome! Use the tabs above to navigate:                                            │
  │  - Home - You are here                                                               │
  │  - Categories - Manage snippet categories                                            │
  │  - Tags - Manage snippet tags                                                        │
  │  - Snippets - Manage code snippets                                                   │
  │                                                                                      │
  │  Navigation:                                                                         │
  │    Ctrl+F      - Toggle between Interactive and Navigation mode                      │
  │    Tab         - Next tab (in Navigation mode)                                       │
  │    Shift+Tab   - Previous tab (in Navigation mode)                                   │
  │    ←/→         - Switch tabs (in Navigation mode)                                    │
  │    PgUp/PgDn   - Scroll content (in Interactive mode)                                │
  │    Ctrl+←/→    - Scroll horizontally                                                 │
  │    Ctrl+C      - Quit                                                                │
  │                                                                                      │
  │  Mode Indicators:                                                                    │
  │    ⌨ (green)   - Interactive mode: Components receive keyboard input                 │
  │    ⇄ (blue)    - Navigation mode: Tab switching and scrolling enabled                │
  └──────────────────────────────────────────────────────────────────────────────────────┘

  Ctrl+F: toggle mode | Ctrl+C: quit
```

## 🚀 Quick Start

### Installation

#### Download Pre-built Binary

Download the latest release for your platform from the [releases page](https://github.com/7-Dany/snip/releases).

**Linux/macOS:**
```bash
# Download and install (replace VERSION with actual version)
curl -L https://github.com/7-Dany/snip/releases/download/vVERSION/snip-linux-amd64 -o snip
chmod +x snip
sudo mv snip /usr/local/bin/
```

**Windows (PowerShell):**
```powershell
# Download from releases page and add to PATH
```

#### Build from Source

```bash
# Clone the repository
git clone https://github.com/7-Dany/snip.git
cd snip

# Build and install
make build
sudo make install

# Or use go install
go install
```

### First Run

```bash
# Launch interactive TUI
snip

# Or use CLI commands, note: to work with commands snip must be followed with arguments.
snip snippet create
snip snippet list
snip help
```

## 📖 Usage

### Interactive TUI Mode

The recommended way to use SNIP is through the interactive TUI:

```bash
snip
```

**Keyboard Shortcuts:**

| Key | Action |
|-----|--------|
| `Ctrl+F` | Toggle between Interactive and Navigation mode |
| `Tab` / `Shift+Tab` | Navigate between tabs (Navigation mode) |
| `←` / `→` | Switch tabs (Navigation mode) |
| `↑` / `↓` | Navigate lists/menus |
| `Enter` | Select item / Confirm action |
| `/` | Start search/filter |
| `a` | Add new item |
| `r` | Refresh list |
| `?` | Show help |
| `Esc` | Cancel / Go back |
| `Ctrl+C` | Quit application |

**In Code Editor:**

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Navigate between fields |
| `Alt+C` | Select category |
| `Alt+T` | Manage tags |
| `Ctrl+S` | Save snippet |
| `Esc` | Cancel editing |

### CLI Commands

#### Snippet Management

```bash
# Create a new snippet interactively
snip snippet create

# List all snippets
snip snippet list

# List snippets with filters
snip snippet list --language go
snip snippet list --category 1
snip snippet list --tag 2

# Show a specific snippet
snip snippet show 5

# Update a snippet
snip snippet update 5

# Delete a snippet
snip snippet delete 5

# Search snippets
snip snippet search "binary tree"
```

#### Category Management

```bash
# Create a category
snip category create algorithms
snip category create  # Interactive mode

# List all categories
snip category list

# Delete a category
snip category delete 3
```

#### Tag Management

```bash
# Create a tag
snip tag create performance
snip tag create  # Interactive mode

# List all tags
snip tag list

# Delete a tag
snip tag delete 7
```

#### Help

```bash
# General help
snip help

# Topic-specific help
snip help snippet
snip help category
snip help tag
```

## 🏗️ Architecture

```
snip/
├── cmd/                   # Application entry points
│   └── main.go            # Main CLI application
├── internal/
│   ├── cli/               # CLI layer
│   │   ├── commands/      # CLI command handlers
│   │   ├── components/    # Reusable Bubble Tea UI components
│   │   └── tui/           # Terminal UI implementation
│   ├── domain/            # Business logic and entities
│   └── storage/           # Data persistence layer
└── main.go
```

### Key Components

- **Domain Layer**: Pure business logic (snippets, categories, tags)
- **Storage Layer**: JSON-based repositories with transaction support
- **CLI Commands**: Traditional command-line interface
- **TUI**: Interactive terminal interface using Bubble Tea
- **Components**: Reusable UI widgets (tables, editors, menus, dialogs)

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean
```

### Running Tests

```bash
go test ./...
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) - amazing TUI framework
- Uses [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling
- Inspired by the need for a simple, elegant snippet manager

## 📬 Contact

- GitHub: [@7-Dany](https://github.com/7-Dany)
- Project Link: [https://github.com/7-Dany/snip](https://github.com/7-Dany/snip)

---

⭐ If you find this project useful, please consider giving it a star!
