# Contributing to Totion

Thanks for your interest in contributing to Totion — a terminal-based notes app built in Go with Bubble Tea. This document covers everything you need to get set up, make a change, and submit it for review.

## Table of contents

- [Code of conduct](#code-of-conduct)
- [Getting started](#getting-started)
- [Project structure](#project-structure)
- [Development workflow](#development-workflow)
- [Coding guidelines](#coding-guidelines)
- [Testing](#testing)
- [Commit messages](#commit-messages)
- [Submitting a pull request](#submitting-a-pull-request)
- [Reporting bugs](#reporting-bugs)
- [Suggesting features](#suggesting-features)
- [Questions](#questions)

## Code of conduct

Be respectful, be constructive, and assume good faith. Disagreements about code or design are fine and expected — personal attacks, harassment, or dismissiveness toward other contributors are not. Maintainers may edit, remove, or reject contributions that don't meet this standard.

## Getting started

### Prerequisites

- Go 1.26.5 or later (check `go.mod` for the exact version this project targets)
- `make`
- Git

### Setup

```bash
git clone https://github.com/Ashoke15/notes-app.git
cd notes-app
make install
make run
```

`make install` downloads all dependencies. `make run` starts the app directly without producing a binary. See the [README](README.md) for the full keybinding reference and configuration options.

## Project structure

```
notes-app/
├── main.go                       # entrypoint — wires everything together
├── internal/
│   ├── config/
│   │   └── config.go               # load/save JSON config
│   ├── vault/
│   │   ├── vault.go                 # note CRUD — pure file I/O, no TUI dependencies
│   │   └── vault_test.go
│   └── handeler/
│       └── ui/
│           ├── model.go             # Model struct, view-state machine
│           ├── initialize.go        # model construction, widget setup
│           ├── update.go            # Update() — all key handling
│           ├── view.go              # View() — rendering
│           ├── file_op.go           # note actions (create/open/save/delete/rename/preview)
│           ├── autosave.go          # debounced autosave via tea.Tick
│           ├── keys.go              # keybinding definitions
│           └── styles.go            # Lipgloss styles
```

Two things worth understanding before you touch code:

- **`internal/vault` has zero Bubble Tea or Lipgloss imports, on purpose.** It's plain Go file I/O and should stay that way — this is what makes it independently unit-testable with `t.TempDir()`. If you're adding a feature that touches the filesystem, the logic belongs in `vault`, not in `ui`.
- **`internal/handeler/ui` is a finite state machine.** `Model.State` (a `viewstate`) determines what's on screen and what keys mean at any given moment. Nearly every bug in this codebase so far has come from a key handler or `View()` branch not accounting for the current state correctly — when adding a new screen or key binding, trace through every state it could fire from, not just the "happy path" state you're building it for.

## Development workflow

1. **Fork the repository** and clone your fork locally.
2. **Create a branch** off `main` for your change:
   ```bash
   git checkout -b feat/short-description
   ```
   Use a prefix that describes the kind of change: `feat/`, `fix/`, `docs/`, `refactor/`, `test/`.
3. **Make your change**, following the guidelines below.
4. **Run the test suite** before opening a PR (see [Testing](#testing)).
5. **Push your branch** and open a pull request against `main`.

## Coding guidelines

- **Format everything with `gofmt`** (or `go fmt ./...`) before committing. PRs with unformatted code will be asked to fix formatting before review.
- **Run `go vet ./...`** and resolve anything it flags.
- **Keep `vault` free of UI dependencies.** No `charm.land/...` imports in `internal/vault`.
- **Errors are returned, not swallowed.** Don't reintroduce `log.Fatal`/`log.Fatalf` in code paths that run after the app has started (`vault` package, key handlers) — a bad file, a missing note, or a failed write should surface as a `StatusMsg` in the UI, not crash the whole program. `log.Fatal` is only acceptable in one-time startup code in `main.go`/`init()`, where there's genuinely no way to continue.
- **Every state transition needs a guard.** If you add a new key binding, check what happens when it's pressed from every other `viewstate` — not just the one you intended it for. Look at how `d`/`r`/`y`/`n` are guarded against firing while the list is in `list.Filtering` state as a reference pattern.
- **Match existing naming and style** in the file you're editing rather than introducing a new convention. If you spot an inconsistency (e.g. a typo in an existing function name) and want to fix it, do that as a separate, clearly-labeled PR rather than folding it into a feature change — it makes both easier to review.

## Testing

```bash
make test
```

This runs `go test -v ./...` across the whole module.

- **New `vault` functions need tests.** Table-driven tests against `t.TempDir()` — see `vault_test.go` for the existing pattern (`TestCreatAndList`, `TestRename`, etc.) and follow the same style: one focused test per behavior, including at least one failure-path test (e.g. "already exists", "not found") where relevant.
- **UI-layer changes** don't currently have automated tests (model-level testing via `teatest` is a planned addition — see the README roadmap). For now, describe how you manually verified the change in your PR description, including which `viewstate` transitions you exercised.
- A PR that changes `vault` behavior without a corresponding test change will likely be asked to add one before merge.

## Commit messages

Keep commits focused and messages descriptive. Prefer:

```
Add rename support with textinput-based dialog
```

over:

```
updates
```

If a commit fixes a specific, previously-discussed bug (e.g. a file descriptor leak on Esc), say so in the message — it makes `git log` and `git blame` actually useful later.

## Submitting a pull request

- **Keep PRs scoped to one change.** A new feature and an unrelated bug fix should be two PRs, not one — much easier to review and revert independently if needed.
- **Describe what changed and why**, not just what files changed. If the PR fixes a bug, describe the reproduction steps. If it adds a feature, describe the new keybinding/behavior and any new states or config fields it introduces.
- **Note any new dependencies.** If your change adds an entry to `go.mod`, explain why it's needed and confirm `go.sum` is up to date (`go mod tidy`).
- **Update the README** if your change affects installation, usage, keybindings, or configuration — docs and code should land in the same PR, not as a follow-up.
- Be responsive to review feedback. It's normal for a PR to go through a few rounds before merge, especially around state-machine edge cases.

## Reporting bugs

Open an issue and include:

- **Steps to reproduce**, as specific as possible (which keys, in what order, from which screen)
- **What you expected to happen** vs. **what actually happened**
- **Your OS and terminal emulator** — TUI rendering issues are frequently terminal-specific
- **A screenshot or terminal recording**, if the bug is visual (sizing, styling, layout)
- Go version (`go version`) and Totion commit/version if relevant

If you're not sure whether something is a bug or intended behavior, open the issue anyway — worst case it gets relabeled.

## Suggesting features

Open an issue describing:

- The problem you're trying to solve (not just the feature itself — context helps evaluate the right shape for it)
- Any relevant prior art (how similar apps handle it)
- Whether you're interested in implementing it yourself

Check the README's roadmap section first — your idea may already be planned (frontmatter/tags, git-backed vault, file watching, command palette, split-pane layout, and packaged releases are all on there).

## Questions

If anything in this guide is unclear, or you're unsure whether a change fits the project's direction before you invest time in it, open an issue to ask first — that's always welcome and preferred over a large PR that turns out to be off-direction.

Thanks again for contributing to Totion!
