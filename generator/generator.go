package generator

import (
	"fmt"
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

	// 4. Detailed Commit Messages
	sb.WriteString("## Detailed Commit Messages\n\n")
	for _, c := range commits {
		sb.WriteString(formatCommit(c))
	}

	return sb.String()
}

func formatCommit(c parser.ParsedCommit) string {
	lines := strings.Split(c.Description, "\n")
	header := lines[0]
	body := ""
	if len(lines) > 1 {
		body = strings.Join(lines[1:], "\n")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### %s", header))
	if c.PRNumber != "" {
		sb.WriteString(fmt.Sprintf(" (#%s)", c.PRNumber))
	}
	sb.WriteString("\n\n")
	if body != "" {
		sb.WriteString(strings.TrimSpace(body) + "\n\n")
	}
	return sb.String()
}
