package storage

import (
	"errors"
	"testing"

	"github.com/7-Dany/snip/internal/domain"
)

// setupCategoryTestStore creates a new store for testing
func setupCategoryTestStore() *categoryStore {
	store := newStore("test.json")
	return newCategoryStore(store)
}

// Test_categoryStore_List_Empty verifies List returns empty slice when no categories exist
func Test_categoryStore_List_Empty(t *testing.T) {
	cs := setupCategoryTestStore()

	categories, err := cs.List()

	if err != nil {
		t.Fatalf("List() error = %v, want nil", err)
	}

	if categories == nil {
		t.Fatal("List() returned nil, want empty slice")
	}

	if len(categories) != 0 {
		t.Errorf("List() returned %d categories, want 0", len(categories))
	}
}

// Test_categoryStore_List_Multiple verifies List returns all categories
func Test_categoryStore_List_Multiple(t *testing.T) {
	cs := setupCategoryTestStore()

	// Create test categories
	cat1, _ := domain.NewCategory("Algorithms")
	cat2, _ := domain.NewCategory("Data Structures")
	cat3, _ := domain.NewCategory("Web Development")

	cs.Create(cat1)
	cs.Create(cat2)
	cs.Create(cat3)

	categories, err := cs.List()

	if err != nil {
		t.Fatalf("List() error = %v, want nil", err)
	}

	if len(categories) != 3 {
		t.Errorf("List() returned %d categories, want 3", len(categories))
	}

	// Verify all categories are present (order doesn't matter with maps)
	found := make(map[int]bool)
	for _, cat := range categories {
		found[cat.ID()] = true
	}

	if !found[cat1.ID()] || !found[cat2.ID()] || !found[cat3.ID()] {
		t.Error("List() missing one or more created categories")
	}
}

// Test_categoryStore_FindByID_Success verifies finding category by valid ID
func Test_categoryStore_FindByID_Success(t *testing.T) {
	cs := setupCategoryTestStore()

	cat, _ := domain.NewCategory("Algorithms")
	cs.Create(cat)

	found, err := cs.FindByID(cat.ID())

	if err != nil {
		t.Fatalf("FindByID() error = %v, want nil", err)
	}

	if found.ID() != cat.ID() {
		t.Errorf("FindByID() returned category ID %d, want %d", found.ID(), cat.ID())
	}

	if found.Name() != "Algorithms" {
		t.Errorf("FindByID() returned category name %q, want %q", found.Name(), "Algorithms")
	}
}

// Test_categoryStore_FindByID_NotFound verifies error when ID doesn't exist
func Test_categoryStore_FindByID_NotFound(t *testing.T) {
	cs := setupCategoryTestStore()

	found, err := cs.FindByID(999)

	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("FindByID() error = %v, want %v", err, domain.ErrNotFound)
	}

	if found != nil {
		t.Errorf("FindByID() returned %v, want nil", found)
	}
}

// Test_categoryStore_FindByName_Success verifies finding category by exact name
func Test_categoryStore_FindByName_Success(t *testing.T) {
	cs := setupCategoryTestStore()

	cat, _ := domain.NewCategory("Algorithms")
	cs.Create(cat)

	found, err := cs.FindByName("Algorithms")

	if err != nil {
		t.Fatalf("FindByName() error = %v, want nil", err)
	}

	if found.ID() != cat.ID() {
		t.Errorf("FindByName() returned category ID %d, want %d", found.ID(), cat.ID())
	}

	if found.Name() != "Algorithms" {
		t.Errorf("FindByName() returned category name %q, want %q", found.Name(), "Algorithms")
	}
}

