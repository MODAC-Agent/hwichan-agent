#!/usr/bin/env bash
set -euo pipefail

BEFORE_REF="${1:-975a30b}"
AFTER_REF="${2:-689ba3c}"
FILE_PATH="${3:-week2/Skill.md}"

cd "$(git rev-parse --show-toplevel)"

JQ_PICK_MAIN_INPUT='[.stats.models[]? | .roles.main?.tokens.input? // empty] | first // empty'

count_tokens() {
  local ref="$1" path="$2" content json tokens
  content=$(git show "${ref}:${path}")
  json=$(printf '%s' "$content" | gemini --approval-mode plan -p "Reply with only the word OK. Do not call any tools." -o json)
  tokens=$(printf '%s' "$json" | jq -r "$JQ_PICK_MAIN_INPUT")
  if [[ -z "$tokens" ]]; then
    echo "ERROR: main-role token count not found in gemini response for ${ref}:${path}" >&2
    echo "--- raw response ---" >&2
    printf '%s\n' "$json" >&2
    exit 1
  fi
  printf '%s' "$tokens"
}

before=$(count_tokens "$BEFORE_REF" "$FILE_PATH")
after=$(count_tokens "$AFTER_REF" "$FILE_PATH")
diff=$(( after - before ))
pct=$(awk -v d="$diff" -v b="$before" 'BEGIN { if (b==0) print "n/a"; else printf "%.2f%%", (d/b)*100 }')

printf 'file   : %s\n' "$FILE_PATH"
printf 'metric : stats.models[].roles.main.tokens.input (non-cached user input)\n'
printf 'before : %s  (%s tokens)\n' "$BEFORE_REF" "$before"
printf 'after  : %s  (%s tokens)\n' "$AFTER_REF" "$after"
printf 'diff   : %+d tokens (%s)\n' "$diff" "$pct"
