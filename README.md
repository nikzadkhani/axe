# Axe

A CLI tool to chop down Git branches that have been squash-merged on GitHub.

## Features

- **List** squash-merged branches that still exist locally
- **Clean** (delete) squash-merged branches with confirmation
- Uses GitHub CLI (`gh`) to check PR merge status

## Prerequisites

- Go 1.21 or later
- Git
- GitHub CLI (`gh`) installed and authenticated

## Installation

```bash
cd ~/Developer/axe
go install
```

This will install the `axe` binary to `$GOPATH/bin` (typically `~/go/bin`).

Make sure `$GOPATH/bin` is in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage

### List squash-merged branches

```bash
# In a git repository
axe list

# With verbose output (shows PR numbers and titles)
axe list -v

# In a specific repository
axe list -r /path/to/repo
```

### Clean squash-merged branches

```bash
# Delete squash-merged branches (with confirmation)
axe clean

# Dry run (see what would be deleted)
axe clean --dry-run

# Skip confirmation prompt
axe clean --force

# In a specific repository
axe clean -r /path/to/repo
```

## How it works

1. Gets all local branches (excluding `main` and `master`)
2. For each branch, uses `gh pr list` to check if there's a merged PR
3. Lists or deletes branches with merged PRs

**Note:** Branches are force-deleted (`git branch -D`) because squash-merged commits have different SHAs than the original commits, so Git doesn't recognize them as merged.

## Examples

```bash
# See which branches can be cleaned up
axe list -v

# Preview what would be deleted
axe clean --dry-run

# Actually clean up the branches
axe clean
```
