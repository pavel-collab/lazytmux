package config

import (
	"testing"
)

func TestGetCategories(t *testing.T) {
	categories := GetCategories()

	if len(categories) == 0 {
		t.Error("GetCategories() should return at least one category")
	}

	// Check that each category has required fields
	for _, cat := range categories {
		if cat.ID == "" {
			t.Error("category ID should not be empty")
		}
		if cat.NameEN == "" {
			t.Error("category NameEN should not be empty")
		}
		if cat.NameRU == "" {
			t.Error("category NameRU should not be empty")
		}
		if len(cat.Options) == 0 {
			t.Errorf("category %q should have at least one option", cat.ID)
		}
	}
}

func TestGetCategoriesExpectedCategories(t *testing.T) {
	categories := GetCategories()

	expectedIDs := []string{"general", "statusbar", "mouse", "visual", "keys", "windows"}

	categoryMap := make(map[string]bool)
	for _, cat := range categories {
		categoryMap[cat.ID] = true
	}

	for _, id := range expectedIDs {
		if !categoryMap[id] {
			t.Errorf("expected category %q to exist", id)
		}
	}
}

func TestGetAllOptions(t *testing.T) {
	options := GetAllOptions()

	if len(options) == 0 {
		t.Error("GetAllOptions() should return at least one option")
	}

	// Verify some expected options exist
	expectedOptions := []string{
		"mouse",
		"history-limit",
		"base-index",
		"status",
		"mode-keys",
		"escape-time",
	}

	for _, key := range expectedOptions {
		if _, ok := options[key]; !ok {
			t.Errorf("expected option %q to exist", key)
		}
	}
}

func TestGetOption(t *testing.T) {
	// Test existing option
	opt, ok := GetOption("mouse")
	if !ok {
		t.Error("GetOption('mouse') should return true")
	}
	if opt.Key != "mouse" {
		t.Errorf("opt.Key = %q, expected 'mouse'", opt.Key)
	}
	if opt.Type != TypeBool {
		t.Errorf("mouse option should be TypeBool, got %v", opt.Type)
	}

	// Test non-existent option
	_, ok = GetOption("nonexistent-option")
	if ok {
		t.Error("GetOption('nonexistent-option') should return false")
	}
}

func TestOptionHasRequiredFields(t *testing.T) {
	options := GetAllOptions()

	for key, opt := range options {
		t.Run(key, func(t *testing.T) {
			if opt.Key == "" {
				t.Error("option Key should not be empty")
			}
			if opt.DescEN == "" {
				t.Error("option DescEN should not be empty")
			}
			if opt.DescRU == "" {
				t.Error("option DescRU should not be empty")
			}
			if opt.Default == "" && opt.Type != TypeString {
				// TypeString can have empty default
				t.Error("option Default should not be empty for non-string types")
			}
		})
	}
}

func TestOptionChoicesForChoiceType(t *testing.T) {
	options := GetAllOptions()

	for key, opt := range options {
		if opt.Type == TypeChoice {
			if len(opt.Choices) == 0 {
				t.Errorf("option %q is TypeChoice but has no choices", key)
			}

			// Default should be one of the choices
			found := false
			for _, choice := range opt.Choices {
				if choice == opt.Default {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("option %q default %q is not in choices %v",
					key, opt.Default, opt.Choices)
			}
		}
	}
}

func TestOptionMinMaxForNumberType(t *testing.T) {
	options := GetAllOptions()

	for key, opt := range options {
		if opt.Type == TypeNumber {
			if opt.Max < opt.Min {
				t.Errorf("option %q has Max (%d) < Min (%d)", key, opt.Max, opt.Min)
			}
		}
	}
}

func TestOptionScopeIsValid(t *testing.T) {
	options := GetAllOptions()
	validScopes := map[Scope]bool{
		ScopeServer:  true,
		ScopeSession: true,
		ScopeWindow:  true,
	}

	for key, opt := range options {
		if !validScopes[opt.Scope] {
			t.Errorf("option %q has invalid scope %q", key, opt.Scope)
		}
	}
}

func TestSpecificOptionConfigurations(t *testing.T) {
	tests := []struct {
		key          string
		expectedType OptionType
		scope        Scope
	}{
		{"mouse", TypeBool, ScopeSession},
		{"history-limit", TypeNumber, ScopeSession},
		{"status-position", TypeChoice, ScopeSession},
		{"mode-keys", TypeChoice, ScopeSession},
		{"escape-time", TypeNumber, ScopeServer},
		{"automatic-rename", TypeBool, ScopeWindow},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			opt, ok := GetOption(tt.key)
			if !ok {
				t.Fatalf("option %q not found", tt.key)
			}

			if opt.Type != tt.expectedType {
				t.Errorf("type = %v, expected %v", opt.Type, tt.expectedType)
			}

			if opt.Scope != tt.scope {
				t.Errorf("scope = %q, expected %q", opt.Scope, tt.scope)
			}
		})
	}
}

func TestMouseOptionDefaults(t *testing.T) {
	opt, ok := GetOption("mouse")
	if !ok {
		t.Fatal("mouse option not found")
	}

	if opt.Default != "off" {
		t.Errorf("mouse default = %q, expected 'off'", opt.Default)
	}
}

func TestHistoryLimitRange(t *testing.T) {
	opt, ok := GetOption("history-limit")
	if !ok {
		t.Fatal("history-limit option not found")
	}

	if opt.Min != 0 {
		t.Errorf("history-limit Min = %d, expected 0", opt.Min)
	}

	if opt.Max < 10000 {
		t.Errorf("history-limit Max = %d, should be at least 10000", opt.Max)
	}
}

func TestStatusPositionChoices(t *testing.T) {
	opt, ok := GetOption("status-position")
	if !ok {
		t.Fatal("status-position option not found")
	}

	expectedChoices := []string{"top", "bottom"}
	if len(opt.Choices) != len(expectedChoices) {
		t.Errorf("status-position has %d choices, expected %d",
			len(opt.Choices), len(expectedChoices))
	}

	for _, expected := range expectedChoices {
		found := false
		for _, choice := range opt.Choices {
			if choice == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected choice %q not found in status-position", expected)
		}
	}
}

func TestModeKeysChoices(t *testing.T) {
	opt, ok := GetOption("mode-keys")
	if !ok {
		t.Fatal("mode-keys option not found")
	}

	expectedChoices := []string{"vi", "emacs"}
	for _, expected := range expectedChoices {
		found := false
		for _, choice := range opt.Choices {
			if choice == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected choice %q not found in mode-keys", expected)
		}
	}
}

func TestAllOptionsHaveBilingualDescriptions(t *testing.T) {
	options := GetAllOptions()

	for key, opt := range options {
		if opt.DescEN == "" {
			t.Errorf("option %q missing English description", key)
		}
		if opt.DescRU == "" {
			t.Errorf("option %q missing Russian description", key)
		}
	}
}

func TestCategoriesHaveBilingualNames(t *testing.T) {
	categories := GetCategories()

	for _, cat := range categories {
		if cat.NameEN == "" {
			t.Errorf("category %q missing English name", cat.ID)
		}
		if cat.NameRU == "" {
			t.Errorf("category %q missing Russian name", cat.ID)
		}
	}
}
