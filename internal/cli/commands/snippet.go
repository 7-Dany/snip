package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/7-Dany/snip/internal/domain"
	"github.com/7-Dany/snip/internal/storage"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedib0t/go-pretty/v6/table"
)

type SnippetCommand struct {
	repos *storage.Repositories
}

func NewSnippetCommand(repos *storage.Repositories) *SnippetCommand {
	return &SnippetCommand{repos: repos}
}

func (sc *SnippetCommand) manage(args []string) {
	if len(args) == 0 {
		PrintError("No subcommand provided. Use 'snip help snippet' for available commands")
		return
	}

	switch strings.ToLower(args[0]) {
	case "list":
		sc.list(args[1:])
	case "show":
		sc.show(args[1:])
	case "create":
		sc.create()
	case "update":
		sc.update(args[1:])
	case "delete":
		sc.delete(args[1:])
	case "search":
		sc.search(args[1:])
	default:
		PrintError(fmt.Sprintf("Unknown command '%s'. Use 'snip help snippet' for available commands", args[0]))
	}
}

// list displays all snippets with optional filters
func (sc *SnippetCommand) list(args []string) {
	var categoryID, tagID int
	var language string

	// Parse filter flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--category":
			if i+1 >= len(args) {
				PrintError("Missing value for --category flag")
				return
			}
			var err error
			categoryID, err = strconv.Atoi(args[i+1])
			if err != nil {
				PrintError(fmt.Sprintf("Validation failed. '--category' must be a number: %v", err))
				return
			}
			i++
		case "--tag":
			if i+1 >= len(args) {
				PrintError("Missing value for --tag flag")
				return
			}
			var err error
			tagID, err = strconv.Atoi(args[i+1])
			if err != nil {
				PrintError(fmt.Sprintf("Validation failed. '--tag' must be a number: %v", err))
				return
			}
			i++
		case "--language":
			if i+1 >= len(args) {
				PrintError("Missing value for --language flag")
				return
			}
			language = args[i+1]
			i++
		}
	}

	var snippets []*domain.Snippet
	var err error

	// Apply filters or list all
	if categoryID > 0 {
		snippets, err = sc.repos.Snippets.FindByCategory(categoryID)
	} else if tagID > 0 {
		snippets, err = sc.repos.Snippets.FindByTag(tagID)
	} else if language != "" {
		snippets, err = sc.repos.Snippets.FindByLanguage(language)
	} else {
		snippets, err = sc.repos.Snippets.List()
	}

	if err != nil {
		PrintError(fmt.Sprintf("failed to list snippets: %v", err))
		return
	}

	if len(snippets) == 0 {
		PrintInfo("no snippets found, create one with 'snip snippet create'")
		return
	}

	sc.displaySnippets(snippets)
}

func (sc *SnippetCommand) show(args []string) {
	if len(args) == 0 {
		PrintError("Missing required argument 'id'. Use 'snip snippet show <id>'")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		PrintError(fmt.Sprintf("Validation failed. 'id' must be a number: %v", err))
		return
	}

	snippet, err := sc.repos.Snippets.FindByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		PrintError(fmt.Sprintf("Snippet with ID %d not found", id))
		return
	}

	if err != nil {
		PrintError(fmt.Sprintf("Failed to find snippet: %v", err))
		return
	}

	categoryName := "N/A"
	if snippet.CategoryID() > 0 {
		category, err := sc.repos.Categories.FindByID(snippet.CategoryID())
		if err == nil && category != nil {
			categoryName = category.Name()
		}
	}

	tagNames := sc.getTagNames(snippet.Tags())

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📝 %s\n", snippet.Title())
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Language:    %s\n", snippet.Language())
	fmt.Printf("Category:    %s\n", categoryName)
	fmt.Printf("Tags:        %s\n", strings.Join(tagNames, ", "))
	fmt.Printf("Description: %s\n", snippet.Description())
	fmt.Println("\n--- Code ---")
	fmt.Println(snippet.Code())
	fmt.Println("--- End ---")
	fmt.Printf("Created: %s\n", snippet.CreatedAt().Format("2006-01-02 15:04"))
	fmt.Printf("Updated: %s\n", snippet.UpdatedAt().Format("2006-01-02 15:04"))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func (sc *SnippetCommand) create() {
	formData := sc.promptForSnippet(nil)
	if formData == nil {
		PrintInfo("Create cancelled")
		return
	}

	snippet, err := domain.NewSnippet(formData.title, formData.language, formData.code)
	if err != nil {
		PrintError(fmt.Sprintf("failed to create snippet: %v", err))
		return
	}

	if formData.categoryID > 0 {
		snippet.SetCategory(formData.categoryID)
	}
	for _, tagID := range formData.tags {
		snippet.AddTag(tagID)
	}
	if formData.description != "" {
		snippet.SetDescription(formData.description)
	}

	err = sc.repos.Snippets.Create(snippet)
	if err != nil {
		PrintError(fmt.Sprintf("failed to save snippet: %v", err))
		return
	}

	PrintSuccess(fmt.Sprintf("Successfully created snippet '%s' (ID: %d)", formData.title, snippet.ID()))
}

