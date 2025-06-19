#!/bin/bash
# Creates/Updates feature branches based on merged-base foundation

# Check for required arguments
if [ $# -eq 0 ]; then
    echo "Usage: $0 FEATURE_NAME [MERGE_BASE]"
    echo "Example: $0 server-push"
    echo "         $0 mobile-api merged-base-v2"
    exit 1
fi

FEATURE_NAME="$1"
MERGE_BASE="${2:-merged-base}"  # Default to merged-base if not specified

echo "⌛ Creating feature branch '$FEATURE_NAME' from merge base '$MERGE_BASE'"

# Validate merge base exists
if ! git show-ref --verify refs/heads/$MERGE_BASE > /dev/null; then
    echo "❌ Merge base branch '$MERGE_BASE' does not exist!"
    echo "   First create it using setup-merge-base.sh"
    exit 1
fi

# Branch operations
if git show-ref --quiet refs/heads/$FEATURE_NAME; then
    echo "♻️  Updating existing feature branch..."
    git checkout $FEATURE_NAME
    git reset --hard HEAD
else
    echo "🆕 Creating new feature branch..."
    git checkout -b $FEATURE_NAME main --no-track
fi

# Apply merge base changes
echo "🔧 Applying merge base changes..."
git merge --squash $MERGE_BASE
git commit -m "BASE: Integrated changes from $MERGE_BASE" --no-verify

echo -e "\n✅ Feature branch '$FEATURE_NAME' ready:"
echo "git push -u origin $FEATURE_NAME  # Push to remote when ready"
echo "To update later:"
echo "$0 '$FEATURE_NAME' '$MERGE_BASE'"
