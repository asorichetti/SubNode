#!/bin/bash
set -e

# ============================================================================
# SubNode Quick Smoke Test
# ============================================================================
# Usage: ./tools/ui-cli/scripts/smoke.sh
# 
# Runs a fast subset of journeys to verify the application is working.
# Designed for pre-commit hooks and quick local validation.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")/../.."
cd "$PROJECT_ROOT"

BINARY="/tmp/ui-cli-smoke-$$"
TIMEOUT=60

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

cleanup() {
    rm -f "$BINARY"
}
trap cleanup EXIT

echo "🧪 SubNode Smoke Test"
echo "===================="
echo ""

# Build the binary
echo "📦 Building UI-CLI..."
cd tools/ui-cli
go build -o "$BINARY" . 2>&1
cd "$PROJECT_ROOT"

# Check if backend is running
echo "🔍 Checking backend..."
if curl -sf http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Backend is running"
else
    echo -e "${RED}✗${NC} Backend not running at localhost:8080"
    echo "   Start with: make dev-backend"
    exit 1
fi

# Check if frontend is running
echo "🔍 Checking frontend..."
if curl -sf http://localhost:3000/ > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Frontend is running"
else
    echo -e "${RED}✗${NC} Frontend not running at localhost:3000"
    echo "   Start with: make dev-frontend"
    exit 1
fi

echo ""
echo "🚀 Running smoke journeys..."
echo ""

# Run a quick smoke test (core + auth only)
SMoke_OUTPUT=$("$BINARY" journey --suite core --suite auth --ux-gates off 2>&1)
SMOKE_EXIT=$?

echo "$SMoke_OUTPUT" | grep -E "(Running journey|✓ Journey|Error)" || true

if [ $SMOKE_EXIT -ne 0 ]; then
    echo ""
    echo -e "${RED}❌ Smoke test failed${NC}"
    exit 1
fi

PASS_COUNT=$(echo "$SMoke_OUTPUT" | grep -c "✓ Journey.*passed" || echo "0")
echo ""
echo -e "${GREEN}✓${NC} $PASS_COUNT/$(echo "$SMoke_OUTPUT" | grep -c "Running journey" || echo "0") smoke journeys passed"
echo ""
echo "✅ Application is healthy!"
echo ""
echo "Full test suite: cd tools/ui-cli && go run . journey --ux-gates warn"
