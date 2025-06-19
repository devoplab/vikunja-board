#!/bin/bash
set -euo pipefail  # Safer script execution: exit on error, unset variables, and pipe fails

## Build vikunja frontend and backend
DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

# Validate required commands
check_command() {
    if ! command -v "$1" &>/dev/null; then
        echo "Error: $1 command not found. Please install it before running this script."
        exit 1
    fi
}

check_command pnpm
check_command mage

echo "****************** Building frontend ******************"
(
    echo "→ Entering frontend directory"
    cd "${DIR}/frontend"
    
    echo "→ Installing dependencies"
    pnpm install
    
    echo "→ Building frontend"
    pnpm build
)
echo "✅ Frontend build completed successfully"
echo

echo "****************** Building backend ******************"
(
    echo "→ Entering project root"
    cd "${DIR}"
    
    echo "→ Building backend"
    mage build
)
echo "✅ Backend build completed successfully"