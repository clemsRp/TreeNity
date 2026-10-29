#!/bin/sh

BRANCH_NAME="$1"

# Define pattern for authorized names
PATTERN="^(feat|fix|chore|docs|refactor|test|ci|style|perf|build|revert)/[a-zA-Z0-9_-]+$"

# Check if branch name is valid
if printf '%s\n' "$BRANCH_NAME" | grep -qE "$PATTERN"; then
    # Return nothing if valid
    echo ""
else
    # Return error string if invalid
    echo "'$BRANCH_NAME' branch does not respect required format (<type>/<fonctionnality>)."
    echo "Authorized types : feat, fix, chore, docs, refactor, test, ci, style, perf."
    exit 1
fi