// Test_categoryStore_FindByName_CaseInsensitive verifies case-insensitive matching
func Test_categoryStore_FindByName_CaseInsensitive(t *testing.T) {
	cs := setupCategoryTestStore()

	cat, _ := domain.NewCategory("Algorithms")
	cs.Create(cat)

	testCases := []struct {
		name      string
		searchFor string
	}{
		{"lowercase", "algorithms"},
		{"uppercase", "ALGORITHMS"},
		{"mixed case", "AlGoRiThMs"},
		{"exact match", "Algorithms"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			found, err := cs.FindByName(tc.searchFor)

			if err != nil {
				t.Fatalf("FindByName(%q) error = %v, want nil", tc.searchFor, err)
			}

			if found.ID() != cat.ID() {
				t.Errorf("FindByName(%q) returned wrong category ID", tc.searchFor)
			}
		})
	}
}

// Test_categoryStore_FindByName_NotFound verifies error when name doesn't exist
func Test_categoryStore_FindByName_NotFound(t *testing.T) {
	cs := setupCategoryTestStore()

	cat, _ := domain.NewCategory("Algorithms")
	cs.Create(cat)

	found, err := cs.FindByName("NonExistent")

	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("FindByName() error = %v, want %v", err, domain.ErrNotFound)
	}

	if found != nil {
		t.Errorf("FindByName() returned %v, want nil", found)
	}
}

// Test_categoryStore_FindByName_FirstMatch verifies first match is returned when duplicates exist
func Test_categoryStore_FindByName_FirstMatch(t *testing.T) {
	cs := setupCategoryTestStore()

	// Create two categories with same name (storage allows this)
	cat1, _ := domain.NewCategory("Duplicate")
	cat2, _ := domain.NewCategory("Duplicate")

	cs.Create(cat1)
	cs.Create(cat2)

	found, err := cs.FindByName("Duplicate")

	if err != nil {
		t.Fatalf("FindByName() error = %v, want nil", err)
	}

	// Should return one of them (implementation returns first found in map iteration)
	if found.ID() != cat1.ID() && found.ID() != cat2.ID() {
		t.Error("FindByName() returned unexpected category")
	}
}

// Test_categoryStore_Create_AssignsID verifies Create assigns auto-incremented ID
func Test_categoryStore_Create_AssignsID(t *testing.T) {
	cs := setupCategoryTestStore()

	cat, _ := domain.NewCategory("Algorithms")

	// Before Create, ID should be 0 (zero value)
	if cat.ID() != 0 {
		t.Errorf("Category ID before Create = %d, want 0", cat.ID())
	}

	err := cs.Create(cat)

	if err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}

	// After Create, ID should be 1 (first ID)
	if cat.ID() != 1 {
		t.Errorf("Category ID after Create = %d, want 1", cat.ID())
	}
}

// Test_categoryStore_Create_IncrementsID verifies IDs increment correctly
func Test_categoryStore_Create_IncrementsID(t *testing.T) {
	cs := setupCategoryTestStore()

	cat1, _ := domain.NewCategory("First")
	cat2, _ := domain.NewCategory("Second")
	cat3, _ := domain.NewCategory("Third")

	cs.Create(cat1)
	cs.Create(cat2)
	cs.Create(cat3)

	if cat1.ID() != 1 {
		t.Errorf("First category ID = %d, want 1", cat1.ID())
	}

	if cat2.ID() != 2 {
		t.Errorf("Second category ID = %d, want 2", cat2.ID())
	}

	if cat3.ID() != 3 {
		t.Errorf("Third category ID = %d, want 3", cat3.ID())
	}
}

// Test_categoryStore_Create_StoresCategory verifies created category can be retrieved
func Test_categoryStore_Create_StoresCategory(t *testing.T) {
	cs := setupCategoryTestStore()

	cat, _ := domain.NewCategory("Algorithms")
	cs.Create(cat)

	found, err := cs.FindByID(cat.ID())

	if err != nil {
		t.Fatalf("FindByID() after Create error = %v, want nil", err)
	}

	if found.Name() != "Algorithms" {
		t.Errorf("Stored category name = %q, want %q", found.Name(), "Algorithms")
	}
}

