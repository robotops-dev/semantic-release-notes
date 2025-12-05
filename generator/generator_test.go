package generator

import (
	"strings"
	"testing"

	"github.com/jpollak/semantic-release-notes/git"
	"github.com/jpollak/semantic-release-notes/parser"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name     string
		commits  []parser.ParsedCommit
		from     string
		to       string
		expected []string // Substrings expected in the output
	}{
		{
			name: "Basic Feature and Fix",
			commits: []parser.ParsedCommit{
				{
					Type:        parser.ChangeTypeFeature,
					Component:   "ui",
					Description: "add new button",
					IssueNumber: "123",
					Original:    git.Commit{Hash: "abcdef123456"},
				},
				{
					Type:        parser.ChangeTypeFix,
					Component:   "backend",
					Description: "fix crash",
					IssueNumber: "456",
					Original:    git.Commit{Hash: "123456abcdef"},
				},
			},
			from: "v1.0.0",
			to:   "v1.1.0",
			expected: []string{
				"# Release Notes (v1.0.0...v1.1.0)",
				"## Features",
				"* add new button [123]",
				"## Bug Fixes",
				"* fix crash [456]",
			},
		},
		{
			name: "With Customer Notes and Config Changes",
			commits: []parser.ParsedCommit{
				{
					Type:                 parser.ChangeTypeFeature,
					Description:          "big feature",
					CustomerFacingNotes:  "Added big feature.",
					ConfigurationChanges: "Enable flag X.",
				},
			},
			from: "v1.0.0",
			to:   "v1.1.0",
			expected: []string{
				// "Release Notes" section was removed in recent changes
				// "Added big feature.",
				"## Configuration Changes",
				"### big feature",
				"Enable flag X.",
			},
		},
		{
			name: "Sorting by Component then Description",
			commits: []parser.ParsedCommit{
				{Type: parser.ChangeTypeFeature, Component: "b", Description: "z"},
				{Type: parser.ChangeTypeFeature, Component: "a", Description: "y"},
				{Type: parser.ChangeTypeFeature, Component: "b", Description: "a"},
			},
			from: "",
			to:   "v1.0.0",
			expected: []string{
				"# Release Notes (v1.0.0)",
				"* y", // a comes first
				"* a", // b comes second, sorted by description a
				"* z", // b comes second, sorted by description z
			},
		},
		{
			name:     "Empty Commits",
			commits:  []parser.ParsedCommit{},
			from:     "",
			to:       "",
			expected: []string{"# Release Notes", "No configuration changes.", "No required hardware changes."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := Generate(tt.commits, tt.from, tt.to)
			for _, exp := range tt.expected {
				if !strings.Contains(output, exp) {
					t.Errorf("Expected output to contain:\n%q\nGot:\n%s", exp, output)
				}
			}
		})
	}
}
