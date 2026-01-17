package storage

import (
	"errors"
	"testing"

	"github.com/7-Dany/snip/internal/domain"
)

func setupTagTestStore() *tagStore {
	store := newStore("test.json")
	return newTagStore(store)
}

// List tests

func TestTagStore_List_Empty(t *testing.T) {
	ts := setupTagTestStore()

	tags, err := ts.List()
	if err != nil {
		t.Fatalf("List() returned error: %v", err)
	}

	if tags == nil {
		t.Fatal("List() returned nil instead of empty slice")
	}

	if len(tags) != 0 {
		t.Errorf("List() on empty store returned %d tags, expected 0", len(tags))
	}
}

func TestTagStore_List_Multiple(t *testing.T) {
	ts := setupTagTestStore()

	tag1, _ := domain.NewTag("sorting")
	tag2, _ := domain.NewTag("algorithm")
	tag3, _ := domain.NewTag("optimization")

	ts.Create(tag1)
	ts.Create(tag2)
	ts.Create(tag3)

	tags, err := ts.List()
	if err != nil {
		t.Fatalf("List() returned error: %v", err)
	}

	if len(tags) != 3 {
		t.Errorf("List() returned %d tags, expected 3", len(tags))
	}
}

// FindByID tests

func TestTagStore_FindByID_Success(t *testing.T) {
	ts := setupTagTestStore()

	tag, _ := domain.NewTag("sorting")
	ts.Create(tag)

	found, err := ts.FindByID(tag.ID())
	if err != nil {
		t.Fatalf("FindByID() returned error: %v", err)
	}

	if found.Name() != "sorting" {
		t.Errorf("FindByID() returned tag with name %q, expected %q", found.Name(), "sorting")
	}
}

func TestTagStore_FindByID_NotFound(t *testing.T) {
	ts := setupTagTestStore()

	_, err := ts.FindByID(999)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("FindByID(999) error = %v, expected domain.ErrNotFound", err)
	}
}

// FindByName tests

func TestTagStore_FindByName_ExactMatch(t *testing.T) {
	ts := setupTagTestStore()

	tag, _ := domain.NewTag("sorting")
	ts.Create(tag)

	found, err := ts.FindByName("sorting")
	if err != nil {
		t.Fatalf("FindByName() returned error: %v", err)
	}

	if found.ID() != tag.ID() {
		t.Errorf("FindByName() returned tag with ID %d, expected %d", found.ID(), tag.ID())
	}
}

func TestTagStore_FindByName_CaseInsensitive(t *testing.T) {
	ts := setupTagTestStore()

	tag, _ := domain.NewTag("Sorting")
	ts.Create(tag)

	tests := []string{"sorting", "SORTING", "SoRtInG"}
	for _, name := range tests {
		found, err := ts.FindByName(name)
		if err != nil {
			t.Errorf("FindByName(%q) returned error: %v", name, err)
			continue
		}

		if found.ID() != tag.ID() {
			t.Errorf("FindByName(%q) returned wrong tag", name)
		}
	}
}

func TestTagStore_FindByName_NotFound(t *testing.T) {
	ts := setupTagTestStore()

	_, err := ts.FindByName("nonexistent")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("FindByName() error = %v, expected domain.ErrNotFound", err)
	}
}

func TestTagStore_FindByName_FirstMatch(t *testing.T) {
	ts := setupTagTestStore()

	// Create two tags with same name (storage allows this)
	tag1, _ := domain.NewTag("duplicate")
	tag2, _ := domain.NewTag("duplicate")
	ts.Create(tag1)
	ts.Create(tag2)

	found, err := ts.FindByName("duplicate")
	if err != nil {
		t.Fatalf("FindByName() returned error: %v", err)
	}

	// Should return one of them (implementation returns first encountered)
	if found.ID() != tag1.ID() && found.ID() != tag2.ID() {
		t.Errorf("FindByName() returned unexpected tag ID %d", found.ID())
	}
}

// Create tests

