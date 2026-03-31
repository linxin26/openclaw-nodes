---
name: verify
description: Run tests and lint to verify code changes
---

# Verify Skill

Run this skill before marking work complete or creating a PR.

## Commands

```bash
# Run all tests
go test ./...

# Build to check for compile errors
go build -o openclaw-node.exe ./cmd

# Run linter (golangci-lint v1.48.0 at ~/go/bin/golangci-lint.exe)
~/go/bin/golangci-lint.exe run --no-config --disable-all -E errcheck ./...
```

## When to Use

- Before marking a task as complete
- Before creating a commit
- After making significant changes
- When fixing bugs or implementing features

## Workflow

1. Run `go test ./...` to ensure all tests pass
2. Run `go build` to verify compilation
3. Run `golangci-lint` to check for issues
4. Fix any issues found before proceeding