func (sc *SnippetCommand) update(args []string) {
	if len(args) == 0 {
		PrintError("Missing required argument 'id'. Use 'snip snippet update <id>'")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		PrintError(fmt.Sprintf("Validation failed. 'id' must be a number: %v", err))
		return
	}

	snippet, err := sc.repos.Snippets.FindByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		PrintError(fmt.Sprintf("Snippet with ID %d not found", id))
		return
	}

	if err != nil {
		PrintError(fmt.Sprintf("Failed to find snippet: %v", err))
		return
	}

	formData := sc.promptForSnippet(snippet)
	if formData == nil {
		PrintInfo("Update cancelled")
		return
	}

	if err := snippet.SetTitle(formData.title); err != nil {
		PrintError(fmt.Sprintf("failed to set title: %v", err))
		return
	}
	if err := snippet.SetLanguage(formData.language); err != nil {
		PrintError(fmt.Sprintf("failed to set language: %v", err))
		return
	}
	if err := snippet.SetCode(formData.code); err != nil {
		PrintError(fmt.Sprintf("failed to set code: %v", err))
		return
	}

	snippet.SetCategory(formData.categoryID)
	snippet.SetDescription(formData.description)

	// Update tags
	for _, existingTag := range snippet.Tags() {
		snippet.RemoveTag(existingTag)
	}
	for _, tagID := range formData.tags {
		snippet.AddTag(tagID)
	}

	err = sc.repos.Snippets.Update(snippet)
	if err != nil {
		PrintError(fmt.Sprintf("failed to update snippet: %v", err))
		return
	}

	PrintSuccess(fmt.Sprintf("Successfully updated snippet '%s'", formData.title))
}

func (sc *SnippetCommand) delete(args []string) {
	if len(args) == 0 {
		PrintError("Missing required argument 'id'. Use 'snip snippet delete <id>'")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		PrintError(fmt.Sprintf("Validation failed. 'id' must be a number: %v", err))
		return
	}

	snippet, err := sc.repos.Snippets.FindByID(id)
	if errors.Is(err, domain.ErrNotFound) {
		PrintError(fmt.Sprintf("Snippet with ID %d not found", id))
		return
	}

	if err != nil {
		PrintError(fmt.Sprintf("Failed to find snippet: %v", err))
		return
	}

	fmt.Printf("Are you sure you want to delete snippet '%s'? (y/n): ", snippet.Title())
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" {
		PrintInfo("Delete cancelled")
		return
	}

	err = sc.repos.Snippets.Delete(id)
	if err != nil {
		PrintError(fmt.Sprintf("Delete failed, couldn't delete snippet: %v", err))
		return
	}

	PrintSuccess(fmt.Sprintf("Deleted snippet with id: '%v'", id))
}

func (sc *SnippetCommand) search(args []string) {
	if len(args) == 0 {
		PrintError("Missing required argument 'keyword'. Use 'snip snippet search <keyword>'")
		return
	}

	keyword := args[0]
	snippets, err := sc.repos.Snippets.Search(keyword)
	if err != nil {
		PrintError(fmt.Sprintf("failed to search snippets: %v", err))
		return
	}

	if len(snippets) == 0 {
		PrintInfo(fmt.Sprintf("no snippets found matching '%s'", keyword))
		return
	}

	sc.displaySnippets(snippets)
}

func (sc *SnippetCommand) getTagNames(tagIDs []int) []string {
	if len(tagIDs) == 0 {
		return []string{"N/A"}
	}

	names := make([]string, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		tag, err := sc.repos.Tags.FindByID(tagID)
		if err == nil && tag != nil {
			names = append(names, tag.Name())
		}
	}

	if len(names) == 0 {
		return []string{"N/A"}
	}

	return names
}

func (sc *SnippetCommand) displaySnippets(snippets []*domain.Snippet) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Title", "Language", "Category", "Tags", "CreatedAt"})

	for _, snippet := range snippets {
		categoryName := "N/A"
		if snippet.CategoryID() > 0 {
			category, err := sc.repos.Categories.FindByID(snippet.CategoryID())
			if err == nil && category != nil {
				categoryName = category.Name()
			}
		}

		tagNames := sc.getTagNames(snippet.Tags())

		t.AppendRow(table.Row{
			snippet.ID(),
			snippet.Title(),
			snippet.Language(),
			categoryName,
			strings.Join(tagNames, ", "),
			snippet.CreatedAt().Format("2006-01-02 15:04"),
		})
	}

	t.SetStyle(table.StyleColoredBright)
	t.Render()
}

// Form models and helpers
type snippetFormData struct {
	title       string
	language    string
	categoryID  int
	tags        []int
	description string
	code        string
}

type snippetFormModel struct {
	inputs    []textinput.Model
	codeArea  textarea.Model
	focusIdx  int
	cancelled bool
}

