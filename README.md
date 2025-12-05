# Semantic Release Notes Generator

A command-line tool written in Go that generates structured release notes from
Git commit history. It is designed to work with **Conventional Commits** and
supports **First-Parent History** to cleanly generate release notes from merge
commits.

## Features

- **Conventional Commits Parsing**: Parses commit messages following the
  `type(component): description [issue]` format.
- **Structured Output**: Groups changes by type (Features, Bug Fixes, etc.) and
  then by component.
- **First-Parent History**: Uses `git log --first-parent` to traverse the main
  branch history, ideal for workflows where features are merged via PRs.
- **Markdown Generation**: Outputs clean, formatted Markdown release notes.
- **Custom Sections**: Extracts specific sections from commit bodies:
  - `## 📣 Customer-Facing Release Notes`
  - `## ⚙️ Configuration Changes`
  - `## 🔌 Required Hardware Changes`
- **Content Cleaning**: Automatically filters out XML comments and "None"/"N/A"
  placeholders.
- **Version Range Header**: Generates a header with the version range (e.g.,
  `# Release Notes (v1.0.0...v1.1.0)`).

## Usage

Build the tool:

```bash
go build -o semantic-release-notes
```

Run the tool by specifying the repository path and the tag range:

```bash
./semantic-release-notes -repo <path_to_repo> -from <start_tag> -to <end_tag>
```

### Arguments

- `-repo`: **(Required)** Path to the local git repository. Supports tilde
  expansion (e.g., `~/src/my-repo`).
- `-from`: **(Required)** The starting tag or commit hash (exclusive).
- `-to`: **(Required)** The ending tag or commit hash (inclusive).

### Example

```bash
./semantic-release-notes -repo ~/go/src/github.com/my-org/my-repo -from v1.0.0 -to v1.1.0
```

### Output Format

The tool outputs Markdown to stdout:

```markdown
# Release Notes (v1.0.0...v1.1.0)

## Release Notes

- Added a shiny new button.

## Configuration Changes

### big feature

Enable flag X.

## Features

- **add new button** [123]

## Bug Fixes

- **fix crash** [456]
```

## Development

### Prerequisites

- Go 1.20 or higher
- Git

### Running Tests

The project includes unit tests for the parser and generator, as well as an
integration verification script.

To run unit tests:

```bash
go test -v ./...
```

To run the integration verification script (requires a Unix-like environment
with bash):

```bash
./verify.sh
```

This script creates a temporary git repository, generates commits with various
scenarios (conventional commits, merge commits, tags), runs the tool, and
verifies the output.

### Project Structure

- `main.go`: Entry point, handles CLI arguments and orchestration.
- `git/`: Handles Git interactions (running `git log`, fetching tags).
- `parser/`: Parses raw commit messages into structured data.
- `generator/`: Generates the Markdown output from parsed commits.
