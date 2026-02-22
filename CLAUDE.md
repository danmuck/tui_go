# CLAUDE.md

This file provides guidance to Claude Code when working with tui_go.

## Commands

```bash
go build ./...
go test ./...
go test -v ./...
go test -run TestName ./...
go vet ./...
go mod tidy
```

## Development

During development, go.mod uses a `replace` directive pointing to `../smplog`.
Before tagging a release, remove the replace directive and use a pinned version.

## Git

- Do not commit changes unless explicitly instructed to do so.
- Include feature size breakpoints in task lists to ask if I would like to commit the changes, giving me time to look them over before I commit them.
- By default, I make all commits.

## TUI Constructor

`NewTUI(w io.Writer)` requires an `io.Writer`. All rendered output flows through this writer.

## API Naming Convention

| Suffix | Category | Example |
|--------|----------|---------|
| `TC` | TUI Component (structured params) | `MenuTC`, `DividerTC` |
| `FU` | Flat Utility (compact single-line) | `FieldFU`, `KeyHintFU` |
| `TERM` | Terminal Control | `ClearScreenTERM`, `BeginFrameTERM` |

Package-level functions (`Clip`, `PadLeft`, `PadRight`, `Center`, config, `NewPhaseTimer`, `NewProgressBar`) have no suffix.

## Component Reference

Each TUI component lives in its own `component_*.go` file. All component tests live in `components_test.go`.

| Component | File | Key Types | Method |
|-----------|------|-----------|--------|
| Menu | `component_menu.go` | `MenuEntry`, `MenuParams`, `TitleParams` | `MenuTC`, `MenuTitleTC` |
| Selector | `component_selector.go` | `SelectorParams` | `SelectorTC` |
| Input | `component_input.go` | `InputParams` | `InputTC` |
| Divider | `component_divider.go` | `DividerParams` | `DividerTC` |
| TreeView | `component_tree.go` | `TreeNode`, `TreeViewEntry`, `TreeViewParams` | `TreeViewTC` |
| Summary | `component_summary.go` | `OperationSummaryParams`, `SummaryField`, `PhaseTimer` | `OperationSummaryTC` |
| Progress | `component_progress.go` | `ProgressBarParams`, `ProgressBar` | `NewProgressBar` |

Shared rendering helpers (`writeComponent`, `writeComposite`, `writeBlock`) are unexported TUI methods in `tui.go`. `effectiveWidth` is a package-level helper.

## Adding New Components

- Create `component_<name>.go` and add tests to `components_test.go`
- Update the Component Reference table above
- Use shared TUI methods from `tui.go` (`t.writeComponent`, `t.writeBlock`, `effectiveWidth`)

## Context Management

- Use `/compact` frequently during long sessions to reduce token surface and keep context focused.
- Use `/clear` when starting a new logical task or when prior context is no longer relevant.
- All agents (including subagents) should `/compact` after completing major steps.
