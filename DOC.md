# SNIP - Development Documentation

> Technical documentation for SNIP development.
> Complete context about architecture, design decisions, and implementation details.

## Table of Contents

1. [Tech Stack](#tech-stack)
2. [Architecture](#architecture)
3. [Domain Layer](#domain-layer)
4. [Storage Layer](#storage-layer)
5. [CLI Layer](#cli-layer)
6. [Testing Strategy](#testing-strategy)
7. [Development Guidelines](#development-guidelines)

---

## Tech Stack

- **Language**: Go 1.25+
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Tables**: [go-pretty](https://github.com/jedib0t/go-pretty)
- **Colors**: [fatih/color](https://github.com/fatih/color)
- **Storage**: JSON files (local filesystem)
- **Testing**: Go standard testing package

---

## Architecture

```
snip/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── cli/
│   │   ├── commands/        # CLI command handlers (COMPLETE)
│   │   ├── components/      # Reusable Bubble Tea UI components
│   │   └── tui/             # Terminal UI implementation
│   ├── domain/              # Business logic and entities (COMPLETE)
│   └── storage/             # Data persistence layer (COMPLETE)
└── main.go
```

### Layer Responsibilities

**Domain Layer** (`internal/domain/`)
- Pure business logic, entities, and domain rules
- No external dependencies (only standard library)
- Validation at entity boundaries
- Immutable where possible

**Storage Layer** (`internal/storage/`)
- Data persistence using JSON files
- Repository pattern with shared internal store
- Thread-safe operations with RWMutex
- Search and filtering capabilities

**CLI Layer** (`internal/cli/commands/`)
- Command-line interface handlers
- Input validation and user interaction
- Interactive prompts using Bubble Tea
- Formatted output and error messages

**TUI Layer** (`internal/cli/tui/` and `internal/cli/components/`)
- Interactive terminal user interface
- Tabbed navigation
- Code editor with syntax highlighting
- Reusable UI components

---

## Domain Layer

Location: `internal/domain/`

### Design Principles

1. **Encapsulation**: All fields are unexported
2. **Immutability**: Getters return values, not pointers (slices are copied)
3. **Validation**: All constructors and setters validate input
4. **Timestamps**: All entities track creation and modification times
5. **JSON Serialization**: Custom MarshalJSON/UnmarshalJSON for unexported fields

### Entities

#### Category

Organizes snippets into logical groups.

**Fields:**
- `id` (int) - Unique identifier
- `name` (string) - Category name (required, non-empty)
- `createdAt`, `updatedAt` (time.Time)

**Key Methods:**
- `NewCategory(name string) (*Category, error)` - Constructor with validation
- `ID()`, `Name()`, `CreatedAt()`, `UpdatedAt()` - Getters
- `SetName(name string) error` - Updates name with validation
- `SetID(id int)` - Sets ID (storage layer only)

**Validation:**
- Name cannot be empty (returns `ErrEmptyName`)

#### Tag

Labels for categorizing snippets.

**Fields:**
- `id` (int) - Unique identifier
- `name` (string) - Tag name (required, non-empty)
- `createdAt`, `updatedAt` (time.Time)

**Key Methods:**
- `NewTag(name string) (*Tag, error)` - Constructor with validation
- `ID()`, `Name()`, `CreatedAt()`, `UpdatedAt()` - Getters
- `SetName(name string) error` - Updates name with validation
- `SetID(id int)` - Sets ID (storage layer only)

**Validation:**
- Name cannot be empty (returns `ErrEmptyName`)

#### Snippet

Code snippet with metadata.

**Fields:**
- `id` (int) - Unique identifier
- `title` (string) - Snippet title (required, non-empty)
- `language` (string) - Programming language (required, non-empty)
- `code` (string) - Code content (required, non-empty)
- `description` (string) - Optional description
- `categoryID` (int) - Category ID (0 if uncategorized)
- `tags` ([]int) - Tag IDs (never nil, always []int{})
- `createdAt`, `updatedAt` (time.Time)

**Key Methods:**
- `NewSnippet(title, language, code string) (*Snippet, error)` - Constructor
- Getters for all fields
- `Tags() []int` - Returns copy of tag IDs
- `SetTitle()`, `SetLanguage()`, `SetCode()` - Update with validation
- `SetDescription()`, `SetCategory()` - Update without validation
- `AddTag(tagID int)`, `RemoveTag(tagID int)`, `HasTag(tagID int) bool`
- `SetID(id int)` - Sets ID (storage layer only)

**Validation:**
- Title, language, and code cannot be empty
- Tags array is never nil (normalized to []int{} on unmarshal)
- AddTag prevents duplicates
- RemoveTag is no-op for non-existent tags

### Domain Errors

```go
var (
    ErrEmptyName     = errors.New("name cannot be empty")
    ErrEmptyTitle    = errors.New("title cannot be empty")
    ErrEmptyLanguage = errors.New("language cannot be empty")
    ErrEmptyCode     = errors.New("code cannot be empty")
)
```

### Code Patterns

**String() Method:**
```go
func (e *Entity) String() string {
    return fmt.Sprintf("Entity{id=%d, name=%q}", e.id, e.name)
}
```

**Equal() Method:**
```go
func (e *Entity) Equal(other *Entity) bool {
    if other == nil {
        return false
    }
    return e.id == other.id &&
        e.name == other.name &&
        e.createdAt.Equal(other.createdAt) &&
        e.updatedAt.Equal(other.updatedAt)
}
```

**JSON Marshaling:**
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

**File Structure:**
```
storage/
├── repositories.go              # Public API
├── repositories_test.go         # Integration tests
├── internal_store.go            # Internal data structure
├── internal_store_test.go       # Store unit tests
├── search.go                    # Search index
├── search_test.go               # Search tests
├── snippet_repository.go        # Snippet operations
├── snippet_repository_test.go
├── category_repository.go       # Category operations
├── category_repository_test.go
├── tag_repository.go            # Tag operations
└── tag_repository_test.go
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

**Methods:**
- `New(filepath string) *Repositories` - Creates repositories with shared store
- `Save() error` - Persists all data atomically to JSON file
- `Load() error` - Reads all data from JSON file

**Usage:**
```go
repos := storage.New("~/.snip/data.json")
if err := repos.Load(); err != nil {
    log.Fatal(err)
}
defer repos.Save()

snippet, _ := domain.NewSnippet("quicksort", "go", "func quicksort() {}")
repos.Snippets.Create(snippet)
```

### Internal Store

**Key Features:**
- Atomic saves using temporary file + rename
- Auto-incrementing IDs
- Thread-safe operations with sync.RWMutex
- Nil slice normalization on load

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

### Search Index

**Methods:**
- `search(query string) []*Snippet` - Full-text search (case-insensitive)
- `findByLanguage(language string) []*Snippet`
- `findByCategory(categoryID int) []*Snippet`
- `findByTag(tagID int) []*Snippet`

**Behavior:**
- Searches across title, language, code, and description
- Returns empty slice for no matches
- Returns nil only for empty query

### Storage Errors

```go
var (
    ErrNotFound      = errors.New("entity not found")
    ErrDuplicateName = errors.New("entity with this name already exists")
)
```

### Repository Pattern

**Common Operations:**
- `List() ([]*Entity, error)` - Returns defensive copy
- `FindByID(id int) (*Entity, error)` - Returns entity or ErrNotFound
- `FindByName(name string) (*Entity, error)` - Returns entity or ErrNotFound
- `Create(entity *Entity) error` - Auto-assigns ID
- `Update(entity *Entity) error` - Validates entity exists
- `Delete(id int) error` - Idempotent (only ErrNotFound if not exists)

**Thread Safety:**
- All operations use RWMutex
- Read operations: `RLock()`/`RUnlock()`
- Write operations: `Lock()`/`Unlock()`

### Design Decisions

1. **Separate internal store from repositories** - Clear separation of concerns
2. **Shared store across repositories** - Ensures data consistency
3. **RWMutex for concurrency** - Prepares for future concurrent usage
4. **Defensive copying in List()** - Prevents external state modification
5. **Search as separate component** - Single responsibility, easy to enhance

---

## CLI Layer

Location: `internal/cli/commands/`

### Architecture

**File Structure:**
```
commands/
├── commands.go              # Main CLI coordinator
├── commands_test.go
├── output.go                # Display utilities
├── output_test.go
├── help.go                  # Help command handler
├── help_test.go
├── category.go              # Category command handler
├── category_test.go
├── tag.go                   # Tag command handler
├── tag_test.go
├── snippet.go               # Snippet command handler
├── snippet_test.go
└── testing_helpers.go       # Shared test utilities
```

### Design Principles

1. **Consistent Structure**: All command handlers follow same pattern
2. **Input Validation**: Validate early, provide clear error messages
3. **User Experience**: Interactive prompts with Bubble Tea
4. **Error Handling**: Use `errors.Is()` for sentinel errors
5. **Testing**: Comprehensive test coverage for all commands

### Command Handler Pattern

**Structure:**
```go
type XCommand struct {
    repos *storage.Repositories
}

func NewXCommand(repos *storage.Repositories) *XCommand {
    return &XCommand{repos: repos}
}

func (xc *XCommand) manage(args []string) {
    if len(args) == 0 {
        PrintError("No subcommand provided. Use 'snip help x' for available commands")
        return
    }

    subcommand := strings.ToLower(args[0])
    subcommandArgs := args[1:]

    switch subcommand {
    case "list":
        xc.list()
    case "create":
        xc.create(subcommandArgs)
    case "delete":
        xc.delete(subcommandArgs)
    default:
        PrintError(fmt.Sprintf("Unknown command '%s'", args[0]))
    }
}
```

### CLI Coordinator

**CLI Struct:**
```go
type CLI struct {
    snippet  *SnippetCommand
    category *CategoryCommand
    tag      *TagCommand
    help     *HelpCommand
}
```

**Run Method:**
- Routes commands to appropriate handler
- Handles backward compatibility (default to snippet)
- Provides clear error messages

### Command Operations

**List:**
- Displays entities in formatted table
- Shows informative message when empty
- Uses go-pretty for table rendering
- Supports filtering (snippets only: --category, --tag, --language)

**Create:**
- Supports both argument and interactive modes
- Validates input (trim whitespace, check duplicates)
- Interactive prompts use Bubble Tea
- Shows success message with ID
- Snippets use multi-field form with code editor

**Update:**
- Interactive form pre-populated with existing data
- Validates all inputs
- Updates all fields including relationships (snippet tags)
- Shows success message with ID

**Delete:**
- Validates ID is provided and numeric
- Confirms before deletion
- Shows clear error if entity not found
- Idempotent (no error on re-delete)

**Show (Snippets only):**
- Displays full snippet details including code
- Resolves category and tag names
- Shows creation and update timestamps

**Search (Snippets only):**
- Full-text search across title, language, code, description
- Case-insensitive matching
- Shows results in table format

### Output Functions

**Available Functions:**
- `PrintLogo()` - Application branding
- `PrintSuccess(msg string)` - Green checkmark message
- `PrintError(msg string)` - Red X message
- `PrintInfo(msg string)` - Cyan informational message
- `PrintCommand(command string)` - Command execution indicator
- `PrintLoading(msg string)` - Loading animation
- `ClearScreen()` - Terminal screen clear

**Color Scheme:**
- Success: Green bold
- Error: Red bold
- Info: Cyan
- Command: Gray
- Logo: Magenta/Yellow

### Interactive Prompts

**Bubble Tea Models:**
```go
type xInputModel struct {
    textInput textinput.Model
    cancelled bool
}

func (m xInputModel) Init() tea.Cmd { return textinput.Blink }
func (m xInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m xInputModel) View() string
```

**Key Bindings:**
- Enter: Confirm input
- Esc/Ctrl+C: Cancel operation
- Tab/Shift+Tab: Navigate fields (in forms)

---

## Testing Strategy

### Test Organization

**Pattern**: Use `t.Run()` for subtests (NOT table-driven tests)

**Structure:**
```go
func TestEntity_Method(t *testing.T) {
    t.Run("description of test case", func(t *testing.T) {
        // Arrange
        entity := mustCreateEntity(t, "params")
        
        // Act
        result := entity.Method()
        
        // Assert
        if result != expected {
            t.Errorf("expected %v, got %v", expected, result)
        }
    })
}
```

### Test Coverage Areas

**Domain Layer:**
- Constructor validation
- Setter validation
- Getter correctness
- JSON marshaling/unmarshaling
- Equal() method
- String() method
- Edge cases (nil, empty, zero values)

**Storage Layer:**
- CRUD operations
- Search functionality
- Error cases (not found, duplicates)
- Thread safety
- Defensive copying
- Save/load round trips
- Nil slice normalization

**CLI Layer:**
- Command routing
- Input validation
- Error handling
- Case insensitivity
- Missing arguments
- Invalid argument types
- Empty/whitespace input
- Duplicate detection

### Helper Functions

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

// setupTestRepos creates a temporary storage repository for testing.
func setupTestRepos(t *testing.T) *storage.Repositories {
    t.Helper()
    repos := storage.New(t.TempDir() + "/test.json")
    if err := repos.Load(); err != nil {
        t.Fatalf("Failed to load test repos: %v", err)
    }
    return repos
}
```

### Test Execution

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package
go test ./internal/storage

# Run specific test
go test -run TestStore_SaveAndLoad ./internal/storage

# Verbose output
go test -v ./...

# With race detector
go test -race ./...
```

---

## Development Guidelines

### Code Organization

1. **Domain layer is pure** - No external dependencies
2. **Validate at boundaries** - Entities validate on construction/mutation
3. **Defensive copying** - Return copies of mutable data
4. **Consistent error handling** - Use sentinel errors with `errors.New()`
5. **Idiomatic Go** - Follow Go conventions and best practices
6. **Thread-safe by default** - Use mutexes even if not immediately needed

### Naming Conventions

**Files:**
- `entity.go` and `entity_test.go`
- Lowercase with underscores

**Constructors:**
- `NewEntity(params) (*Entity, error)` - Exported
- `newEntity(params)` - Unexported (internal)

**Getters/Setters:**
- `Field()` - Not `GetField()`
- `SetField(value) error` - With validation
- `SetID(id int)` - No validation (storage layer only)

**Private Methods:**
- Unexported (lowercase first letter)
- Used only within package

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

**Key Points:**
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

**Examples:**
```
refactor(cli): improve command structure and validation
test(cli): add comprehensive tests for all commands
feat(storage): add search index functionality
docs: update development documentation
```

### Code Review Checklist

- [ ] All tests pass (`go test ./...`)
- [ ] Code formatted (`go fmt ./...`)
- [ ] No vet warnings (`go vet ./...`)
- [ ] New code has tests
- [ ] Comments follow idiomatic style
- [ ] Commit messages follow convention
- [ ] No breaking changes (or documented)
- [ ] Thread safety considered

---

## Project Status

### Completed Layers

✅ **Domain Layer** (`internal/domain/`)
- Category, Tag, Snippet entities
- Full validation and error handling
- JSON serialization
- Comprehensive tests (>100.0% coverage)

✅ **Storage Layer** (`internal/storage/`)
- Repository pattern implementation
- Thread-safe operations
- Search functionality
- Comprehensive tests (>99.5% coverage)

✅ **CLI Commands** (`internal/cli/commands/`)
- All command handlers (snippet, category, tag, help)
- Input validation and interactive prompts
- Error handling and user feedback
- Comprehensive tests (>85% coverage)
