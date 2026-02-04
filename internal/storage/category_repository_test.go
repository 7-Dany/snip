package storage

import (
	"testing"
)

func TestCategoryRepository_List(t *testing.T) {
	t.Run("returns empty slice when no categories", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		categories, err := repo.List()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if categories == nil {
			t.Fatal("expected empty slice, got nil")
		}

		if len(categories) != 0 {
			t.Errorf("expected 0 categories, got %d", len(categories))
		}
	})

	t.Run("returns all categories", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		cat1 := mustCreateCategory(t, "algorithms")
		cat2 := mustCreateCategory(t, "data-structures")
		cat3 := mustCreateCategory(t, "web-dev")

		repo.Create(cat1)
		repo.Create(cat2)
		repo.Create(cat3)

		categories, err := repo.List()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(categories) != 3 {
			t.Errorf("expected 3 categories, got %d", len(categories))
		}
	})

	t.Run("returns defensive copy", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		original := mustCreateCategory(t, "test")
		repo.Create(original)

		categories, _ := repo.List()
		categories[0] = nil

		// Original should be unchanged
		stored, _ := repo.FindByID(original.ID())
		if stored == nil {
			t.Error("modifying returned slice affected internal store")
		}
	})
}

func TestCategoryRepository_FindByID(t *testing.T) {
	t.Run("finds existing category", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		category := mustCreateCategory(t, "algorithms")
		repo.Create(category)

		found, err := repo.FindByID(category.ID())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if found.ID() != category.ID() {
			t.Errorf("expected ID %d, got %d", category.ID(), found.ID())
		}

		if found.Name() != "algorithms" {
			t.Errorf("expected name %q, got %q", "algorithms", found.Name())
		}
	})

	t.Run("returns ErrNotFound for nonexistent ID", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		_, err := repo.FindByID(999)

		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestCategoryRepository_FindByName(t *testing.T) {
	t.Run("finds existing category by name", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		category := mustCreateCategory(t, "algorithms")
		repo.Create(category)

		found, err := repo.FindByName("algorithms")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if found.ID() != category.ID() {
			t.Errorf("expected ID %d, got %d", category.ID(), found.ID())
		}
	})

	t.Run("returns ErrNotFound for nonexistent name", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		_, err := repo.FindByName("nonexistent")

		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestCategoryRepository_Create(t *testing.T) {
	t.Run("assigns ID to category", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		category := mustCreateCategory(t, "algorithms")

		if category.ID() != 0 {
			t.Errorf("expected ID 0 before create, got %d", category.ID())
		}

		err := repo.Create(category)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if category.ID() != 1 {
			t.Errorf("expected ID 1 after create, got %d", category.ID())
		}
	})

	t.Run("increments IDs for multiple categories", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		cat1 := mustCreateCategory(t, "algorithms")
		cat2 := mustCreateCategory(t, "data-structures")
		cat3 := mustCreateCategory(t, "web-dev")

		repo.Create(cat1)
		repo.Create(cat2)
		repo.Create(cat3)

		if cat1.ID() != 1 {
			t.Errorf("expected first ID 1, got %d", cat1.ID())
		}

		if cat2.ID() != 2 {
			t.Errorf("expected second ID 2, got %d", cat2.ID())
		}

		if cat3.ID() != 3 {
			t.Errorf("expected third ID 3, got %d", cat3.ID())
		}
	})

	t.Run("stores category in repository", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		category := mustCreateCategory(t, "algorithms")
		repo.Create(category)

		found, err := repo.FindByID(category.ID())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if found.Name() != "algorithms" {
			t.Errorf("expected name %q, got %q", "algorithms", found.Name())
		}
	})

	t.Run("returns ErrDuplicateName for duplicate name", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		cat1 := mustCreateCategory(t, "algorithms")
		repo.Create(cat1)

		cat2 := mustCreateCategory(t, "algorithms")
		err := repo.Create(cat2)

		if err != ErrDuplicateName {
			t.Errorf("expected ErrDuplicateName, got %v", err)
		}
	})
}

func TestCategoryRepository_Update(t *testing.T) {
	t.Run("updates existing category", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		category := mustCreateCategory(t, "original")
		repo.Create(category)

		category.SetName("updated")
		err := repo.Update(category)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		found, _ := repo.FindByID(category.ID())
		if found.Name() != "updated" {
			t.Errorf("expected name %q, got %q", "updated", found.Name())
		}
	})

	t.Run("returns ErrNotFound for nonexistent category", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		category := mustCreateCategory(t, "test")
		category.SetID(999)

		err := repo.Update(category)

		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrDuplicateName when name conflicts with another category", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		cat1 := mustCreateCategory(t, "algorithms")
		cat2 := mustCreateCategory(t, "data-structures")
		repo.Create(cat1)
		repo.Create(cat2)

		cat2.SetName("algorithms")
		err := repo.Update(cat2)

		if err != ErrDuplicateName {
			t.Errorf("expected ErrDuplicateName, got %v", err)
		}
	})

	t.Run("allows updating to same name", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		category := mustCreateCategory(t, "algorithms")
		repo.Create(category)

		// Update with same name should succeed
		category.SetName("algorithms")
		err := repo.Update(category)

		if err != nil {
			t.Errorf("expected no error when updating to same name, got %v", err)
		}
	})
}

func TestCategoryRepository_Delete(t *testing.T) {
	t.Run("deletes existing category", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		category := mustCreateCategory(t, "algorithms")
		repo.Create(category)
		id := category.ID()

		err := repo.Delete(id)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err = repo.FindByID(id)
		if err != ErrNotFound {
			t.Error("category still exists after delete")
		}
	})

	t.Run("returns ErrNotFound for nonexistent category", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		err := repo.Delete(999)

		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("does not affect other categories", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		cat1 := mustCreateCategory(t, "keep1")
		cat2 := mustCreateCategory(t, "delete")
		cat3 := mustCreateCategory(t, "keep2")

		repo.Create(cat1)
		repo.Create(cat2)
		repo.Create(cat3)

		repo.Delete(cat2.ID())

		// Verify other categories still exist
		found1, err1 := repo.FindByID(cat1.ID())
		found3, err3 := repo.FindByID(cat3.ID())

		if err1 != nil || found1 == nil {
			t.Error("delete affected unrelated category 1")
		}

		if err3 != nil || found3 == nil {
			t.Error("delete affected unrelated category 3")
		}
	})

	t.Run("does not reuse deleted IDs", func(t *testing.T) {
		s := newStore("test.json")
		repo := newCategoryRepository(s)

		cat1 := mustCreateCategory(t, "first")
		repo.Create(cat1)

		cat2 := mustCreateCategory(t, "second")
		repo.Create(cat2)

		repo.Delete(cat1.ID())

		cat3 := mustCreateCategory(t, "third")
		repo.Create(cat3)

		if cat3.ID() != 3 {
			t.Errorf("expected ID 3, got %d (should not reuse deleted ID)", cat3.ID())
		}
	})
}
