package storage

import (
	"testing"
)

func TestSnippetRepository_List(t *testing.T) {
	t.Run("returns empty slice when no snippets", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippets, err := repo.List()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if snippets == nil {
			t.Fatal("expected empty slice, got nil")
		}

		if len(snippets) != 0 {
			t.Errorf("expected 0 snippets, got %d", len(snippets))
		}
	})

	t.Run("returns all snippets", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet1 := mustCreateSnippet(t, "quicksort", "go", "code1")
		snippet2 := mustCreateSnippet(t, "mergesort", "go", "code2")
		snippet3 := mustCreateSnippet(t, "bubblesort", "go", "code3")

		repo.Create(snippet1)
		repo.Create(snippet2)
		repo.Create(snippet3)

		snippets, err := repo.List()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 3 {
			t.Errorf("expected 3 snippets, got %d", len(snippets))
		}
	})

	t.Run("returns defensive copy", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		original := mustCreateSnippet(t, "test", "go", "code")
		repo.Create(original)

		snippets, _ := repo.List()
		snippets[0] = nil

		// Original should be unchanged
		stored, _ := repo.FindByID(original.ID())
		if stored == nil {
			t.Error("modifying returned slice affected internal store")
		}
	})
}

func TestSnippetRepository_FindByID(t *testing.T) {
	t.Run("finds existing snippet", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet := mustCreateSnippet(t, "quicksort", "go", "func quicksort() {}")
		repo.Create(snippet)

		found, err := repo.FindByID(snippet.ID())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if found.ID() != snippet.ID() {
			t.Errorf("expected ID %d, got %d", snippet.ID(), found.ID())
		}

		if found.Title() != "quicksort" {
			t.Errorf("expected title %q, got %q", "quicksort", found.Title())
		}
	})

	t.Run("returns ErrNotFound for nonexistent ID", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		_, err := repo.FindByID(999)

		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestSnippetRepository_FindByCategory(t *testing.T) {
	t.Run("finds snippets in category", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet1 := mustCreateSnippet(t, "test1", "go", "code1")
		snippet1.SetCategory(5)
		snippet2 := mustCreateSnippet(t, "test2", "go", "code2")
		snippet2.SetCategory(10)
		snippet3 := mustCreateSnippet(t, "test3", "go", "code3")
		snippet3.SetCategory(5)

		repo.Create(snippet1)
		repo.Create(snippet2)
		repo.Create(snippet3)

		snippets, err := repo.FindByCategory(5)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 2 {
			t.Errorf("expected 2 snippets, got %d", len(snippets))
		}
	})

	t.Run("returns empty slice for category with no snippets", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippets, err := repo.FindByCategory(999)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 0 {
			t.Errorf("expected 0 snippets, got %d", len(snippets))
		}
	})

	t.Run("finds uncategorized snippets", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet := mustCreateSnippet(t, "test", "go", "code")
		repo.Create(snippet)

		snippets, _ := repo.FindByCategory(0)

		if len(snippets) != 1 {
			t.Errorf("expected 1 uncategorized snippet, got %d", len(snippets))
		}
	})
}

func TestSnippetRepository_FindByTag(t *testing.T) {
	t.Run("finds snippets with tag", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet1 := mustCreateSnippet(t, "test1", "go", "code1")
		snippet1.AddTag(5)
		snippet2 := mustCreateSnippet(t, "test2", "go", "code2")
		snippet2.AddTag(10)
		snippet3 := mustCreateSnippet(t, "test3", "go", "code3")
		snippet3.AddTag(5)
		snippet3.AddTag(10)

		repo.Create(snippet1)
		repo.Create(snippet2)
		repo.Create(snippet3)

		snippets, err := repo.FindByTag(5)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 2 {
			t.Errorf("expected 2 snippets, got %d", len(snippets))
		}
	})

	t.Run("returns empty slice for tag with no snippets", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippets, err := repo.FindByTag(999)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 0 {
			t.Errorf("expected 0 snippets, got %d", len(snippets))
		}
	})
}

func TestSnippetRepository_FindByLanguage(t *testing.T) {
	t.Run("finds snippets with language", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet1 := mustCreateSnippet(t, "test1", "go", "code1")
		snippet2 := mustCreateSnippet(t, "test2", "python", "code2")
		snippet3 := mustCreateSnippet(t, "test3", "go", "code3")

		repo.Create(snippet1)
		repo.Create(snippet2)
		repo.Create(snippet3)

		snippets, err := repo.FindByLanguage("go")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 2 {
			t.Errorf("expected 2 snippets, got %d", len(snippets))
		}
	})

	t.Run("returns empty slice for language with no snippets", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippets, err := repo.FindByLanguage("rust")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 0 {
			t.Errorf("expected 0 snippets, got %d", len(snippets))
		}
	})
}

