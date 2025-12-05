package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Commit struct {
	Hash    string
	Subject string
	Body    string
}

// GetMergeCommits returns a list of merge commits from the git log.
// It uses the "git log --merges" command.
func GetMergeCommits(repoPath, from, to string) ([]Commit, error) {
	// Format: Hash%nSubject%nBody%n---COMMIT-END---%n
	// Use --first-parent to follow the main branch history, avoiding intermediate commits from merged branches.
	// We also remove --merges to include direct commits to main if any, though typically PRs are merges.
	// Actually, the user asked to remove merge parsing support, implying we treat everything as a commit.
	// But usually we still want to see the merge commits themselves if they contain the release notes.
	// Let's stick to --first-parent.
	args := []string{"log", "--first-parent", "--pretty=format:%H%n%s%n%b%n---COMMIT-END---%n"}
	if from != "" || to != "" {
		rangeSpec := ""
		if from != "" && to != "" {
			rangeSpec = fmt.Sprintf("%s..%s", from, to)
		} else if from != "" {
			rangeSpec = fmt.Sprintf("%s..HEAD", from)
		} else if to != "" {
			rangeSpec = to
		}
		args = append(args, rangeSpec)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run git log: %w, stderr: %s", err, stderr.String())
	}

	return parseGitLogOutput(out.String()), nil
}

// FetchTags fetches all tags from the remote.
func FetchTags(repoPath string) error {
	cmd := exec.Command("git", "fetch", "--tags")
	cmd.Dir = repoPath
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch tags: %w, stderr: %s", err, stderr.String())
	}
	return nil
}

// TagExists checks if a tag or commit hash exists in the repository.
func TagExists(repoPath, tag string) (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--verify", tag)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return false, nil
	}
	return true, nil
}

func parseGitLogOutput(output string) []Commit {
	var commits []Commit
	// Split by the delimiter.
	entries := strings.Split(output, "\n---COMMIT-END---\n")

	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		lines := strings.SplitN(entry, "\n", 3)
		commit := Commit{}
		if len(lines) >= 1 {
			commit.Hash = lines[0]
		}
		if len(lines) >= 2 {
			commit.Subject = lines[1]
		}
		if len(lines) >= 3 {
			commit.Body = strings.TrimSpace(lines[2])
		}
		commits = append(commits, commit)
	}
	return commits
}
