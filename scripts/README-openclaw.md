# OpenClaw repository-context launcher

From the repository root:

    chmod +x scripts/launch-openclaw-context.sh
    scripts/launch-openclaw-context.sh

Requirements:
- OpenClaw CLI installed and configured with a working provider.
- Git installed (the launcher uses git status to describe the current checkout).
- Repository changes still require the project's build workflow.

The script builds a temporary Markdown prompt from the root README, root CLAUDE.md when present, Git status, and a bounded directory listing. It launches openclaw agent exec --cwd <repo> --message-file <prompt>, so the agent's workspace is the current checkout and the prompt includes project conventions.

The prompt is stored in a temporary file and removed on exit. It does not read .env files, print environment variables, or automatically modify files. Review OpenClaw's tool policy before giving an agent write-capable access.

After code changes: Always run python3 build.py and include the diagnostic artifacts in your PR. Use .github/pull_request_template.md. The launcher has not been validated against a live OpenClaw installation in this draft.
