# Totion 🧠

A terminal-based notes app for Go, styled after Notion/Obsidian — built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) v2, [Bubbles](https://github.com/charmbracelet/bubbles) v2, and [Lipgloss](https://github.com/charmbracelet/lipgloss) v2.

Create, edit, rename, delete, and preview Markdown notes without leaving your terminal.

## Features

- **Create, edit, save, delete, rename notes** — all stored as plain `.md` files
- **Markdown preview** — render notes with [Glamour](https://github.com/charmbracelet/glamour) instead of viewing raw text
- **List + fuzzy filter** — browse all notes and filter by name on the fly
- **Autosave** — debounced background saves while you type, on top of manual save
- **Config file** — customize vault location and theme via a JSON config
- **Full keybinding help** — press `?` for a complete list of shortcuts
- **Confirm-before-delete** — no accidental data loss

## Installation

```bash
git clone https://github.com/Ashoke15/notes-app.git
cd notes-app
make install
make run
```

`make install` downloads dependencies; `make run` starts the app directly. Run `make build` instead if you want a standalone `totion` binary, or `make test` to run the test suite.

## Usage

| Key | Action |
|---|---|
| `Ctrl+N` | Create a new note |
| `Ctrl+L` | List all notes |
| `Enter` | Open selected note / confirm new note name |
| `Ctrl+S` | Save the current note |
| `Ctrl+P` | Toggle Markdown preview |
| `r` | Rename the selected note |
| `d` | Delete the selected note (with confirmation) |
| `/` | Filter the note list |
| `Esc` | Go back / cancel |
| `?` | Toggle full help |
| `q` / `Ctrl+C` | Quit |

Notes are saved as autosave is on by default — edits are written to disk automatically a moment after you stop typing, in addition to manual `Ctrl+S` saves.

## Configuration

Totion works with **zero configuration** — by default it stores notes in `~/.totion` and follows your terminal's light/dark background automatically.

To customize, create an optional JSON config file at your platform's standard config directory (e.g. `~/.config/totion/config.json` on Linux):

```json
{
  "vault_dir": "~/Documents/notes",
  "theme": "dark"
}
```

| Field | Values | Description |
|---|---|---|
| `vault_dir` | any path, `~` expands to home | Where notes are stored. Defaults to `~/.totion` if omitted |
| `theme` | `"auto"`, `"dark"`, `"light"` | Defaults to `"auto"` if omitted |

Both fields are optional — omit either one (or the whole file) to use the default.

## Project structure

```
notes-app/
├── main.go                      # entrypoint — wires everything together
├── internal/
│   ├── config/
│   │   └── config.go             # load/save JSON config
│   ├── vault/
│   │   ├── vault.go               # note CRUD — pure file I/O, no TUI dependencies
│   │   └── vault_test.go
│   └── handeler/
│       └── ui/
│           ├── model.go           # Model struct, view-state machine
│           ├── initialize.go      # model construction, widget setup
│           ├── update.go          # Update() — all key handling
│           ├── view.go            # View() — rendering
│           ├── file_op.go         # note actions (create/open/save/delete/rename/preview)
│           ├── autosave.go        # debounced autosave via tea.Tick
│           ├── keys.go            # keybinding definitions (bubbles/help)
│           └── styles.go          # Lipgloss styles
```

`internal/vault` has zero Bubble Tea or Lipgloss dependencies by design — it's plain Go file I/O, fully unit-testable with `t.TempDir()` and independent of the terminal UI layer.

## Development

```bash
make test     # run the test suite
make build    # build a standalone binary
```

## Roadmap

Built iteratively through a structured beginner → intermediate → advanced progression:

- [x] Core CRUD (create, list, open, edit, save)
- [x] Delete with confirmation
- [x] Rename
- [x] List filtering
- [x] Config file (vault path, theme)
- [x] Markdown preview
- [x] Full help view
- [x] Debounced autosave
- [ ] YAML frontmatter + tags
- [ ] Git-backed vault (auto-commit on save)
- [ ] Live vault refresh via file watching
- [ ] Command palette (fuzzy action search)
- [ ] Split-pane list + preview layout
- [ ] Model-level tests (`teatest`)
- [ ] Packaged releases (`goreleaser` + CI)

## License

MIT
