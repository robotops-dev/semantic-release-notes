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
	// The description should include the subject description and the body description
	expectedDescStart := "Support Next Generation Model"
	if len(parsed.Description) < len(expectedDescStart) || parsed.Description[:len(expectedDescStart)] != expectedDescStart {
		t.Errorf("Expected Description to start with '%s', got '%s'", expectedDescStart, parsed.Description)
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