// Test_categoryStore_Update_Success verifies updating existing category
func Test_categoryStore_Update_Success(t *testing.T) {
	cs := setupCategoryTestStore()

	cat, _ := domain.NewCategory("Old Name")
	cs.Create(cat)

	// Update name
	cat.SetName("New Name")
	err := cs.Update(cat)

	if err != nil {
		t.Fatalf("Update() error = %v, want nil", err)
	}

	// Verify update persisted
	found, _ := cs.FindByID(cat.ID())
	if found.Name() != "New Name" {
		t.Errorf("Updated category name = %q, want %q", found.Name(), "New Name")
	}
}

// Test_categoryStore_Update_NotFound verifies error when updating non-existent category
func Test_categoryStore_Update_NotFound(t *testing.T) {
	cs := setupCategoryTestStore()

	cat, _ := domain.NewCategory("Test")
	cat.SetID(999) // Non-existent ID

	err := cs.Update(cat)

	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Update() error = %v, want %v", err, domain.ErrNotFound)
	}
}

// Test_categoryStore_Delete_Success verifies deleting existing category
func Test_categoryStore_Delete_Success(t *testing.T) {
	cs := setupCategoryTestStore()

	cat, _ := domain.NewCategory("ToDelete")
	cs.Create(cat)

	err := cs.Delete(cat.ID())

	if err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}

	// Verify category no longer exists
	found, err := cs.FindByID(cat.ID())
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("FindByID() after Delete error = %v, want %v", err, domain.ErrNotFound)
	}

	if found != nil {
		t.Error("FindByID() after Delete returned non-nil category")
	}
}

// Test_categoryStore_Delete_NotFound verifies error when deleting non-existent category
func Test_categoryStore_Delete_NotFound(t *testing.T) {
	cs := setupCategoryTestStore()

	err := cs.Delete(999)

	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Delete() error = %v, want %v", err, domain.ErrNotFound)
	}
}

// Test_categoryStore_Delete_DoesNotAffectOthers verifies delete only removes target category
func Test_categoryStore_Delete_DoesNotAffectOthers(t *testing.T) {
	cs := setupCategoryTestStore()

	cat1, _ := domain.NewCategory("Keep1")
	cat2, _ := domain.NewCategory("Delete")
	cat3, _ := domain.NewCategory("Keep2")

	cs.Create(cat1)
	cs.Create(cat2)
	cs.Create(cat3)

	cs.Delete(cat2.ID())

	// Verify other categories still exist
	found1, err1 := cs.FindByID(cat1.ID())
	found3, err3 := cs.FindByID(cat3.ID())

	if err1 != nil || found1 == nil {
		t.Error("Delete() affected unrelated category 1")
	}

	if err3 != nil || found3 == nil {
		t.Error("Delete() affected unrelated category 3")
	}

	// Verify deleted category is gone
	_, err2 := cs.FindByID(cat2.ID())
	if !errors.Is(err2, domain.ErrNotFound) {
		t.Error("Delete() did not remove target category")
	}
}

// Test_categoryStore_Delete_DoesNotReuseID verifies deleted IDs are not reused
func Test_categoryStore_Delete_DoesNotReuseID(t *testing.T) {
	cs := setupCategoryTestStore()

	cat1, _ := domain.NewCategory("First")
	cs.Create(cat1) // ID = 1

	cat2, _ := domain.NewCategory("Second")
	cs.Create(cat2) // ID = 2

	cs.Delete(cat1.ID()) // Delete ID 1

	cat3, _ := domain.NewCategory("Third")
	cs.Create(cat3) // Should get ID = 3 (not reuse ID 1)

	if cat3.ID() != 3 {
		t.Errorf("New category after delete got ID %d, want 3 (should not reuse deleted ID)", cat3.ID())
	}

	// Verify metadata incremented correctly
	if cs.store.metadata.next_category_id != 4 {
		t.Errorf("Metadata next_category_id = %d, want 4", cs.store.metadata.next_category_id)
	}
}
