// Package exercises provides a static library of Garmin strength training exercises.
package exercises

import (
	_ "embed"
	"encoding/json"
	"slices"
	"strings"
	"sync"
)

//go:embed data.json
var dataJSON []byte

// Exercise represents a strength training exercise.
type Exercise struct {
	Key              string   `json:"key"`
	Category         string   `json:"category"`
	PrimaryMuscles   []string `json:"primaryMuscles"`
	SecondaryMuscles []string `json:"secondaryMuscles"`
	Equipment        []string `json:"equipment"`
}

// Library provides access to the exercise database with pre-built indexes.
type Library struct {
	exercises  []Exercise
	byCategory map[string][]Exercise
	byMuscle   map[string][]Exercise
	byEquip    map[string][]Exercise
	byKey      map[string][]Exercise
	categories []string
	muscles    []string
	equipment  []string
}

var (
	library     *Library
	libraryOnce sync.Once
)

// Get returns the global exercise library.
func Get() *Library {
	libraryOnce.Do(func() {
		library = mustLoadLibrary()
	})
	return library
}

func mustLoadLibrary() *Library {
	var exercises []Exercise
	if err := json.Unmarshal(dataJSON, &exercises); err != nil {
		panic("exercises: failed to load embedded data: " + err.Error())
	}

	l := &Library{
		exercises:  exercises,
		byCategory: make(map[string][]Exercise),
		byMuscle:   make(map[string][]Exercise),
		byEquip:    make(map[string][]Exercise),
		byKey:      make(map[string][]Exercise),
	}

	categorySet := make(map[string]bool)
	muscleSet := make(map[string]bool)
	equipSet := make(map[string]bool)

	for _, ex := range exercises {
		// Index by category
		l.byCategory[ex.Category] = append(l.byCategory[ex.Category], ex)
		categorySet[ex.Category] = true

		// Index by key
		l.byKey[ex.Key] = append(l.byKey[ex.Key], ex)

		// Index by muscle (both primary and secondary)
		for _, m := range ex.PrimaryMuscles {
			l.byMuscle[m] = append(l.byMuscle[m], ex)
			muscleSet[m] = true
		}
		for _, m := range ex.SecondaryMuscles {
			l.byMuscle[m] = append(l.byMuscle[m], ex)
			muscleSet[m] = true
		}

		// Index by equipment
		for _, e := range ex.Equipment {
			l.byEquip[e] = append(l.byEquip[e], ex)
			equipSet[e] = true
		}
	}

	// Build sorted lists
	l.categories = sortedKeys(categorySet)
	l.muscles = sortedKeys(muscleSet)
	l.equipment = sortedKeys(equipSet)

	return l
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Simple insertion sort for small sets
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}

// All returns all exercises.
func (l *Library) All() []Exercise {
	result := make([]Exercise, len(l.exercises))
	copy(result, l.exercises)
	return result
}

// Categories returns all exercise category names.
func (l *Library) Categories() []string {
	result := make([]string, len(l.categories))
	copy(result, l.categories)
	return result
}

// Muscles returns all muscle group names.
func (l *Library) Muscles() []string {
	result := make([]string, len(l.muscles))
	copy(result, l.muscles)
	return result
}

// Equipment returns all equipment type names.
func (l *Library) Equipment() []string {
	result := make([]string, len(l.equipment))
	copy(result, l.equipment)
	return result
}

// ByCategory returns exercises in the given category (case-insensitive).
func (l *Library) ByCategory(category string) []Exercise {
	return copySlice(l.byCategory[strings.ToUpper(category)])
}

// ByMuscle returns exercises targeting the given muscle (primary or secondary, case-insensitive).
func (l *Library) ByMuscle(muscle string) []Exercise {
	return copySlice(l.byMuscle[strings.ToUpper(muscle)])
}

// ByEquipment returns exercises using the given equipment (case-insensitive).
func (l *Library) ByEquipment(equipment string) []Exercise {
	return copySlice(l.byEquip[strings.ToUpper(equipment)])
}

// ByKey returns all exercises with the given key (case-insensitive).
// Multiple exercises may share a key across different categories.
func (l *Library) ByKey(key string) []Exercise {
	return copySlice(l.byKey[strings.ToUpper(key)])
}

// Search returns exercises whose key contains the query (case-insensitive).
func (l *Library) Search(query string) []Exercise {
	if query == "" {
		return l.All()
	}
	query = strings.ToUpper(query)
	var result []Exercise
	for _, ex := range l.exercises {
		if strings.Contains(ex.Key, query) {
			result = append(result, ex)
		}
	}
	return result
}

// Filter returns exercises matching all specified criteria.
// Empty strings are ignored. Multiple criteria are ANDed together.
func (l *Library) Filter(category, muscle, equipment, search string) []Exercise {
	result := l.exercises

	if category != "" {
		result = filterByCategory(result, strings.ToUpper(category))
	}
	if muscle != "" {
		result = filterByMuscle(result, strings.ToUpper(muscle))
	}
	if equipment != "" {
		result = filterByEquipment(result, strings.ToUpper(equipment))
	}
	if search != "" {
		result = filterBySearch(result, strings.ToUpper(search))
	}

	return copySlice(result)
}

func filterByCategory(exercises []Exercise, category string) []Exercise {
	var result []Exercise
	for _, ex := range exercises {
		if ex.Category == category {
			result = append(result, ex)
		}
	}
	return result
}

func filterByMuscle(exercises []Exercise, muscle string) []Exercise {
	var result []Exercise
	for _, ex := range exercises {
		if containsString(ex.PrimaryMuscles, muscle) || containsString(ex.SecondaryMuscles, muscle) {
			result = append(result, ex)
		}
	}
	return result
}

func filterByEquipment(exercises []Exercise, equipment string) []Exercise {
	var result []Exercise
	for _, ex := range exercises {
		if containsString(ex.Equipment, equipment) {
			result = append(result, ex)
		}
	}
	return result
}

func filterBySearch(exercises []Exercise, query string) []Exercise {
	var result []Exercise
	for _, ex := range exercises {
		if strings.Contains(ex.Key, query) {
			result = append(result, ex)
		}
	}
	return result
}

func containsString(slice []string, s string) bool {
	return slices.Contains(slice, s)
}

func copySlice(src []Exercise) []Exercise {
	if src == nil {
		return []Exercise{}
	}
	result := make([]Exercise, len(src))
	copy(result, src)
	return result
}
