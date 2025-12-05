package parser

import (
	"testing"

	"github.com/jpollak/semantic-release-notes/git"
)

func TestParseCommit_ComplexExample(t *testing.T) {
	subject := "feat(config): Support Next Generation Model [SW-1928] (#1330)"
	body := `## 📝 Description

Add support for Next Generation Model via config changes.

## 📣 Customer-Facing Release Notes

We have added support for Next Generation Model via config changes.

## ⚙️ Configuration Changes

Be sure to enable the new model in the config.

## 🔌 Required Hardware Changes

<!-- List of Hardware Changes required or none if not required -->`

	commit := git.Commit{
		Subject: subject,
		Body:    body,
	}

	parsed := ParseCommit(commit)

	if parsed.Type != ChangeTypeFeature {
		t.Errorf("Expected Type to be Feature, got %v", parsed.Type)
	}
	if parsed.Component != "config" {
		t.Errorf("Expected Component to be 'config', got '%s'", parsed.Component)
	}
	if parsed.IssueNumber != "SW-1928" {
		t.Errorf("Expected IssueNumber to be 'SW-1928', got '%s'", parsed.IssueNumber)
	}
	// The description should match the subject description exactly, as we no longer append the body
	expectedDesc := "Support Next Generation Model"
	if parsed.Description != expectedDesc {
		t.Errorf("Expected Description to be '%s', got '%s'", expectedDesc, parsed.Description)
	}

	expectedCustomerNotes := "We have added support for Next Generation Model via config changes."
	if parsed.CustomerFacingNotes != expectedCustomerNotes {
		t.Errorf("Expected CustomerFacingNotes to be '%s', got '%s'", expectedCustomerNotes, parsed.CustomerFacingNotes)
	}

	expectedConfigChanges := "Be sure to enable the new model in the config."
	if parsed.ConfigurationChanges != expectedConfigChanges {
		t.Errorf("Expected ConfigurationChanges to be '%s', got '%s'", expectedConfigChanges, parsed.ConfigurationChanges)
	}

	if parsed.RequiredHardwareChanges != "" {
		t.Errorf("Expected RequiredHardwareChanges to be empty, got '%s'", parsed.RequiredHardwareChanges)
	}
}

func TestCleanContent(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Valid content", "Valid content"},
		{"  Valid content with spaces  ", "Valid content with spaces"},
		{"None", ""},
		{"none", ""},
		{"N/A", ""},
		{"n/a", ""},
		{"None required", ""},
		{"N/A - not applicable", ""},
		{"none.", ""},
		{"<!-- comment --> Real content", "Real content"},
		{"<!-- comment --> None", ""},
		{"", ""},
		{"   ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := cleanContent(tt.input); got != tt.expected {
				t.Errorf("cleanContent(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
