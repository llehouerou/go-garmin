package exercises

import (
	"slices"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	lib := Get()
	if lib == nil {
		t.Fatal("Get() returned nil")
	}

	// Should return same instance
	lib2 := Get()
	if lib != lib2 {
		t.Error("Get() should return singleton")
	}
}

func TestAll(t *testing.T) {
	lib := Get()
	all := lib.All()

	if len(all) == 0 {
		t.Error("All() returned empty slice")
	}

	// Should be 1794 exercises
	if len(all) != 1794 {
		t.Errorf("All() returned %d exercises, want 1794", len(all))
	}
}

func TestCategories(t *testing.T) {
	lib := Get()
	cats := lib.Categories()

	if len(cats) == 0 {
		t.Error("Categories() returned empty slice")
	}

	// Check sorted
	for i := 1; i < len(cats); i++ {
		if cats[i] < cats[i-1] {
			t.Error("Categories() not sorted")
			break
		}
	}

	// Check known category exists
	if !slices.Contains(cats, "BENCH_PRESS") {
		t.Error("Categories() missing BENCH_PRESS")
	}
}

func TestMuscles(t *testing.T) {
	lib := Get()
	muscles := lib.Muscles()

	if len(muscles) == 0 {
		t.Error("Muscles() returned empty slice")
	}

	// Check typos are fixed
	for _, m := range muscles {
		if m == "ADBUCTORS" || m == "ADDDUCTORS" || m == "SHOULDER" {
			t.Errorf("Muscles() contains unfixed typo: %s", m)
		}
	}

	// Check normalized versions exist
	found := map[string]bool{"ABDUCTORS": false, "ADDUCTORS": false, "SHOULDERS": false}
	for _, m := range muscles {
		if _, ok := found[m]; ok {
			found[m] = true
		}
	}
	for m, ok := range found {
		if !ok {
			t.Errorf("Muscles() missing normalized muscle: %s", m)
		}
	}
}

func TestEquipment(t *testing.T) {
	lib := Get()
	equip := lib.Equipment()

	if len(equip) == 0 {
		t.Error("Equipment() returned empty slice")
	}

	// Check known equipment exists
	if !slices.Contains(equip, "DUMBBELL") {
		t.Error("Equipment() missing DUMBBELL")
	}
}

func TestByCategory(t *testing.T) {
	lib := Get()

	exercises := lib.ByCategory("BENCH_PRESS")
	if len(exercises) == 0 {
		t.Error("ByCategory(BENCH_PRESS) returned empty")
	}

	for _, ex := range exercises {
		if ex.Category != "BENCH_PRESS" {
			t.Errorf("ByCategory returned exercise with wrong category: %s", ex.Category)
		}
	}

	// Case insensitive
	exercises2 := lib.ByCategory("bench_press")
	if len(exercises2) != len(exercises) {
		t.Error("ByCategory should be case-insensitive")
	}

	// Unknown category returns empty slice
	unknown := lib.ByCategory("UNKNOWN_CATEGORY")
	if unknown == nil {
		t.Error("ByCategory should return empty slice, not nil")
	}
	if len(unknown) != 0 {
		t.Error("ByCategory(unknown) should return empty slice")
	}
}

func TestByMuscle(t *testing.T) {
	lib := Get()

	exercises := lib.ByMuscle("CHEST")
	if len(exercises) == 0 {
		t.Error("ByMuscle(CHEST) returned empty")
	}

	// All should target chest (primary or secondary)
	for _, ex := range exercises {
		found := containsString(ex.PrimaryMuscles, "CHEST") || containsString(ex.SecondaryMuscles, "CHEST")
		if !found {
			t.Errorf("ByMuscle returned exercise not targeting CHEST: %s", ex.Key)
		}
	}
}

func TestByEquipment(t *testing.T) {
	lib := Get()

	exercises := lib.ByEquipment("DUMBBELL")
	if len(exercises) == 0 {
		t.Error("ByEquipment(DUMBBELL) returned empty")
	}

	for _, ex := range exercises {
		if !containsString(ex.Equipment, "DUMBBELL") {
			t.Errorf("ByEquipment returned exercise without DUMBBELL: %s", ex.Key)
		}
	}
}

func TestByKey(t *testing.T) {
	lib := Get()

	// CHEST_PRESS exists in multiple categories
	exercises := lib.ByKey("CHEST_PRESS")
	if len(exercises) == 0 {
		t.Error("ByKey(CHEST_PRESS) returned empty")
	}
	if len(exercises) < 2 {
		t.Error("ByKey(CHEST_PRESS) should return multiple exercises (different categories)")
	}

	for _, ex := range exercises {
		if ex.Key != "CHEST_PRESS" {
			t.Errorf("ByKey returned exercise with wrong key: %s", ex.Key)
		}
	}

	// Case insensitive
	exercises2 := lib.ByKey("chest_press")
	if len(exercises2) != len(exercises) {
		t.Error("ByKey should be case-insensitive")
	}
}

func TestSearch(t *testing.T) {
	lib := Get()

	// Search for CURL should match multiple exercises
	exercises := lib.Search("CURL")
	if len(exercises) == 0 {
		t.Error("Search(CURL) returned empty")
	}

	for _, ex := range exercises {
		if !containsSubstr(ex.Key, "CURL") {
			t.Errorf("Search returned exercise not matching CURL: %s", ex.Key)
		}
	}

	// Empty search returns all
	all := lib.Search("")
	if len(all) != len(lib.All()) {
		t.Error("Search('') should return all exercises")
	}
}

func TestFilter(t *testing.T) {
	lib := Get()

	// Filter by multiple criteria
	exercises := lib.Filter("CURL", "BICEPS", "DUMBBELL", "")
	if len(exercises) == 0 {
		t.Error("Filter(CURL, BICEPS, DUMBBELL) returned empty")
	}

	for _, ex := range exercises {
		if ex.Category != "CURL" {
			t.Errorf("Filter returned wrong category: %s", ex.Category)
		}
		if !containsString(ex.PrimaryMuscles, "BICEPS") && !containsString(ex.SecondaryMuscles, "BICEPS") {
			t.Errorf("Filter returned exercise not targeting BICEPS: %s", ex.Key)
		}
		if !containsString(ex.Equipment, "DUMBBELL") {
			t.Errorf("Filter returned exercise without DUMBBELL: %s", ex.Key)
		}
	}

	// Empty filters return all
	all := lib.Filter("", "", "", "")
	if len(all) != len(lib.All()) {
		t.Error("Filter with empty params should return all")
	}
}

func containsSubstr(s, substr string) bool {
	return strings.Contains(s, substr)
}
