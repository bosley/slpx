#!/usr/bin/env zsh

# WARNING: This file was auto-generated

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

TEST_DIRS=(
    "primitive"
    "complex"
)

echo ""
echo "${BLUE}${BOLD}╔═══════════════════════════════════════════════════╗${RESET}"
echo "${BLUE}${BOLD}║           SLPX Test Suite Runner                 ║${RESET}"
echo "${BLUE}${BOLD}╚═══════════════════════════════════════════════════╝${RESET}"
echo ""

TOTAL_START=$(gdate +%s.%N 2>/dev/null || date +%s)
PASSED=0
FAILED=0
FAILED_TESTS=()

for TEST_DIR in "${TEST_DIRS[@]}"; do
    if [[ ! -d "$TEST_DIR" ]]; then
        echo "${YELLOW}⚠️  Skipping ${TEST_DIR}: directory not found${RESET}"
        continue
    fi

    if [[ ! -f "$TEST_DIR/main.slpx" ]]; then
        echo "${YELLOW}⚠️  Skipping ${TEST_DIR}: no main.slpx found${RESET}"
        continue
    fi

    echo "${CYAN}${BOLD}📦 Running ${TEST_DIR} tests...${RESET}"
    echo ""

    START_TIME=$(gdate +%s.%N 2>/dev/null || date +%s)
    
    (cd "$TEST_DIR" && ../../build/slpx main.slpx)
    EXIT_CODE=$?
    
    END_TIME=$(gdate +%s.%N 2>/dev/null || date +%s)

    if command -v gdate &> /dev/null; then
        DURATION=$(echo "$END_TIME - $START_TIME" | bc)
        FORMATTED_TIME=$(printf "%.3f" $DURATION)
    else
        DURATION=$((END_TIME - START_TIME))
        FORMATTED_TIME="${DURATION}"
    fi

    echo ""
    if [ $EXIT_CODE -eq 0 ]; then
        echo "${GREEN}✅ ${TEST_DIR} passed (${FORMATTED_TIME}s)${RESET}"
        ((PASSED++))
    else
        echo "${RED}❌ ${TEST_DIR} failed with exit code ${EXIT_CODE} (${FORMATTED_TIME}s)${RESET}"
        ((FAILED++))
        FAILED_TESTS+=("$TEST_DIR")
    fi
    echo ""
done

TOTAL_END=$(gdate +%s.%N 2>/dev/null || date +%s)

if command -v gdate &> /dev/null; then
    TOTAL_DURATION=$(echo "$TOTAL_END - $TOTAL_START" | bc)
    TOTAL_FORMATTED=$(printf "%.3f" $TOTAL_DURATION)
else
    TOTAL_DURATION=$((TOTAL_END - TOTAL_START))
    TOTAL_FORMATTED="${TOTAL_DURATION}"
fi

echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo "${BOLD}📊 Test Summary${RESET}"
echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo "${GREEN}  ✅ Passed: ${PASSED}${RESET}"
echo "${RED}  ❌ Failed: ${FAILED}${RESET}"
echo "${BOLD}  ⏱️  Total Runtime: ${TOTAL_FORMATTED}s${RESET}"

if [ $FAILED -gt 0 ]; then
    echo ""
    echo "${RED}${BOLD}Failed test suites:${RESET}"
    for FAILED_TEST in "${FAILED_TESTS[@]}"; do
        echo "${RED}  • ${FAILED_TEST}${RESET}"
    done
fi

echo "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo "${GREEN}${BOLD}🎉 All test suites passed!${RESET}"
    echo ""
    exit 0
else
    exit 1
fi

