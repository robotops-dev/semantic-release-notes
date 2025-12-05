package generator

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jpollak/semantic-release-notes/parser"
)

func Generate(commits []parser.ParsedCommit) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Release Notes (%s)\n\n", time.Now().Format("2006-01-02")))

	// 1. Release Notes (Customer-Facing)
	hasReleaseNotes := false
	for _, c := range commits {
		if c.CustomerFacingNotes != "" {
			sb.WriteString(c.CustomerFacingNotes + "\n\n")
			hasReleaseNotes = true
		}
	}
	if !hasReleaseNotes {
		sb.WriteString("No release notes provided.\n\n")
	}

	// 2. Configuration Changes
	sb.WriteString("## Configuration Changes\n\n")
	hasConfigChanges := false
	for _, c := range commits {
		if c.ConfigurationChanges != "" {
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
			sb.WriteString(c.RequiredHardwareChanges + "\n\n")
			hasHardwareChanges = true
		}
	}
	if !hasHardwareChanges {
		sb.WriteString("No required hardware changes.\n\n")
	}

	// 4. Grouped Commits
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
		sb.WriteString("## Features\n\n")
		sb.WriteString(formatGroupedCommits(features))
	}

	if len(fixes) > 0 {
		sb.WriteString("## Bug Fixes\n\n")
		sb.WriteString(formatGroupedCommits(fixes))
	}

	if len(others) > 0 {
		sb.WriteString("## Other Changes\n\n")
		sb.WriteString(formatGroupedCommits(others))
	}

	return sb.String()
}

func formatGroupedCommits(commits []parser.ParsedCommit) string {
	var sb strings.Builder

	// Group by component
	byComponent := make(map[string][]parser.ParsedCommit)
	var noComponent []parser.ParsedCommit

	for _, c := range commits {
		if c.Component != "" {
			byComponent[c.Component] = append(byComponent[c.Component], c)
		} else {
			noComponent = append(noComponent, c)
		}
	}

	// Sort components for deterministic output
	var components []string
	for k := range byComponent {
		components = append(components, k)
	}
	sort.Strings(components)

	for _, comp := range components {
		sb.WriteString(fmt.Sprintf("### %s\n\n", comp))
		for _, c := range byComponent[comp] {
			sb.WriteString(formatCommit(c))
		}
	}

	if len(noComponent) > 0 {
		if len(components) > 0 {
			sb.WriteString("### General\n\n")
		}
		for _, c := range noComponent {
			sb.WriteString(formatCommit(c))
		}
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
	// Use bold for header instead of H3 since we are inside H2/H3 sections
	sb.WriteString(fmt.Sprintf("**%s**", header))

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
	sb.WriteString("\n\n")
	if body != "" {
		sb.WriteString(strings.TrimSpace(body) + "\n\n")
	}
	return sb.String()
}
