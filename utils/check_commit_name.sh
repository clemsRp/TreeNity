#!/bin/sh

COMMIT_MSG="$1"

# Define pattern for authorized commit messages
PATTERN="^(feat|fix|chore|docs|refactor|test|ci|style|perf|build|revert)(\([a-zA-Z0-9_-]+\))?:? .+$"

# Check if commit message is valid
if printf '%s\n' "$COMMIT_MSG" | grep -qE "$PATTERN"; then
    # Return nothing if valid
    echo ""
else
    # Return error string if invalid
    echo "'$COMMIT_MSG' commit message does not respect required format (<type>(<scope>): <description> or <type> <description>)[cite: 1]."
    echo "Authorized types : feat, fix, chore, docs, refactor, test, ci, style, perf."
    exit 1
fi