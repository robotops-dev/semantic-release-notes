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
}

var (
	// Example: Merge pull request #123 from user/feature/awesome-thing
	// We allow any type (feature, fix, chore, etc.)
	mergePRRegex = regexp.MustCompile(`Merge pull request #(\d+) from .*/([^/]+)/(.+)`)
	// Example: Merge branch 'feature/awesome-thing'
	mergeBranchRegex = regexp.MustCompile(`Merge branch '([^/]+)/(.+)'`)
)

func ParseCommit(commit git.Commit) ParsedCommit {
	parsed := ParsedCommit{
		Type:     ChangeTypeOther,
		Original: commit,
	}

	// Try to match PR merge pattern
	if matches := mergePRRegex.FindStringSubmatch(commit.Subject); len(matches) > 3 {
		parsed.PRNumber = matches[1]
		parsed.Type = mapType(matches[2])
		parsed.Description = matches[3] // Default description from branch name
	} else if matches := mergeBranchRegex.FindStringSubmatch(commit.Subject); len(matches) > 2 {
		parsed.Type = mapType(matches[1])
		parsed.Description = matches[2]
	}

	// Clean up description from branch name
	parsed.Description = strings.ReplaceAll(parsed.Description, "-", " ")
	parsed.Description = strings.Title(parsed.Description)

	// Parse body for sections
	if desc := extractSection(commit.Body, "## 📝 Description"); desc != "" {
		parsed.Description = desc
	}
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
	if lower == "none" || lower == "n/a" || lower == "" {
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
