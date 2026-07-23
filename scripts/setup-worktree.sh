#!/usr/bin/env bash
# Run inside a freshly created worktree (after `git worktree add`).
# Copies gitignored local files listed in .worktreeinclude from the
# main worktree, since git worktree add doesn't carry them over.
set -euo pipefail

current_root=$(git rev-parse --show-toplevel)
main_root=$(git worktree list --porcelain | awk '/^worktree /{print $2; exit}')

if [ "$main_root" = "$current_root" ]; then
  echo "already in main worktree, nothing to copy"
  exit 0
fi

include_file="$current_root/.worktreeinclude"
if [ ! -f "$include_file" ]; then
  echo "no .worktreeinclude found"
  exit 0
fi

while IFS= read -r path; do
  [ -z "$path" ] && continue
  case "$path" in \#*) continue ;; esac

  src="$main_root/$path"
  dest="$current_root/$path"
  if [ -f "$src" ]; then
    mkdir -p "$(dirname "$dest")"
    cp "$src" "$dest"
    echo "copied $path"
  fi
done < "$include_file"