func TestSnippetRepository_Search(t *testing.T) {
	t.Run("finds snippets matching query", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet1 := mustCreateSnippet(t, "quicksort algorithm", "go", "func quicksort() {}")
		snippet2 := mustCreateSnippet(t, "hello world", "go", "fmt.Println()")

		repo.Create(snippet1)
		repo.Create(snippet2)

		snippets, err := repo.Search("quicksort")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 1 {
			t.Errorf("expected 1 snippet, got %d", len(snippets))
		}

		if len(snippets) > 0 && snippets[0].Title() != "quicksort algorithm" {
			t.Error("wrong snippet returned")
		}
	})

	t.Run("returns empty slice for no matches", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet := mustCreateSnippet(t, "test", "go", "code")
		repo.Create(snippet)

		snippets, err := repo.Search("nonexistent")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 0 {
			t.Errorf("expected 0 snippets, got %d", len(snippets))
		}
	})
}

func TestSnippetRepository_Create(t *testing.T) {
	t.Run("assigns ID to snippet", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet := mustCreateSnippet(t, "test", "go", "code")

		if snippet.ID() != 0 {
			t.Errorf("expected ID 0 before create, got %d", snippet.ID())
		}

		err := repo.Create(snippet)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if snippet.ID() != 1 {
			t.Errorf("expected ID 1 after create, got %d", snippet.ID())
		}
	})

	t.Run("increments IDs for multiple snippets", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet1 := mustCreateSnippet(t, "test1", "go", "code1")
		snippet2 := mustCreateSnippet(t, "test2", "go", "code2")
		snippet3 := mustCreateSnippet(t, "test3", "go", "code3")

		repo.Create(snippet1)
		repo.Create(snippet2)
		repo.Create(snippet3)

		if snippet1.ID() != 1 {
			t.Errorf("expected first ID 1, got %d", snippet1.ID())
		}

		if snippet2.ID() != 2 {
			t.Errorf("expected second ID 2, got %d", snippet2.ID())
		}

		if snippet3.ID() != 3 {
			t.Errorf("expected third ID 3, got %d", snippet3.ID())
		}
	})

	t.Run("stores snippet in repository", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet := mustCreateSnippet(t, "test", "go", "code")
		repo.Create(snippet)

		found, err := repo.FindByID(snippet.ID())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if found.Title() != "test" {
			t.Errorf("expected title %q, got %q", "test", found.Title())
		}
	})
}

func TestSnippetRepository_Update(t *testing.T) {
	t.Run("updates existing snippet", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet := mustCreateSnippet(t, "original", "go", "code")
		repo.Create(snippet)

		snippet.SetTitle("updated")
		err := repo.Update(snippet)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		found, _ := repo.FindByID(snippet.ID())
		if found.Title() != "updated" {
			t.Errorf("expected title %q, got %q", "updated", found.Title())
		}
	})

	t.Run("returns ErrNotFound for nonexistent snippet", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet := mustCreateSnippet(t, "test", "go", "code")
		snippet.SetID(999)

		err := repo.Update(snippet)

		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestSnippetRepository_Delete(t *testing.T) {
	t.Run("deletes existing snippet", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet := mustCreateSnippet(t, "test", "go", "code")
		repo.Create(snippet)
		id := snippet.ID()

		err := repo.Delete(id)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err = repo.FindByID(id)
		if err != ErrNotFound {
			t.Error("snippet still exists after delete")
		}
	})

	t.Run("returns ErrNotFound for nonexistent snippet", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		err := repo.Delete(999)

		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("does not reuse deleted IDs", func(t *testing.T) {
		s := newStore("test.json")
		repo := newSnippetRepository(s)

		snippet1 := mustCreateSnippet(t, "test1", "go", "code1")
		repo.Create(snippet1)

		snippet2 := mustCreateSnippet(t, "test2", "go", "code2")
		repo.Create(snippet2)

		repo.Delete(snippet1.ID())

		snippet3 := mustCreateSnippet(t, "test3", "go", "code3")
		repo.Create(snippet3)

		if snippet3.ID() != 3 {
			t.Errorf("expected ID 3, got %d (should not reuse deleted ID)", snippet3.ID())
		}
	})
}
