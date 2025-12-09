package generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jpollak/semantic-release-notes/parser"
)

func Generate(commits []parser.ParsedCommit, from, to string) string {
	var sb strings.Builder

	header := "# Release Notes"
	if from != "" && to != "" {
		header += fmt.Sprintf(" (%s...%s)", from, to)
	} else if to != "" {
		header += fmt.Sprintf(" (%s)", to)
	}
	sb.WriteString(header + "\n\n")

	// Group Commits
	features := []parser.ParsedCommit{}
	fixes := []parser.ParsedCommit{}
	others := []parser.ParsedCommit{}

	for _, c := range commits {
		switch c.Type {
		case parser.ChangeTypeFeature:
			features = append(features, c)
		case parser.ChangeTypeFix:
			fixes = append(fixes, c)
		default:
			others = append(others, c)
		}
	}

	if len(features) > 0 {
		sb.WriteString("## Features\n")
		sb.WriteString(formatGroupedCommits(features))
		sb.WriteString("\n")
	}

	if len(fixes) > 0 {
		sb.WriteString("## Bug Fixes\n")
		sb.WriteString(formatGroupedCommits(fixes))
		sb.WriteString("\n")
	}

	if len(others) > 0 {
		sb.WriteString("## Other Changes\n")
		sb.WriteString(formatGroupedCommits(others))
		sb.WriteString("\n")
	}

	// 2. Configuration Changes
	sb.WriteString("## Configuration Changes\n\n")
	hasConfigChanges := false
	for _, c := range commits {
		if c.ConfigurationChanges != "" {
			sb.WriteString("### " + c.Description + "\n\n")
			sb.WriteString(c.ConfigurationChanges + "\n\n")
			hasConfigChanges = true
		}
	}
	if !hasConfigChanges {
		sb.WriteString("No configuration changes.\n\n")
	}

	// 3. Required Hardware Changes
	sb.WriteString("## Required Hardware Changes\n\n")
	hasHardwareChanges := false
	for _, c := range commits {
		if c.RequiredHardwareChanges != "" {
			sb.WriteString("### " + c.Description + "\n\n")
			sb.WriteString(c.RequiredHardwareChanges + "\n\n")
			hasHardwareChanges = true
		}
	}
	if !hasHardwareChanges {
		sb.WriteString("No required hardware changes.\n\n")
	}

	return sb.String()
}

func formatGroupedCommits(commits []parser.ParsedCommit) string {
	var sb strings.Builder

	// Sort commits by component, then description
	sort.Slice(commits, func(i, j int) bool {
		if commits[i].Component != commits[j].Component {
			if commits[i].Component == "other" {
				return false
			}
			if commits[j].Component == "other" {
				return true
			}
			return commits[i].Component < commits[j].Component
		}
		return commits[i].Description < commits[j].Description
	})

	currentComponent := ""
	for _, c := range commits {
		if c.Component != currentComponent {
			currentComponent = c.Component
			r := []rune(currentComponent)
			sb.WriteString("\n### " + strings.ToUpper(string(r[0])) + string(r[1:]) + "\n\n")
		}
		sb.WriteString(formatCommit(c))
	}

	return sb.String()
}

func formatCommit(c parser.ParsedCommit) string {
	if c.Description == "" {
		return ""
	}
	lines := strings.Split(c.Description, "\n")
	header := lines[0]
	body := ""
	if len(lines) > 1 {
		body = strings.Join(lines[1:], "\n")
	}

	var sb strings.Builder
	// Bullet point per change
	sb.WriteString(fmt.Sprintf("* %s", header))

	if c.IssueNumber != "" {
		sb.WriteString(fmt.Sprintf(" [%s]", c.IssueNumber))
	} else if c.PRNumber != "" {
		sb.WriteString(fmt.Sprintf(" (#%s)", c.PRNumber))
	} else if c.Original.Hash != "" {
		hash := c.Original.Hash
		if len(hash) > 8 {
			hash = hash[:8]
		}
		sb.WriteString(fmt.Sprintf(" (%s)", hash))
	}
	sb.WriteString("\n")
	if body != "" {
		sb.WriteString(strings.TrimSpace(body) + "\n")
	}

	return sb.String()
}