func newSnippetFormModel(existing *domain.Snippet) snippetFormModel {
	inputs := make([]textinput.Model, 5)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "e.g., Binary Search Implementation"
	inputs[0].CharLimit = 100
	inputs[0].Width = 50
	inputs[0].Focus()

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "e.g., go, python, javascript"
	inputs[1].CharLimit = 50
	inputs[1].Width = 50

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "e.g., 1 (or leave empty)"
	inputs[2].CharLimit = 10
	inputs[2].Width = 50

	inputs[3] = textinput.New()
	inputs[3].Placeholder = "e.g., 1,2,3 (comma-separated IDs)"
	inputs[3].CharLimit = 100
	inputs[3].Width = 50

	inputs[4] = textinput.New()
	inputs[4].Placeholder = "Brief description of the snippet"
	inputs[4].CharLimit = 200
	inputs[4].Width = 50

	codeArea := textarea.New()
	codeArea.Placeholder = "Paste or type your code here..."
	codeArea.CharLimit = 10000
	codeArea.SetWidth(80)
	codeArea.SetHeight(10)

	if existing != nil {
		inputs[0].SetValue(existing.Title())
		inputs[1].SetValue(existing.Language())
		if existing.CategoryID() > 0 {
			inputs[2].SetValue(strconv.Itoa(existing.CategoryID()))
		}
		if len(existing.Tags()) > 0 {
			tagStrs := make([]string, len(existing.Tags()))
			for i, t := range existing.Tags() {
				tagStrs[i] = strconv.Itoa(t)
			}
			inputs[3].SetValue(strings.Join(tagStrs, ","))
		}
		inputs[4].SetValue(existing.Description())
		codeArea.SetValue(existing.Code())
	}

	return snippetFormModel{
		inputs:   inputs,
		codeArea: codeArea,
		focusIdx: 0,
	}
}

func (m snippetFormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m snippetFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.cancelled = true
			return m, tea.Quit
		case tea.KeyEnter:
			if m.focusIdx == 5 {
				m.codeArea, cmd = m.codeArea.Update(msg)
				return m, cmd
			}
			return m, tea.Quit
		case tea.KeyTab, tea.KeyShiftTab:
			if msg.Type == tea.KeyTab {
				m.focusIdx++
			} else {
				m.focusIdx--
			}

			if m.focusIdx > 5 {
				m.focusIdx = 0
			} else if m.focusIdx < 0 {
				m.focusIdx = 5
			}

			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIdx {
					m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}

			if m.focusIdx == 5 {
				m.codeArea.Focus()
			} else {
				m.codeArea.Blur()
			}

			return m, nil
		}
	}

	if m.focusIdx < 5 {
		m.inputs[m.focusIdx], cmd = m.inputs[m.focusIdx].Update(msg)
	} else {
		m.codeArea, cmd = m.codeArea.Update(msg)
	}

	return m, cmd
}

func (m snippetFormModel) View() string {
	return fmt.Sprintf(
		"\n📝 Snippet Form\n\n"+
			"Title:       %s\n"+
			"Language:    %s\n"+
			"Category ID: %s\n"+
			"Tag IDs:     %s\n"+
			"Description: %s\n\n"+
			"Code:\n%s\n\n"+
			"(Tab/Shift+Tab to navigate, Enter to submit, Esc to cancel)\n",
		m.inputs[0].View(),
		m.inputs[1].View(),
		m.inputs[2].View(),
		m.inputs[3].View(),
		m.inputs[4].View(),
		m.codeArea.View(),
	)
}

func (sc *SnippetCommand) promptForSnippet(existing *domain.Snippet) *snippetFormData {
	p := tea.NewProgram(newSnippetFormModel(existing))
	finalModel, err := p.Run()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to run form: %v", err))
		return nil
	}

	m := finalModel.(snippetFormModel)
	if m.cancelled {
		return nil
	}

	title := strings.TrimSpace(m.inputs[0].Value())
	language := strings.TrimSpace(m.inputs[1].Value())
	description := strings.TrimSpace(m.inputs[4].Value())
	code := strings.TrimSpace(m.codeArea.Value())

	if title == "" || language == "" || code == "" {
		PrintError("Title, language, and code are required fields")
		return nil
	}

	categoryID := 0
	if catStr := strings.TrimSpace(m.inputs[2].Value()); catStr != "" {
		categoryID, err = strconv.Atoi(catStr)
		if err != nil {
			PrintError("Category ID must be a number")
			return nil
		}
	}

	var tags []int
	if tagStr := strings.TrimSpace(m.inputs[3].Value()); tagStr != "" {
		tagParts := strings.Split(tagStr, ",")
		tags = make([]int, 0, len(tagParts))
		for _, part := range tagParts {
			tagID, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				PrintError("Tag IDs must be comma-separated numbers")
				return nil
			}
			tags = append(tags, tagID)
		}
	}

	return &snippetFormData{
		title:       title,
		language:    language,
		categoryID:  categoryID,
		tags:        tags,
		description: description,
		code:        code,
	}
}
