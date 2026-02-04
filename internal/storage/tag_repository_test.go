package storage

import (
	"testing"
)

func TestTagRepository_List(t *testing.T) {
	t.Run("returns empty slice when no tags", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tags, err := repo.List()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if tags == nil {
			t.Fatal("expected empty slice, got nil")
		}

		if len(tags) != 0 {
			t.Errorf("expected 0 tags, got %d", len(tags))
		}
	})

	t.Run("returns all tags", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tag1 := mustCreateTag(t, "sorting")
		tag2 := mustCreateTag(t, "algorithm")
		tag3 := mustCreateTag(t, "recursion")

		repo.Create(tag1)
		repo.Create(tag2)
		repo.Create(tag3)

		tags, err := repo.List()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(tags) != 3 {
			t.Errorf("expected 3 tags, got %d", len(tags))
		}
	})

	t.Run("returns defensive copy", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		original := mustCreateTag(t, "test")
		repo.Create(original)

		tags, _ := repo.List()
		tags[0] = nil

		// Original should be unchanged
		stored, _ := repo.FindByID(original.ID())
		if stored == nil {
			t.Error("modifying returned slice affected internal store")
		}
	})
}

func TestTagRepository_FindByID(t *testing.T) {
	t.Run("finds existing tag", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tag := mustCreateTag(t, "sorting")
		repo.Create(tag)

		found, err := repo.FindByID(tag.ID())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if found.ID() != tag.ID() {
			t.Errorf("expected ID %d, got %d", tag.ID(), found.ID())
		}

		if found.Name() != "sorting" {
			t.Errorf("expected name %q, got %q", "sorting", found.Name())
		}
	})

	t.Run("returns ErrNotFound for nonexistent ID", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		_, err := repo.FindByID(999)

		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestTagRepository_FindByName(t *testing.T) {
	t.Run("finds existing tag by name", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tag := mustCreateTag(t, "sorting")
		repo.Create(tag)

		found, err := repo.FindByName("sorting")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if found.ID() != tag.ID() {
			t.Errorf("expected ID %d, got %d", tag.ID(), found.ID())
		}
	})

	t.Run("returns ErrNotFound for nonexistent name", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		_, err := repo.FindByName("nonexistent")

		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestTagRepository_Create(t *testing.T) {
	t.Run("assigns ID to tag", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tag := mustCreateTag(t, "sorting")

		if tag.ID() != 0 {
			t.Errorf("expected ID 0 before create, got %d", tag.ID())
		}

		err := repo.Create(tag)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if tag.ID() != 1 {
			t.Errorf("expected ID 1 after create, got %d", tag.ID())
		}
	})

	t.Run("increments IDs for multiple tags", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tag1 := mustCreateTag(t, "sorting")
		tag2 := mustCreateTag(t, "algorithm")
		tag3 := mustCreateTag(t, "recursion")

		repo.Create(tag1)
		repo.Create(tag2)
		repo.Create(tag3)

		if tag1.ID() != 1 {
			t.Errorf("expected first ID 1, got %d", tag1.ID())
		}

		if tag2.ID() != 2 {
			t.Errorf("expected second ID 2, got %d", tag2.ID())
		}

		if tag3.ID() != 3 {
			t.Errorf("expected third ID 3, got %d", tag3.ID())
		}
	})

	t.Run("stores tag in repository", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tag := mustCreateTag(t, "sorting")
		repo.Create(tag)

		found, err := repo.FindByID(tag.ID())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if found.Name() != "sorting" {
			t.Errorf("expected name %q, got %q", "sorting", found.Name())
		}
	})

	t.Run("returns ErrDuplicateName for duplicate name", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tag1 := mustCreateTag(t, "sorting")
		repo.Create(tag1)

		tag2 := mustCreateTag(t, "sorting")
		err := repo.Create(tag2)

		if err != ErrDuplicateName {
			t.Errorf("expected ErrDuplicateName, got %v", err)
		}
	})
}

func TestTagRepository_Update(t *testing.T) {
	t.Run("updates existing tag", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tag := mustCreateTag(t, "original")
		repo.Create(tag)

		tag.SetName("updated")
		err := repo.Update(tag)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		found, _ := repo.FindByID(tag.ID())
		if found.Name() != "updated" {
			t.Errorf("expected name %q, got %q", "updated", found.Name())
		}
	})

	t.Run("returns ErrNotFound for nonexistent tag", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tag := mustCreateTag(t, "test")
		tag.SetID(999)

		err := repo.Update(tag)

		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns ErrDuplicateName when name conflicts with another tag", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tag1 := mustCreateTag(t, "sorting")
		tag2 := mustCreateTag(t, "algorithm")
		repo.Create(tag1)
		repo.Create(tag2)

		tag2.SetName("sorting")
		err := repo.Update(tag2)

		if err != ErrDuplicateName {
			t.Errorf("expected ErrDuplicateName, got %v", err)
		}
	})

	t.Run("allows updating to same name", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tag := mustCreateTag(t, "sorting")
		repo.Create(tag)

		// Update with same name should succeed
		tag.SetName("sorting")
		err := repo.Update(tag)

		if err != nil {
			t.Errorf("expected no error when updating to same name, got %v", err)
		}
	})
}

func TestTagRepository_Delete(t *testing.T) {
	t.Run("deletes existing tag", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tag := mustCreateTag(t, "sorting")
		repo.Create(tag)
		id := tag.ID()

		err := repo.Delete(id)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err = repo.FindByID(id)
		if err != ErrNotFound {
			t.Error("tag still exists after delete")
		}
	})

	t.Run("returns ErrNotFound for nonexistent tag", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		err := repo.Delete(999)

		if err != ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("does not affect other tags", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tag1 := mustCreateTag(t, "keep1")
		tag2 := mustCreateTag(t, "delete")
		tag3 := mustCreateTag(t, "keep2")

		repo.Create(tag1)
		repo.Create(tag2)
		repo.Create(tag3)

		repo.Delete(tag2.ID())

		// Verify other tags still exist
		found1, err1 := repo.FindByID(tag1.ID())
		found3, err3 := repo.FindByID(tag3.ID())

		if err1 != nil || found1 == nil {
			t.Error("delete affected unrelated tag 1")
		}

		if err3 != nil || found3 == nil {
			t.Error("delete affected unrelated tag 3")
		}
	})

	t.Run("does not reuse deleted IDs", func(t *testing.T) {
		s := newStore("test.json")
		repo := newTagRepository(s)

		tag1 := mustCreateTag(t, "first")
		repo.Create(tag1)

		tag2 := mustCreateTag(t, "second")
		repo.Create(tag2)

		repo.Delete(tag1.ID())

		tag3 := mustCreateTag(t, "third")
		repo.Create(tag3)

		if tag3.ID() != 3 {
			t.Errorf("expected ID 3, got %d (should not reuse deleted ID)", tag3.ID())
		}
	})
}
