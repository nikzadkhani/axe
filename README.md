# Axe 🪓

A CLI tool to chop down Git branches that have been squash-merged on GitHub.

## Features

- 🔍 **Find** branches that are ready to axe (squash-merged but still local)
- 🪓 **Chop** merged branches with confirmation
- ⚡ **Fast** parallel GitHub API calls (10x faster than sequential)
- 🎨 **Pretty** colored output with progress spinners
- 🔒 **Safe** validates repository and confirms before deleting

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

### Find branches ready to axe

```bash
# List all squash-merged branches
axe branches

# Show PR details (verbose)
axe branches -v

# Show all branches with their PR status
axe branches --all

# Also works with aliases
axe list
axe ls

# Check a specific repository
axe branches -r /path/to/repo
```

### Chop down merged branches

```bash
# Chop with confirmation prompt
axe chop

# Dry run (see what would be chopped)
axe chop --dry-run

# Skip confirmation and chop immediately
axe chop --force

# Also works with aliases
axe clean
axe delete
axe rm

# Check a specific repository
axe chop -r /path/to/repo
```

### Disable colors (for CI/CD)

```bash
axe branches --no-color
axe chop --no-color --force
```

## How it works

1. Validates you're in a git repository
2. Fetches all local branches (excluding `main` and `master`)
3. Checks GitHub in parallel (10 workers) for merged PRs using `gh pr list`
4. Lists or chops branches that have been squash-merged

**Note:** Branches are force-deleted (`git branch -D`) because squash-merged commits have different SHAs than the original commits, so Git doesn't recognize them as merged.

## Examples

```bash
# See which branches are ready to axe
axe branches -v

# Preview what would be chopped
axe chop --dry-run

# Chop them down!
axe chop

# Quick chop without confirmation
axe chop -f
```

## Output Examples

### Default: List merged branches

```
⠋ Fetching local branches...
✓ Found 45 local branches
⠙ Looking for branches to chop (43 to check)...
✓ Found 5 branches ready to axe

🪓 Found 5 branch(es) to axe:
  feature/old-login
  bugfix/typo-fix
  refactor/cleanup
  ...
```

### Show all branch statuses

```bash
axe branches --all
```

```
⠋ Fetching local branches...
✓ Found 45 local branches
⠙ Checking PR status for 43 branches...
✓ Completed status check for 43 branches

🪓 Merged (ready to axe): 5 branch(es)
  feature/old-login (#123) Fix login flow
  bugfix/typo-fix (#124) Fix typo in header
  refactor/cleanup (#125) Code cleanup

📂 Open PR: 3 branch(es)
  feature/new-ui (#126) Add new dashboard UI
  feature/auth (#127) Implement OAuth

✏️ Draft PR: 2 branch(es)
  feature/experimental (#128) Testing new API

❌ Closed (not merged): 1 branch(es)
  feature/abandoned (#129) Abandoned work

🔍 No PR: 32 branch(es)
  temp/test-branch
  experimental/new-feature
  ...
```

### Chopping branches

```
🪓 Chop these branches? [y/N]: y
⠹ Chopping 5 branches...
✓ Chopped 5 branches

✓ Chopped: feature/old-login
✓ Chopped: bugfix/typo-fix
...

🪓 Chopped 5 branch(es)!
```
