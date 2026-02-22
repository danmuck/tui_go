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

## Component Reference

Each TUI component lives in its own `component_*.go` file. All component tests live in `components_test.go`.

| Component | File | Key Types |
|-----------|------|-----------|
| Menu | `component_menu.go` | `MenuEntry`, `MenuParams`, `TitleParams` |
| Selector | `component_selector.go` | `SelectorParams` |
| Input | `component_input.go` | `InputParams` |
| Divider | `component_divider.go` | `DividerParams` |
| TreeView | `component_tree.go` | `TreeNode`, `TreeViewEntry`, `TreeViewParams` |
| Summary | `component_summary.go` | `OperationSummaryParams`, `SummaryField`, `PhaseTimer` |
| Progress | `component_progress.go` | `ProgressBarParams`, `ProgressBar` |

Shared rendering helpers (`writeComponent`, `writeComposite`, `writeBlock`, `effectiveWidth`) live in `tui.go`.

## Adding New Components

- Create `component_<name>.go` and add tests to `components_test.go`
- Update the Component Reference table above
- Use shared helpers from `tui.go` (`writeComponent`, `writeBlock`, `effectiveWidth`)

## Context Management

- Use `/compact` frequently during long sessions to reduce token surface and keep context focused.
- Use `/clear` when starting a new logical task or when prior context is no longer relevant.
- All agents (including subagents) should `/compact` after completing major steps.
