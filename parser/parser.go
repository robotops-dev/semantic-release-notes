package parser

import (
	"regexp"
	"strings"

	"github.com/jpollak/semantic-release-notes/git"
)

type ChangeType string

const (
	ChangeTypeFeature ChangeType = "Feature"
	ChangeTypeFix     ChangeType = "Fix"
	ChangeTypeOther   ChangeType = "Other"
)

type ParsedCommit struct {
	Type                    ChangeType
	Scope                   string
	Description             string
	PRNumber                string
	Original                git.Commit
	CustomerFacingNotes     string
	ConfigurationChanges    string
	RequiredHardwareChanges string
	Component               string
	IssueNumber             string
}

var (
	// Conventional commit regex: type(component): description [issue] (#pr)
	// Component is optional. Issue is optional. PR is optional.
	// We allow alphanumeric issues (e.g. SW-123).
	conventionalRegex = regexp.MustCompile(`^([a-z]+)(?:\((.+)\))?: (.+?)(?: \[#?([a-zA-Z0-9-]+)\])?(?: \(#\d+\))?$`)
)

func ParseCommit(commit git.Commit) ParsedCommit {
	parsed := ParsedCommit{
		Type:     ChangeTypeOther,
		Original: commit,
	}

	// Default description to subject
	parsed.Description = commit.Subject

	// Try to parse conventional commit from the subject
	if matches := conventionalRegex.FindStringSubmatch(commit.Subject); len(matches) > 3 {
		parsed.Type = mapType(matches[1])
		parsed.Component = matches[2]
		parsed.Description = matches[3]
		if len(matches) > 4 {
			parsed.IssueNumber = matches[4]
		}
	}

	// Parse body for sections
	parsed.CustomerFacingNotes = extractSection(commit.Body, "## 📣 Customer-Facing Release Notes")
	parsed.ConfigurationChanges = extractSection(commit.Body, "## ⚙️ Configuration Changes")
	parsed.RequiredHardwareChanges = extractSection(commit.Body, "## 🔌 Required Hardware Changes")

	return parsed
}

func extractSection(body, header string) string {
	if !strings.Contains(body, header) {
		return ""
	}
	parts := strings.Split(body, header)
	if len(parts) < 2 {
		return ""
	}
	// The content is after the header. We need to stop at the next header.
	content := parts[1]
	// Find the next header (starts with ##)
	if idx := strings.Index(content, "\n## "); idx != -1 {
		content = content[:idx]
	}
	return cleanContent(content)
}

func cleanContent(content string) string {
	// Remove XML comments
	xmlCommentRegex := regexp.MustCompile(`(?s)<!--.*?-->`)
	content = xmlCommentRegex.ReplaceAllString(content, "")

	// Trim whitespace
	content = strings.TrimSpace(content)

	// Check for "None" or "N/A" (case-insensitive)
	lower := strings.ToLower(content)
	if strings.HasPrefix(lower, "none") || strings.HasPrefix(lower, "n/a") || lower == "" {
		return ""
	}

	return content
}

func mapType(t string) ChangeType {
	switch strings.ToLower(t) {
	case "feature", "feat":
		return ChangeTypeFeature
	case "fix", "bugfix":
		return ChangeTypeFix
	default:
		return ChangeTypeOther
	}
}
