package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jpollak/semantic-release-notes/generator"
	"github.com/jpollak/semantic-release-notes/git"
	"github.com/jpollak/semantic-release-notes/parser"
)

func main() {
	repoPath := flag.String("repo", ".", "Path to the git repository")
	fromTag := flag.String("from", "", "Git tag/commit to start from (older)")
	toTag := flag.String("to", "", "Git tag/commit to end at (newer)")
	flag.Parse()

	repo := expandPath(*repoPath)
	fmt.Printf("Using repo path: %s\n", repo)

	// Fetch tags to ensure we have the latest
	if err := git.FetchTags(repo); err != nil {
		// Log warning but continue, as we might be offline or in a state where fetch fails
		// but the tags might already exist locally.
		fmt.Printf("Warning: failed to fetch tags: %v\n", err)
	}

	// Validate tags if provided
	if *fromTag != "" {
		exists, err := git.TagExists(repo, *fromTag)
		if err != nil {
			log.Fatalf("Error checking if tag %s exists: %v", *fromTag, err)
		}
		if !exists {
			log.Fatalf("Error: Tag or commit '%s' does not exist in the repository.", *fromTag)
		}
	}
	if *toTag != "" {
		exists, err := git.TagExists(repo, *toTag)
		if err != nil {
			log.Fatalf("Error checking if tag %s exists: %v", *toTag, err)
		}
		if !exists {
			log.Fatalf("Error: Tag or commit '%s' does not exist in the repository.", *toTag)
		}
	}

	commits, err := git.GetMergeCommits(repo, *fromTag, *toTag)
	if err != nil {
		log.Fatalf("Error getting merge commits: %v", err)
	}

	if len(commits) == 0 {
		fmt.Println("No merge commits found.")
		os.Exit(0)
	}

	var parsedCommits []parser.ParsedCommit
	for _, c := range commits {
		parsedCommits = append(parsedCommits, parser.ParseCommit(c))
	}

	notes := generator.Generate(parsedCommits, *fromTag, *toTag)
	fmt.Println(notes)
}

func expandPath(path string) string {
	if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return home
	}
	if len(path) > 1 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return home + path[1:]
	}
	return path
}
