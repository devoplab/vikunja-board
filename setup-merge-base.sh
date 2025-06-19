#!/bin/bash
# Merge multiple branches into a disposable merged-base foundation branch

# Configure local Git settings for this repo
git config user.name "Your Name"
git config user.email "your.email@example.com"

# Refresh all remote branches
git fetch --all

# Create/Recreate merged-base branch
BRANCH_NAME="merged-base"
git checkout -B $BRANCH_NAME origin/main --no-track > /dev/null

# Ordered list of branches to merge (customize order as needed)
MERGE_BRANCHES=(
    origin/cosmetic-change
    origin/back-util
    origin/dept-user-setup
    # origin/attachment-summary
)

# Perform merges with conflict markers if needed
for branch in "${MERGE_BRANCHES[@]}"; do
    echo "┗ Merging: $branch"
    git merge --no-commit --no-ff "$branch" || {
        echo "⚠️  Conflict detected! Resolve conflicts then:"
        echo "git commit -m 'Merge $branch'" 
        exit 1
    }
    git commit -m "Integrated: $branch" --no-verify
done

echo -e "\n✅ Merged-base '$BRANCH_NAME' ready:"
echo "git push origin $BRANCH_NAME --force  # Only if sharing with team"