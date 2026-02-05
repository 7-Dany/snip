# SNIP - Code Snippet Manager

> Complete technical documentation and project context for AI assistants.
> This document contains all essential information about the SNIP project's architecture, design decisions, implementation details, and development guidelines.

## Quick Reference

- **Project Name**: SNIP
- **Type**: Terminal-based code snippet manager
- **Language**: Go 1.24+
- **Repository**: https://github.com/7-Dany/snip
- **Author**: 7-Dany (ali.ali.ameen245@gmail.com)
- **License**: MIT
- **Current Version**: v0.1.1
- **Latest Commit**: `1d8fead` - refactor(cli): replace table view with searchable table view

## Table of Contents

1. [Project Overview](#project-overview)
2. [Tech Stack](#tech-stack)
3. [Architecture](#architecture)
4. [Domain Layer](#domain-layer)
5. [Storage Layer](#storage-layer)
6. [CLI Layer](#cli-layer)
7. [TUI Layer](#tui-layer)
8. [Configuration](#configuration)
9. [Testing Strategy](#testing-strategy)
10. [Development Guidelines](#development-guidelines)
11. [Build and Release](#build-and-release)

---

## Project Overview

SNIP is a beautiful terminal-based code snippet manager built with Go and Bubble Tea. It provides both an interactive TUI (Terminal User Interface) and traditional CLI commands for managing code snippets, categories, and tags.

### Key Features

- 🎨 **Beautiful TUI** - Interactive terminal interface with tabbed navigation
- 📁 **Category Management** - Organize snippets into logical categories
- 🏷️ **Tag System** - Multi-tag support for flexible organization
- 🔍 **Full-Text Search** - Search across title, description, and code
- ⌨️ **Syntax Highlighting** - Code editor with line numbers
- 🚀 **Dual Interface** - Both interactive TUI and CLI commands
- 💾 **JSON Storage** - Simple, human-readable data format
- 🧪 **Comprehensive Tests** - ~85% overall coverage

### Project Status

✅ **Completed**:
- Domain layer (100% coverage)
- Storage layer (99.5% coverage)
- CLI commands (85% coverage)
- Interactive TUI with tabbed navigation
- All core features implemented

🔄 **Active Branch**: `main`
📦 **Latest Release**: v0.1.1

---

## Tech Stack

### Core Dependencies

```go
// UI Framework
github.com/charmbracelet/bubbletea v1.3.10    // TUI framework
github.com/charmbracelet/bubbles v0.21.0      // TUI components
github.com/charmbracelet/lipgloss v1.1.0      // Styling

// CLI & Output
github.com/fatih/color v1.18.0                // Terminal colors
github.com/jedib0t/go-pretty/v6 v6.7.8        // Tables
```

### Development Tools

- **Go**: 1.24+
- **Make**: Build automation
- **Git**: Version control
- **Testing**: Go standard library

---

## Architecture

### Project Structure

```
snip/
├── cmd/
│   └── main.go                  # CLI entry point & TUI launcher
├── internal/
│   ├── cli/
│   │   ├── commands/            # CLI command handlers
│   │   │   ├── commands.go      # Main CLI coordinator
│   │   │   ├── snippet.go       # Snippet commands
│   │   │   ├── category.go      # Category commands
│   │   │   ├── tag.go           # Tag commands
│   │   │   ├── help.go          # Help system
│   │   │   ├── output.go        # Display utilities
│   │   │   └── input_helpers.go # Interactive prompts
│   │   ├── components/          # Reusable Bubble Tea components
│   │   │   ├── tabs.go          # Tab navigation
│   │   │   ├── searchable_table_view.go
│   │   │   ├── form_view.go     # Multi-field forms
│   │   │   ├── code_editor.go   # Code editing
│   │   │   ├── menu_view.go     # Menu selections
│   │   │   └── selector_view.go # Item selection
│   │   ├── tui/                 # Terminal UI screens
│   │   │   ├── home_tab.go
│   │   │   ├── categories_tab.go
│   │   │   ├── tags_tab.go
│   │   │   └── snippets_tab.go
│   │   └── config/              # Configuration management
│   │       └── config.go
│   ├── domain/                  # Business entities
│   │   ├── snippet.go
│   │   ├── category.go
│   │   ├── tag.go
│   │   ├── errors.go
│   │   └── repository.go        # Interface definitions
│   └── storage/                 # Data persistence
│       ├── repositories.go      # Public API
│       ├── internal_store.go    # Shared data store
│       ├── search.go            # Search functionality
│       ├── snippet_repository.go
│       ├── category_repository.go
│       └── tag_repository.go
├── main.go                      # Application entry point
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── DOC.md                       # This file
└── LICENSE
```

### Layer Responsibilities

**Entry Point** (`cmd/main.go`)
- Loads configuration from `~/.snip/config.json`
- Initializes storage repositories
- Routes to TUI (no args) or CLI (with args)
- Ensures data is saved on exit

**Domain Layer** (`internal/domain/`)
- Pure business logic and entities
- No external dependencies (standard library only)
- Immutable entities with validation
- Repository interfaces

**Storage Layer** (`internal/storage/`)
- JSON file persistence
- Repository pattern with shared internal store
- Thread-safe operations (RWMutex)
- Full-text search capabilities

**CLI Layer** (`internal/cli/commands/`)
- Command-line interface handlers
- Input validation and prompts
- Formatted output and error messages
- Interactive flows using Bubble Tea

**TUI Layer** (`internal/cli/tui/` + `components/`)
- Interactive terminal user interface
- Tabbed navigation system
- Reusable UI components
- Code editor with syntax highlighting

**Configuration** (`internal/cli/config/`)
- Manages `~/.snip/config.json`
- Auto-creates config on first run
- Stores data file path

---

## Domain Layer

Location: `internal/domain/`

### Design Principles

1. **Encapsulation**: All fields are unexported (private)
2. **Immutability**: Getters return values, not pointers
3. **Defensive Copying**: Slices are copied when returned
4. **Validation at Boundaries**: All constructors and setters validate
5. **Timestamps**: All entities track `createdAt` and `updatedAt`
6. **JSON Serialization**: Custom Marshal/Unmarshal for unexported fields

### Entity: Category

Organizes snippets into logical groups.

**Fields:**
```go
id        int       // Unique identifier (set by storage layer)
name      string    // Category name (required, non-empty)
createdAt time.Time // Creation timestamp
updatedAt time.Time // Last modification timestamp
```

**Constructor:**
```go
NewCategory(name string) (*Category, error)
// Returns ErrEmptyName if name is empty
```

**Methods:**
```go
// Getters
ID() int
Name() string
CreatedAt() time.Time
UpdatedAt() time.Time

// Setters
SetName(name string) error    // Validates, updates updatedAt
SetID(id int)                 // Storage layer only

// Utility
String() string               // "Category{id=1, name="Go"}"
Equal(other *Category) bool   // Deep equality check
```

**Validation Rules:**
- Name cannot be empty (returns `ErrEmptyName`)
- Whitespace is preserved (trimming done in CLI layer)

### Entity: Tag

Labels for flexible snippet categorization.

**Fields:**
```go
id        int       // Unique identifier
name      string    // Tag name (required, non-empty)
createdAt time.Time
updatedAt time.Time
```

**Constructor:**
```go
NewTag(name string) (*Tag, error)
// Returns ErrEmptyName if name is empty
```

**Methods:**
```go
// Same interface as Category
ID(), Name(), CreatedAt(), UpdatedAt()
SetName(name string) error
SetID(id int)
String() string
Equal(other *Tag) bool
```

### Entity: Snippet

Code snippet with rich metadata.

**Fields:**
```go
id          int       // Unique identifier
title       string    // Snippet title (required)
language    string    // Programming language (required)
code        string    // Code content (required)
description string    // Optional description
categoryID  int       // Category ID (0 if uncategorized)
tags        []int     // Tag IDs (never nil, always []int{})
createdAt   time.Time
updatedAt   time.Time
```

**Constructor:**
```go
NewSnippet(title, language, code string) (*Snippet, error)
// Returns ErrEmptyTitle, ErrEmptyLanguage, or ErrEmptyCode
// Sets createdAt and updatedAt to current time
// Initializes tags as empty slice []int{}
```

**Methods:**
```go
// Getters
ID() int
Title() string
Language() string
Code() string
Description() string
CategoryID() int
Tags() []int              // Returns defensive copy
CreatedAt() time.Time
UpdatedAt() time.Time

// Setters (with validation)
SetTitle(title string) error
SetLanguage(language string) error
SetCode(code string) error
SetDescription(description string)  // No validation
SetCategory(catID int)               // No validation
SetID(id int)                        // Storage layer only

// Tag Management
AddTag(tagID int)           // Prevents duplicates, silently ignores
RemoveTag(tagID int)        // No-op if tag doesn't exist
HasTag(tagID int) bool      // Check tag membership

// Utility
String() string
Equal(other *Snippet) bool
```

**Validation Rules:**
- Title, language, and code cannot be empty
- Tags array is never nil (normalized to `[]int{}` on unmarshal)
- AddTag prevents duplicates
- All setters update `updatedAt` timestamp

### Domain Errors

```go
var (
    ErrEmptyName     = errors.New("name cannot be empty")
    ErrEmptyTitle    = errors.New("title cannot be empty")
    ErrEmptyLanguage = errors.New("language cannot be empty")
    ErrEmptyCode     = errors.New("code cannot be empty")
)
```

### Repository Interfaces

```go
type CategoryRepository interface {
    List() ([]*Category, error)
    FindByID(id int) (*Category, error)
    FindByName(name string) (*Category, error)
    Create(category *Category) error
    Update(category *Category) error
    Delete(id int) error
}

type TagRepository interface {
    List() ([]*Tag, error)
    FindByID(id int) (*Tag, error)
    FindByName(name string) (*Tag, error)
    Create(tag *Tag) error
    Update(tag *Tag) error
    Delete(id int) error
}

type SnippetRepository interface {
    List() ([]*Snippet, error)
    FindByID(id int) (*Snippet, error)
    Search(query string) ([]*Snippet, error)
    FindByLanguage(language string) ([]*Snippet, error)
    FindByCategory(categoryID int) ([]*Snippet, error)
    FindByTag(tagID int) ([]*Snippet, error)
    Create(snippet *Snippet) error
    Update(snippet *Snippet) error
    Delete(id int) error
}
```

### JSON Marshaling Pattern

All entities use custom JSON marshaling to handle unexported fields:

```go
func (e *Entity) MarshalJSON() ([]byte, error) {
    return json.Marshal(&struct {
        ID        int       `json:"id"`
        Name      string    `json:"name"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
    }{
        ID:        e.id,
        Name:      e.name,
        CreatedAt: e.createdAt,
        UpdatedAt: e.updatedAt,
    })
}

func (e *Entity) UnmarshalJSON(data []byte) error {
    aux := &struct {
        ID        int       `json:"id"`
        Name      string    `json:"name"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
    }{}
    if err := json.Unmarshal(data, aux); err != nil {
        return err
    }
    if aux.Name == "" {
        return ErrEmptyName
    }
    e.id = aux.ID
    e.name = aux.Name
    e.createdAt = aux.CreatedAt
    e.updatedAt = aux.UpdatedAt
    return nil
}
```

---

## Storage Layer

Location: `internal/storage/`

### Architecture

**Design Pattern**: Repository pattern with shared internal store

**Key Features**:
- Atomic file saves (temp file + rename)
- Auto-incrementing IDs
- Thread-safe operations (RWMutex)
- Full-text search
- Defensive copying

### File Structure

```
storage/
├── repositories.go              # Public API
├── internal_store.go            # Shared data structure
├── search.go                    # Search index
├── snippet_repository.go        # Snippet operations
├── category_repository.go       # Category operations
└── tag_repository.go            # Tag operations
```

### Public API

**Repositories Struct:**
```go
type Repositories struct {
    Snippets   domain.SnippetRepository
    Categories domain.CategoryRepository
    Tags       domain.TagRepository
    store      *store  // Internal, unexported
}
```

**Constructor:**
```go
func New(filepath string) *Repositories
// Creates repositories with shared internal store
// filepath: Path to JSON data file (e.g., ~/.snip/snippets.json)
```

**Methods:**
```go
func (r *Repositories) Save() error
// Atomically saves all data to JSON file
// Uses temp file + rename for safety

func (r *Repositories) Load() error
// Loads all data from JSON file
// Not an error if file doesn't exist (creates empty store)
// Normalizes nil slices to empty slices
```

**Usage Example:**
```go
repos := storage.New("~/.snip/snippets.json")
if err := repos.Load(); err != nil {
    log.Fatal(err)
}
defer repos.Save()

snippet, _ := domain.NewSnippet("quicksort", "go", "func quicksort() {}")
repos.Snippets.Create(snippet)
```

### Internal Store

**Structure:**
```go
type store struct {
    filepath   string
    mu         sync.RWMutex
    snippets   map[int]*domain.Snippet
    categories map[int]*domain.Category
    tags       map[int]*domain.Tag
    nextSnippetID  int
    nextCategoryID int
    nextTagID      int
}
```

**JSON Format:**
```json
{
  "snippets": [...],
  "categories": [...],
  "tags": [...],
  "next_snippet_id": 10,
  "next_category_id": 5,
  "next_tag_id": 8
}
```

**Key Methods:**
```go
func (s *store) save() error
// Atomic save: marshal → write to temp → rename
// Ensures data integrity even on crash

func (s *store) load() error
// Loads JSON file if exists
// Normalizes nil slices to empty slices
// Creates empty maps if file doesn't exist

func (s *store) nextID(idType string) int
// Auto-increments and returns next available ID
```

### Search Functionality

Location: `internal/storage/search.go`

**Methods:**
```go
func (s *store) search(query string) []*domain.Snippet
// Full-text search across title, language, code, description
// Case-insensitive
// Returns empty slice for no matches
// Returns nil only for empty query

func (s *store) findByLanguage(language string) []*domain.Snippet
// Filters by exact language match (case-sensitive)

func (s *store) findByCategory(categoryID int) []*domain.Snippet
// Filters by category ID

func (s *store) findByTag(tagID int) []*domain.Snippet
// Filters by tag ID (checks Tags() slice)
```

**Search Implementation:**
```go
// Converts to lowercase for case-insensitive search
lowerQuery := strings.ToLower(query)

// Searches across multiple fields
matches := strings.Contains(strings.ToLower(s.Title()), lowerQuery) ||
           strings.Contains(strings.ToLower(s.Language()), lowerQuery) ||
           strings.Contains(strings.ToLower(s.Code()), lowerQuery) ||
           strings.Contains(strings.ToLower(s.Description()), lowerQuery)
```

### Repository Operations

**Common Pattern:**
```go
// List - Returns defensive copy of all entities
func (r *repository) List() ([]*Entity, error) {
    r.store.mu.RLock()
    defer r.store.mu.RUnlock()
    
    entities := make([]*Entity, 0, len(r.store.entities))
    for _, e := range r.store.entities {
        entities = append(entities, e)
    }
    return entities, nil
}

// FindByID - Returns entity or ErrNotFound
func (r *repository) FindByID(id int) (*Entity, error) {
    r.store.mu.RLock()
    defer r.store.mu.RUnlock()
    
    entity, exists := r.store.entities[id]
    if !exists {
        return nil, ErrNotFound
    }
    return entity, nil
}

// Create - Auto-assigns ID and stores
func (r *repository) Create(entity *Entity) error {
    r.store.mu.Lock()
    defer r.store.mu.Unlock()
    
    id := r.store.nextID("entity")
    entity.SetID(id)
    r.store.entities[id] = entity
    return nil
}

// Update - Validates entity exists
func (r *repository) Update(entity *Entity) error {
    r.store.mu.Lock()
    defer r.store.mu.Unlock()
    
    if _, exists := r.store.entities[entity.ID()]; !exists {
        return ErrNotFound
    }
    r.store.entities[entity.ID()] = entity
    return nil
}

// Delete - Idempotent (only error if not found)
func (r *repository) Delete(id int) error {
    r.store.mu.Lock()
    defer r.store.mu.Unlock()
    
    if _, exists := r.store.entities[id]; !exists {
        return ErrNotFound
    }
    delete(r.store.entities, id)
    return nil
}
```

### Storage Errors

```go
var (
    ErrNotFound      = errors.New("entity not found")
    ErrDuplicateName = errors.New("entity with this name already exists")
)
```

### Thread Safety

All operations use `sync.RWMutex`:
- **Read operations**: `RLock()` / `RUnlock()`
- **Write operations**: `Lock()` / `Unlock()`

This prepares for future concurrent usage without requiring refactoring.

---

## CLI Layer

Location: `internal/cli/commands/`

### File Structure

```
commands/
├── commands.go              # Main CLI coordinator
├── snippet.go               # Snippet command handler
├── category.go              # Category command handler
├── tag.go                   # Tag command handler
├── help.go                  # Help system
├── output.go                # Display utilities
├── input_helpers.go         # Interactive prompts
└── testing_helpers.go       # Test utilities
```

### CLI Coordinator

**Structure:**
```go
type CLI struct {
    snippet  *SnippetCommand
    category *CategoryCommand
    tag      *TagCommand
    help     *HelpCommand
}

func NewCLI(repos *storage.Repositories) *CLI
```

**Run Method:**
```go
func (c *CLI) Run(args []string) {
    if len(args) < 2 {
        c.snippet.manage([]string{}) // Default to snippet
        return
    }

    topic := strings.ToLower(args[1])
    topicArgs := args[2:]

    switch topic {
    case "snippet", "s":
        c.snippet.manage(topicArgs)
    case "category", "cat", "c":
        c.category.manage(topicArgs)
    case "tag", "t":
        c.tag.manage(topicArgs)
    case "help", "h":
        c.help.show(topicArgs)
    default:
        // Backward compatibility: treat as snippet command
        c.snippet.manage(args[1:])
    }
}
```

### Command Handler Pattern

All command handlers follow this structure:

```go
type XCommand struct {
    repos *storage.Repositories
}

func NewXCommand(repos *storage.Repositories) *XCommand {
    return &XCommand{repos: repos}
}

func (xc *XCommand) manage(args []string) {
    if len(args) == 0 {
        PrintError("No subcommand provided")
        return
    }

    subcommand := strings.ToLower(args[0])
    subcommandArgs := args[1:]

    switch subcommand {
    case "list", "l":
        xc.list()
    case "create", "c":
        xc.create(subcommandArgs)
    case "delete", "d":
        xc.delete(subcommandArgs)
    default:
        PrintError(fmt.Sprintf("Unknown command '%s'", args[0]))
    }
}
```

### Command Operations

#### List Operation

**Purpose**: Display all entities in formatted table

**Implementation:**
```go
func (xc *XCommand) list() {
    entities, err := xc.repos.Entities.List()
    if err != nil {
        PrintError("Failed to list entities")
        return
    }

    if len(entities) == 0 {
        PrintInfo("No entities found")
        return
    }

    // Create table with go-pretty
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.AppendHeader(table.Row{"ID", "Name", "Created At"})

    for _, e := range entities {
        t.AppendRow(table.Row{
            e.ID(),
            e.Name(),
            e.CreatedAt().Format("2006-01-02 15:04"),
        })
    }

    t.Render()
}
```

**Features:**
- Shows informative message when empty
- Uses `go-pretty` for table rendering
- Sorts by ID (insertion order)
- Snippets support filtering: `--category`, `--tag`, `--language`

#### Create Operation

**Purpose**: Add new entity with validation

**Modes:**
1. **Argument Mode**: `snip category create "My Category"`
2. **Interactive Mode**: `snip category create` (prompts user)

**Implementation:**
```go
func (xc *XCommand) create(args []string) {
    var name string

    if len(args) > 0 {
        // Argument mode
        name = strings.TrimSpace(args[0])
    } else {
        // Interactive mode
        var err error
        name, err = promptForName()
        if err != nil {
            PrintError("Input cancelled")
            return
        }
    }

    if name == "" {
        PrintError("Name cannot be empty")
        return
    }

    // Check for duplicates
    existing, _ := xc.repos.Entities.FindByName(name)
    if existing != nil {
        PrintError(fmt.Sprintf("Entity '%s' already exists", name))
        return
    }

    // Create entity
    entity, err := domain.NewEntity(name)
    if err != nil {
        PrintError(err.Error())
        return
    }

    if err := xc.repos.Entities.Create(entity); err != nil {
        PrintError("Failed to create entity")
        return
    }

    PrintSuccess(fmt.Sprintf("Created entity '%s' (ID: %d)", name, entity.ID()))
}
```

**Validation:**
- Trims whitespace from input
- Checks for empty names
- Prevents duplicate names
- Shows clear error messages

#### Update Operation

**Purpose**: Modify existing entity

**Implementation:**
```go
func (sc *SnippetCommand) update(args []string) {
    // Validate ID
    if len(args) == 0 {
        PrintError("Please provide snippet ID")
        return
    }

    id, err := strconv.Atoi(args[0])
    if err != nil {
        PrintError("Invalid snippet ID")
        return
    }

    // Fetch existing snippet
    snippet, err := sc.repos.Snippets.FindByID(id)
    if err != nil {
        PrintError(fmt.Sprintf("Snippet %d not found", id))
        return
    }

    // Interactive form with pre-filled values
    updated, err := promptForSnippet(snippet)
    if err != nil {
        PrintError("Update cancelled")
        return
    }

    // Update in database
    if err := sc.repos.Snippets.Update(updated); err != nil {
        PrintError("Failed to update snippet")
        return
    }

    PrintSuccess(fmt.Sprintf("Updated snippet %d", id))
}
```

#### Delete Operation

**Purpose**: Remove entity with confirmation

**Implementation:**
```go
func (xc *XCommand) delete(args []string) {
    if len(args) == 0 {
        PrintError("Please provide entity ID")
        return
    }

    id, err := strconv.Atoi(args[0])
    if err != nil {
        PrintError("Invalid entity ID")
        return
    }

    // Confirm deletion
    fmt.Printf("Delete entity %d? (y/n): ", id)
    var confirm string
    fmt.Scanln(&confirm)

    if strings.ToLower(confirm) != "y" {
        PrintInfo("Deletion cancelled")
        return
    }

    // Delete entity
    if err := xc.repos.Entities.Delete(id); err != nil {
        if errors.Is(err, storage.ErrNotFound) {
            PrintError(fmt.Sprintf("Entity %d not found", id))
        } else {
            PrintError("Failed to delete entity")
        }
        return
    }

    PrintSuccess(fmt.Sprintf("Deleted entity %d", id))
}
```

### Output Functions

Location: `internal/cli/commands/output.go`

**Available Functions:**
```go
func PrintLogo()                    // ASCII art branding
func PrintSuccess(msg string)       // ✓ Green message
func PrintError(msg string)         // ✗ Red message
func PrintInfo(msg string)          // ⓘ Cyan message
func PrintCommand(command string)   // ▶ Gray message
func PrintLoading(msg string)       // Spinner animation
func ClearScreen()                  // Clear terminal
```

**Color Scheme:**
- Success: Green + Bold
- Error: Red + Bold
- Info: Cyan
- Command: Gray
- Logo: Magenta/Yellow

### Interactive Prompts

Location: `internal/cli/commands/input_helpers.go`

Uses Bubble Tea for interactive input:

**Simple Text Input:**
```go
func promptForName() (string, error) {
    p := tea.NewProgram(initialNameInputModel())
    m, err := p.Run()
    if err != nil {
        return "", err
    }

    model := m.(nameInputModel)
    if model.cancelled {
        return "", errors.New("cancelled")
    }
    return model.textInput.Value(), nil
}
```

**Multi-Field Form:**
```go
func promptForSnippet(existing *domain.Snippet) (*domain.Snippet, error) {
    // Create form with pre-filled values
    form := components.NewFormView([]components.FormField{
        {Label: "Title", Value: existing.Title()},
        {Label: "Language", Value: existing.Language()},
        {Label: "Description", Value: existing.Description()},
    })

    p := tea.NewProgram(form)
    m, err := p.Run()
    // ... process result
}
```

**Key Bindings:**
- `Enter`: Confirm input
- `Esc` / `Ctrl+C`: Cancel operation
- `Tab` / `Shift+Tab`: Navigate fields (in forms)

### Help System

Location: `internal/cli/commands/help.go`

**Structure:**
```go
type HelpCommand struct {
    repos *storage.Repositories
}

func (h *HelpCommand) show(args []string) {
    if len(args) == 0 {
        h.showGeneral()
        return
    }

    topic := strings.ToLower(args[0])
    switch topic {
    case "snippet", "s":
        h.showSnippet()
    case "category", "cat", "c":
        h.showCategory()
    case "tag", "t":
        h.showTag()
    default:
        PrintError(fmt.Sprintf("Unknown topic '%s'", topic))
        h.showGeneral()
    }
}
```

**Output Format:**
- General help: Lists all topics
- Topic-specific: Shows commands and examples
- Includes aliases (e.g., `s` for `snippet`)

---

## TUI Layer

Location: `internal/cli/tui/` and `internal/cli/components/`

### Architecture

**Components:**
- **Tabs** (`components/tabs.go`): Tabbed navigation system
- **Home** (`tui/home_tab.go`): Welcome screen with instructions
- **Categories** (`tui/categories_tab.go`): Category management
- **Tags** (`tui/tags_tab.go`): Tag management
- **Snippets** (`tui/snippets_tab.go`): Snippet management

### Reusable Components

#### Tabs Component

```go
type Tabs struct {
    labels      []string
    tabs        []TabModel
    activeTab   int
    focusMode   FocusMode  // Interactive vs Navigation
    labelWidth  int
}

const (
    FocusModeInteractive FocusMode = iota  // Tab content receives input
    FocusModeNavigation                     // Tab switching enabled
)
```

**Key Bindings:**
- `Ctrl+F`: Toggle focus mode
- `Tab` / `Shift+Tab`: Navigate tabs (Navigation mode)
- `←` / `→`: Switch tabs (Navigation mode)
- `Ctrl+C`: Quit

**Features:**
- Mode indicator (⌨ Interactive, ⇄ Navigation)
- Smooth tab transitions
- Delegates input to active tab in Interactive mode
- Custom label width per tab

#### Searchable Table View

```go
type SearchableTableView struct {
    table       table.Model
    searchInput textinput.Model
    searching   bool
    allData     []table.Row
}
```

**Key Bindings:**
- `/`: Enter search mode
- `Esc`: Exit search mode
- `↑` / `↓`: Navigate rows
- `Enter`: Select row

**Features:**
- Real-time filtering
- Case-insensitive search
- Preserves all data for reset
- Integrated with table navigation

#### Form View

```go
type FormView struct {
    fields      []FormField
    activeField int
    submitted   bool
    cancelled   bool
}

type FormField struct {
    Label    string
    Value    string
    Input    textinput.Model
}
```

**Key Bindings:**
- `Tab` / `Shift+Tab`: Navigate fields
- `Enter`: Submit form
- `Esc`: Cancel

**Features:**
- Multi-field input
- Pre-filled values for editing
- Validation support
- Visual focus indicators

#### Code Editor

```go
type CodeEditor struct {
    textarea textarea.Model
    language string
    saved    bool
}
```

**Key Bindings:**
- All textarea bindings
- `Ctrl+S`: Save
- `Esc`: Cancel

**Features:**
- Syntax highlighting (via language hint)
- Line numbers
- Multi-line editing
- Save/cancel state

### Tab Implementations

#### Home Tab

**Purpose**: Welcome screen and navigation guide

**Content:**
- ASCII art logo
- Feature list
- Keyboard shortcuts guide
- Mode explanations

#### Categories Tab

**Features:**
- List all categories in table
- Create new category
- Delete category with confirmation
- Search/filter categories
- Refresh list

**Key Bindings:**
- `a`: Add category
- `d`: Delete selected
- `r`: Refresh list
- `/`: Search

#### Tags Tab

**Features:**
- List all tags in table
- Create new tag
- Delete tag with confirmation
- Search/filter tags
- Refresh list

**Key Bindings:**
- Same as Categories Tab

#### Snippets Tab

**Features:**
- List all snippets with metadata
- Create new snippet (multi-step)
- View snippet details
- Edit snippet
- Delete snippet
- Filter by category, tag, language
- Full-text search

**Key Bindings:**
- `a`: Add snippet
- `v`: View selected snippet
- `e`: Edit selected snippet
- `d`: Delete selected
- `/`: Search
- `r`: Refresh

**Snippet Creation Flow:**
1. Enter title, language, description (form)
2. Edit code (code editor)
3. Select category (menu)
4. Select tags (multi-select)
5. Save to database

### TUI Entry Point

Location: `cmd/main.go`

```go
func Run() {
    // Load config
    config, err := config.LoadConfig()
    handleError(err)

    // Initialize storage
    repos := storage.New(config.StoragePath)
    repos.Load()
    defer repos.Save()

    // Create CLI for command mode
    app := commands.NewCLI(repos)

    if len(os.Args) > 1 {
        // CLI mode: run command and exit
        app.Run(os.Args)
        return
    }

    // TUI mode: create tabs and run
    tabs := components.NewTabs(
        []string{"Home", "Categories", "Tags", "Snippets"},
        []components.TabModel{
            tui.NewHomeTab(),
            tui.NewCategoriesTab(repos),
            tui.NewTagsTab(repos),
            tui.NewSnippetsTab(repos),
        },
    )

    p := tea.NewProgram(tabs, tea.WithAltScreen())
    p.Run()
}
```

---

## Configuration

Location: `internal/cli/config/`

### Config File

**Path**: `~/.snip/config.json`

**Format:**
```json
{
  "storage_path": "/home/user/.snip/snippets.json"
}
```

### Behavior

**First Run:**
1. Checks for `~/.snip/config.json`
2. If not found, creates:
   - `~/.snip/` directory
   - `config.json` with default values
   - `snippets.json` (empty, created on first save)

**Default Values:**
- `storage_path`: `~/.snip/snippets.json`

### API

```go
func LoadConfig() (*Config, error)
// Loads config from ~/.snip/config.json
// Creates file with defaults if missing
// Returns error if:
//   - Cannot determine home directory
//   - Cannot create .snip directory
//   - Cannot read/write config file
//   - JSON is invalid
```

---

## Testing Strategy

### Test Organization

**Pattern**: Use `t.Run()` for subtests (NOT table-driven tests)

**Structure:**
```go
func TestEntity_Method(t *testing.T) {
    t.Run("success case description", func(t *testing.T) {
        // Arrange
        entity := mustCreateEntity(t, "params")
        
        // Act
        result := entity.Method()
        
        // Assert
        if result != expected {
            t.Errorf("expected %v, got %v", expected, result)
        }
    })

    t.Run("error case description", func(t *testing.T) {
        // Test error scenarios
    })
}
```

### Coverage by Layer

| Layer | Coverage | Files Tested |
|-------|----------|--------------|
| **Domain** | ~100% | All entity files |
| **Storage** | ~99.5% | All repository files |
| **CLI Commands** | ~85% | All command handlers |
| **TUI Components** | Manual | Interactive components |
| **Overall** | **~85%** | Excluding TUI |

### What We Test

✅ **Domain Layer:**
- Constructor validation
- Setter validation
- Getter correctness
- JSON marshaling/unmarshaling
- Equal() method
- String() method
- Edge cases (nil, empty, zero values)

✅ **Storage Layer:**
- CRUD operations
- Search functionality
- Error cases (not found, duplicates)
- Thread safety
- Defensive copying
- Save/load round trips
- Nil slice normalization

✅ **CLI Layer:**
- Command routing
- Input validation
- Error handling
- Case insensitivity
- Missing/invalid arguments
- Duplicate detection

⚠️ **Not Tested (Manual QA):**
- TUI interactive components
- User input prompts (stdin mocking is complex)
- Visual rendering
- Keyboard navigation

### Test Utilities

**Helper Functions:**
```go
// mustCreateEntity creates an entity or fails the test.
func mustCreateEntity(t *testing.T, params) *Entity {
    t.Helper()
    entity, err := NewEntity(params)
    if err != nil {
        t.Fatalf("failed to create entity: %v", err)
    }
    return entity
}

// setupTestRepos creates a temporary storage repository.
func setupTestRepos(t *testing.T) *storage.Repositories {
    t.Helper()
    repos := storage.New(t.TempDir() + "/test.json")
    if err := repos.Load(); err != nil {
        t.Fatalf("Failed to load test repos: %v", err)
    }
    return repos
}
```

### Running Tests

```bash
# Run all tests
make test
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/domain
go test ./internal/storage
go test ./internal/cli/commands

# Specific test
go test -run TestSnippet_Create ./internal/domain

# Verbose output
go test -v ./...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Race detector
go test -race ./...
```

---

## Development Guidelines

### Code Organization

1. **Domain layer is pure** - No external dependencies except standard library
2. **Validate at boundaries** - Entities validate on construction/mutation
3. **Defensive copying** - Return copies of mutable data (slices)
4. **Consistent error handling** - Use sentinel errors with `errors.New()`
5. **Idiomatic Go** - Follow Go conventions and best practices
6. **Thread-safe by default** - Use mutexes even if not immediately needed

### Naming Conventions

**Files:**
```
entity.go
entity_test.go
```

**Constructors:**
```go
NewEntity(params) (*Entity, error)   // Exported
newEntity(params)                     // Unexported (internal)
```

**Getters/Setters:**
```go
Field()                    // Not GetField()
SetField(value) error      // With validation
SetID(id int)              // No validation (storage layer only)
```

### Comment Style

```go
// Package name provides description.
package name

// Entity is description of entity.
// Additional details if needed.
type Entity struct {
    // ...
}

// NewEntity creates a new entity with validation.
func NewEntity(params) (*Entity, error) {
    // ...
}
```

**Guidelines:**
- Package comment at top of file
- Type comments before declaration
- Function/method comments start with name
- Concise, focus on what not how
- No verbose "Parameters:", "Returns:" sections

### Error Handling

**Use Sentinel Errors:**
```go
var (
    ErrNotFound = errors.New("entity not found")
    ErrEmptyName = errors.New("name cannot be empty")
)
```

**Check with errors.Is():**
```go
if errors.Is(err, storage.ErrNotFound) {
    // Handle not found
}
```

**Provide Context:**
```go
return fmt.Errorf("failed to save category: %w", err)
```

### Git Workflow

**Branch Naming:**
```
feature/your-feature-name
refactor/layer-name
fix/bug-description
docs/documentation-updates
```

**Commit Message Convention:**
```
type(scope): subject

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `refactor`: Code refactoring
- `test`: Adding/updating tests
- `docs`: Documentation
- `style`: Formatting
- `perf`: Performance
- `chore`: Build/tools

**Scopes:**
- `domain`: Domain layer
- `storage`: Storage layer
- `cli`: CLI layer
- `tui`: TUI layer
- `config`: Configuration

**Examples:**
```
refactor(cli): improve command structure and validation
test(storage): add search functionality tests
feat(tui): add code editor component
docs: update development documentation
fix(domain): handle nil tags array in unmarshal
```

### Code Review Checklist

- [ ] All tests pass (`go test ./...`)
- [ ] Code formatted (`go fmt ./...`)
- [ ] No vet warnings (`go vet ./...`)
- [ ] New code has tests (where applicable)
- [ ] Comments follow idiomatic style
- [ ] Commit messages follow convention
- [ ] No breaking changes (or documented)
- [ ] Thread safety considered
- [ ] Error handling is consistent

---

## Build and Release

### Makefile Targets

**Development:**
```bash
make build        # Build for current platform
make run          # Build and run
make test         # Run all tests
make coverage     # Generate coverage report
make fmt          # Format code
make vet          # Run go vet
make lint         # Run golangci-lint (if installed)
```

**Release:**
```bash
make build-all    # Build for all platforms (Windows, Linux, macOS)
make clean        # Remove build artifacts
make install      # Install to system (requires admin)
make uninstall    # Remove from system
```

**Dependencies:**
```bash
make deps         # Download dependencies
make tidy         # Tidy go.mod and go.sum
```

### Build Output

**Current Platform:**
- Output: `build/snip.exe` (Windows) or `build/snip` (Unix)

**All Platforms:**
```
dist/
├── snip-linux-amd64
├── snip-linux-arm64
├── snip-darwin-amd64
├── snip-darwin-arm64
└── snip-windows-amd64.exe
```

### Release Process

**Version Tags:**
```bash
git tag -a v0.1.0 -m "Release version 0.1.0"
git push origin v0.1.0
```

**GitHub Actions:**
- `.github/workflows/cl.yml`: Continuous integration
- `.github/workflows/release.yml`: Automated releases

**Manual Release:**
1. Update version in code
2. Run tests: `make test`
3. Build all platforms: `make build-all`
4. Create Git tag
5. Push tag to GitHub
6. Upload binaries to GitHub Releases

### Installation

**Pre-built Binary:**
```bash
# Download from releases page
curl -L https://github.com/7-Dany/snip/releases/download/v0.1.1/snip-linux-amd64 -o snip
chmod +x snip
sudo mv snip /usr/local/bin/
```

**From Source:**
```bash
git clone https://github.com/7-Dany/snip.git
cd snip
make build
sudo make install
```

### System Requirements

- **OS**: Windows, Linux, macOS
- **Go**: 1.24+ (for building from source)
- **Terminal**: Any modern terminal with Unicode support
- **Storage**: Minimal (~5MB for binary, varies for data)

---

## Common Patterns

### Creating Domain Entities

```go
// Category
category, err := domain.NewCategory("Algorithms")
if err != nil {
    // Handle validation error
}

// Tag
tag, err := domain.NewTag("golang")
if err != nil {
    // Handle validation error
}

// Snippet
snippet, err := domain.NewSnippet("QuickSort", "go", code)
if err != nil {
    // Handle validation error
}
snippet.SetDescription("Fast sorting algorithm")
snippet.SetCategory(categoryID)
snippet.AddTag(tagID)
```

### Using Repositories

```go
// Initialize
repos := storage.New("~/.snip/snippets.json")
repos.Load()
defer repos.Save()

// Create
category, _ := domain.NewCategory("Web")
repos.Categories.Create(category)

// List
categories, _ := repos.Categories.List()
for _, cat := range categories {
    fmt.Println(cat.Name())
}

// Find
cat, err := repos.Categories.FindByName("Web")
if errors.Is(err, storage.ErrNotFound) {
    // Not found
}

// Update
cat.SetName("Web Development")
repos.Categories.Update(cat)

// Delete
repos.Categories.Delete(catID)
```

### Handling Errors

```go
// Domain errors
if err != nil {
    if errors.Is(err, domain.ErrEmptyName) {
        // Handle empty name
    }
}

// Storage errors
if err != nil {
    if errors.Is(err, storage.ErrNotFound) {
        // Handle not found
    }
    if errors.Is(err, storage.ErrDuplicateName) {
        // Handle duplicate
    }
}
```

### Working with Bubble Tea

```go
// Simple model
type myModel struct {
    input textinput.Model
}

func (m myModel) Init() tea.Cmd {
    return textinput.Blink
}

func (m myModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            return m, tea.Quit
        case "esc", "ctrl+c":
            return m, tea.Quit
        }
    }
    var cmd tea.Cmd
    m.input, cmd = m.input.Update(msg)
    return m, cmd
}

func (m myModel) View() string {
    return m.input.View()
}

// Run
p := tea.NewProgram(initialModel())
finalModel, err := p.Run()
```

---

## Troubleshooting

### Common Issues

**Config/Storage Errors:**
```
Error: failed to load config
Solution: Check ~/.snip/ directory permissions
```

**Data File Corruption:**
```
Error: failed to unmarshal config
Solution: Check ~/.snip/snippets.json for valid JSON
```

**Build Errors:**
```
Error: go.mod requires Go 1.24+
Solution: Upgrade Go to 1.24 or higher
```

**Test Failures:**
```
Error: TestStore_SaveAndLoad failed
Solution: Check if test has write permissions in temp dir
```

### Debugging Tips

**Enable Verbose Logging:**
```go
// In development, add debug prints
fmt.Printf("DEBUG: loaded %d snippets\n", len(snippets))
```

**Check Data File:**
```bash
# View current data
cat ~/.snip/snippets.json | jq
```

**Reset Configuration:**
```bash
rm -rf ~/.snip
# Run snip to recreate
```

---

## Future Enhancements

### Planned Features

- [ ] Export/import snippets (JSON, Markdown)
- [ ] Syntax highlighting in code viewer
- [ ] Snippet templates
- [ ] Cloud sync (optional)
- [ ] Snippet sharing (GitHub Gist integration)
- [ ] Advanced search (regex, filters)
- [ ] Snippet versioning
- [ ] Custom themes
- [ ] Plugin system

### Performance Optimizations

- [ ] Lazy loading for large snippet collections
- [ ] Indexed search (full-text index)
- [ ] Caching for frequently accessed snippets
- [ ] Pagination for large lists

### Developer Experience

- [ ] Better error messages
- [ ] Progress indicators for long operations
- [ ] Undo/redo support
- [ ] Keyboard shortcuts customization

---

## Resources

### External Documentation

- **Bubble Tea**: https://github.com/charmbracelet/bubbletea
- **Bubbles**: https://github.com/charmbracelet/bubbles
- **Lipgloss**: https://github.com/charmbracelet/lipgloss
- **Go Pretty**: https://github.com/jedib0t/go-pretty

### Project Links

- **Repository**: https://github.com/7-Dany/snip
- **Issues**: https://github.com/7-Dany/snip/issues
- **Releases**: https://github.com/7-Dany/snip/releases

### Related Projects

- **Pet**: Another snippet manager in Go
- **Nap**: Code snippet organizer
- **SnipKit**: Cross-platform snippet manager

---

## Changelog

### v0.1.1 (2025-02-05)

- refactor(cli): Replace table view with searchable table view
- refactor(cli): Delete unnecessary comments
- docs: Update development documentation

### v0.1.0 (2024-12-30)

- Initial release
- Complete TUI implementation
- CLI commands for snippets, categories, tags
- JSON-based storage
- Full-text search
- Comprehensive test coverage

---

## Contact & Support

**Author**: 7-Dany
**Email**: ali.ali.ameen245@gmail.com
**GitHub**: [@7-Dany](https://github.com/7-Dany)

**Contributing**:
- Fork the repository
- Create a feature branch
- Make your changes
- Add tests
- Submit a pull request

**Bug Reports**:
- Use GitHub Issues
- Include OS, Go version, and steps to reproduce
- Provide error messages and logs

---

## License

MIT License - see LICENSE file for details.

Copyright (c) 2024 7-Dany

---

*Last Updated: 2025-02-05*
*Document Version: 2.0*
