#!/usr/bin/env bash
# Launch OpenClaw with a bounded, current snapshot of this repository's contributor context.
set -euo pipefail

ROOT="$(pwd)"
if [[ ! -f "$ROOT/build.py" || ! -f "$ROOT/README.md" ]]; then
  printf 'error: run this script from the repository root\n' >&2
  exit 2
fi
if ! command -v openclaw >/dev/null 2>&1; then
  printf 'error: OpenClaw CLI not found. Install/configure it first.\n' >&2
  exit 127
fi

PROMPT_FILE="$(mktemp /tmp/tentoftrials-openclaw.XXXXXX.md)"
cleanup() { rm -f -- "$PROMPT_FILE"; }
trap cleanup EXIT INT TERM

{
  cat <<'HEADER'
You are working in the Tent of Trials repository. Use this context as orientation, not as a substitute for inspecting the relevant source files.

## Project purpose and repository map
HEADER
  sed -n '1,180p' README.md
  printf '\n## Agent instructions\n'
  if [[ -f CLAUDE.md ]]; then
    cat CLAUDE.md
  else
    printf 'No root CLAUDE.md found. Read the module build files before making changes.\n'
  fi
  printf '\n## Repository status\n'
  git status --short 2>/dev/null || true
  printf '\n## Top-level directories\n'
  find . -maxdepth 2 -type d \
    -not -path './.git*' \
    -not -path './node_modules*' \
    -not -path './target*' \
    -not -path './.venv*' \
    | sort | head -n 100
  cat <<'FOOTER'

## Mandatory workflow
- Read the relevant module's build file and tests before editing.
- Always run python3 build.py after changes and include the diagnostic artifacts in your PR.
- Use .github/pull_request_template.md for any PR.
- Never claim a build or test passed unless it actually ran; call out missing toolchains.
- Keep changes scoped and never weaken financial, security, or risk controls to make tests pass.
- Do not expose secrets from environment files, generated diagnostics, or logs.
FOOTER
} > "$PROMPT_FILE"

exec openclaw agent exec --cwd "$ROOT" --message-file "$PROMPT_FILE" \
  "Inspect the repository context above, ask clarifying questions if the task is ambiguous, and propose a small evidence-backed next step."