func TestTagStore_Create_AssignsID(t *testing.T) {
	ts := setupTagTestStore()

	tag, _ := domain.NewTag("sorting")

	err := ts.Create(tag)
	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	if tag.ID() != 1 {
		t.Errorf("Create() assigned ID %d, expected 1", tag.ID())
	}
}

func TestTagStore_Create_IncrementsMetadata(t *testing.T) {
	ts := setupTagTestStore()

	tag1, _ := domain.NewTag("sorting")
	tag2, _ := domain.NewTag("algorithm")

	ts.Create(tag1)
	ts.Create(tag2)

	if tag1.ID() != 1 {
		t.Errorf("First tag ID = %d, expected 1", tag1.ID())
	}

	if tag2.ID() != 2 {
		t.Errorf("Second tag ID = %d, expected 2", tag2.ID())
	}

	if ts.store.metadata.next_tag_id != 3 {
		t.Errorf("Metadata next_tag_id = %d, expected 3", ts.store.metadata.next_tag_id)
	}
}

func TestTagStore_Create_StoresInMap(t *testing.T) {
	ts := setupTagTestStore()

	tag, _ := domain.NewTag("sorting")
	ts.Create(tag)

	stored, ok := ts.store.tags[tag.ID()]
	if !ok {
		t.Fatal("Created tag not found in store.tags map")
	}

	if stored.Name() != "sorting" {
		t.Errorf("Stored tag has name %q, expected %q", stored.Name(), "sorting")
	}
}

// Update tests

func TestTagStore_Update_Success(t *testing.T) {
	ts := setupTagTestStore()

	tag, _ := domain.NewTag("sorting")
	ts.Create(tag)

	tag.SetName("quick-sort")

	err := ts.Update(tag)
	if err != nil {
		t.Fatalf("Update() returned error: %v", err)
	}

	found, _ := ts.FindByID(tag.ID())
	if found.Name() != "quick-sort" {
		t.Errorf("Updated tag has name %q, expected %q", found.Name(), "quick-sort")
	}
}

func TestTagStore_Update_NotFound(t *testing.T) {
	ts := setupTagTestStore()

	tag, _ := domain.NewTag("sorting")
	tag.SetID(999) // Non-existent ID

	err := ts.Update(tag)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Update() error = %v, expected domain.ErrNotFound", err)
	}
}

// Delete tests

func TestTagStore_Delete_Success(t *testing.T) {
	ts := setupTagTestStore()

	tag, _ := domain.NewTag("sorting")
	ts.Create(tag)

	err := ts.Delete(tag.ID())
	if err != nil {
		t.Fatalf("Delete() returned error: %v", err)
	}

	_, err = ts.FindByID(tag.ID())
	if !errors.Is(err, domain.ErrNotFound) {
		t.Error("Deleted tag still found in store")
	}
}

func TestTagStore_Delete_NotFound(t *testing.T) {
	ts := setupTagTestStore()

	err := ts.Delete(999)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Delete() error = %v, expected domain.ErrNotFound", err)
	}
}

func TestTagStore_Delete_DoesNotAffectOthers(t *testing.T) {
	ts := setupTagTestStore()

	tag1, _ := domain.NewTag("sorting")
	tag2, _ := domain.NewTag("algorithm")
	ts.Create(tag1)
	ts.Create(tag2)

	ts.Delete(tag1.ID())

	// tag2 should still exist
	found, err := ts.FindByID(tag2.ID())
	if err != nil {
		t.Fatal("Delete() affected other tags")
	}

	if found.Name() != "algorithm" {
		t.Error("Delete() corrupted other tags")
	}
}

func TestTagStore_Delete_IDNotReused(t *testing.T) {
	ts := setupTagTestStore()

	tag1, _ := domain.NewTag("sorting")
	ts.Create(tag1) // ID = 1

	ts.Delete(tag1.ID())

	tag2, _ := domain.NewTag("algorithm")
	ts.Create(tag2) // Should be ID = 2, not reusing 1

	if tag2.ID() == tag1.ID() {
		t.Error("Delete() caused ID reuse")
	}

	if tag2.ID() != 2 {
		t.Errorf("New tag after delete has ID %d, expected 2", tag2.ID())
	}
}
