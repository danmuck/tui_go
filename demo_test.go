package tui

import (
	"bytes"
	"os"
	"testing"
	"time"

	smplog "github.com/danmuck/smplog"
)

// demoConfig loads tui.config.toml and applies it via Configure.
// On cleanup it restores DefaultConfig.
func demoConfig(t *testing.T) {
	t.Helper()
	cfg, err := ConfigFromFile("tui.config.toml")
	if err != nil {
		t.Fatalf("demoConfig: %v", err)
	}
	Configure(cfg)
	t.Cleanup(func() { Configure(DefaultConfig()) })
}

// TestVisualFlatHelpers renders all flat output helpers to stdout for visual
// inspection. Run with:
//
//	go test -v -run TestVisualFlatHelpers ./...
func TestVisualFlatHelpers(t *testing.T) {
	demoConfig(t)
	tui := NewTUI(os.Stdout)

	t.Log("── flat helpers ──")
	nl := func() { smplog.Fprintln(os.Stdout, "") }
	tui.DividerTC(&DividerParams{Rune: '=', Width: 48}); nl()
	tui.MenuItemFU(1, "selected item", true); nl()
	tui.MenuItemFU(2, "normal item", false); nl()
	tui.DividerTC(&DividerParams{Rune: '-', Width: 48}); nl()
	tui.FieldFU("host", "localhost:8080"); nl()
	tui.KeyHintFU("q", "quit"); nl()
	tui.KeyHintFU("r", "refresh"); nl()
	tui.DividerTC(&DividerParams{Rune: '-', Width: 48}); nl()
	tui.InputLineFU("search> ", "foo", true); nl()
	tui.InputLineFU("filter> ", "bar", false); nl()
	tui.DividerTC(&DividerParams{Rune: '-', Width: 48}); nl()
	tui.StatusInfoFU("everything is fine"); nl()
	tui.StatusWarnFU("something looks off"); nl()
	tui.StatusErrorFU("something broke"); nl()
	tui.DividerTC(&DividerParams{Rune: '=', Width: 48}); nl()
}

// TestVisualComponents renders three full scenes to stdout for visual inspection.
// Scene 1 uses hardcoded overrides; Scene 2 is a component reference showing all
// components and their signatures; Scene 3 uses tui.config.toml so it renders
// last and you can immediately see the effect of TOML tweaks.
// Run with:
//
//	go test -v -run TestVisualComponents ./...
func TestVisualComponents(t *testing.T) {
	tui := NewTUI(os.Stdout)

	// Scene 1: left-aligned, hardcoded palette
	t.Log("Scene 1: Left-aligned layout (hardcoded overrides)")
	Configure(Config{
		NoColor: false,
		Colors: ColorConfig{
			Title:   smplog.StyleColor256(11), // bright yellow
			Menu:    smplog.StyleColor256(12), // bright blue
			Prompt:  smplog.StyleColor256(13), // bright magenta
			Data:    smplog.StyleColor256(7),  // white
			Divider: smplog.StyleColor256(6),  // cyan
		},
		TUI: TUIConfig{
			MaxWidth:             0,
			Centered:             false,
			MenuSelectedPrefix:   "▶",
			MenuUnselectedPrefix: " ",
			MenuIndexWidth:       2,
			InputCursor:          "█",
			DividerWidth:         40,
		},
	})

	tui.DividerTC(&DividerParams{Rune: '='})
	tui.MenuTitleTC(&TitleParams{Text: "Settings"})
	tui.DividerTC(&DividerParams{Rune: '='})
	tui.MenuTC(&MenuParams{Items: []MenuEntry{
		{Label: "Network", Selected: false},
		{Label: "Storage", Selected: true},
		{Label: "Security", Selected: false},
	}})
	tui.DividerTC(&DividerParams{Rune: '='})
	tui.SelectorTC(&SelectorParams{
		Label:   "theme",
		Items:   []string{"dark", "light", "system"},
		Current: 0,
	})
	tui.InputTC(&InputParams{Label: "alias", Value: "dev-box", Active: true})
	tui.DividerTC(&DividerParams{Rune: '='})

	// Scene 2: Component reference — all components with their signatures
	t.Log("Scene 2: Component Reference (all components & signatures)")
	demoConfig(t)

	tui.DividerTC(&DividerParams{Rune: '='})
	tui.MenuTitleTC(&TitleParams{Text: "Component Reference"})
	tui.DividerTC(&DividerParams{Rune: '='})

	nl := func() { smplog.Fprintln(os.Stdout, "") }

	// TUI.MenuTitleTC(p *TitleParams)
	//   TitleParams{Text string, Width int}
	tui.FieldFU("MenuTitle", "TitleParams{Text, Width}"); nl()
	tui.MenuTitleTC(&TitleParams{Text: "Example Title"})
	tui.DividerTC(&DividerParams{})

	// TUI.MenuTC(p *MenuParams)
	//   MenuParams{Items []MenuEntry, Width int}
	//   MenuEntry{Label string, Selected bool}
	tui.FieldFU("Menu", "MenuParams{Items []MenuEntry, Width}"); nl()
	tui.FieldFU("MenuEntry", "{Label, Selected}"); nl()
	tui.MenuTC(&MenuParams{Items: []MenuEntry{
		{Label: "Selected item", Selected: true},
		{Label: "Normal item", Selected: false},
	}})
	tui.DividerTC(&DividerParams{})

	// TUI.SelectorTC(p *SelectorParams)
	//   SelectorParams{Label string, Items []string, Current int, Width int}
	tui.FieldFU("Selector", "SelectorParams{Label, Items []string, Current, Width}"); nl()
	tui.SelectorTC(&SelectorParams{
		Label:   "example",
		Items:   []string{"option-a", "option-b", "option-c"},
		Current: 1,
	})
	tui.DividerTC(&DividerParams{})

	// TUI.InputTC(p *InputParams)
	//   InputParams{Label string, Value string, Active bool, Width int}
	tui.FieldFU("Input", "InputParams{Label, Value, Active, Width}"); nl()
	tui.InputTC(&InputParams{Label: "active", Value: "typing", Active: true})
	tui.InputTC(&InputParams{Label: "inactive", Value: "static", Active: false})
	tui.DividerTC(&DividerParams{})

	// TUI.DividerTC(p *DividerParams)
	//   DividerParams{Rune rune, Width int}
	tui.FieldFU("Divider", "DividerParams{Rune, Width}"); nl()
	tui.DividerTC(&DividerParams{Rune: '~', Width: 20})
	tui.DividerTC(&DividerParams{})

	// Flat helpers (TUI methods with FU suffix)
	tui.MenuTitleTC(&TitleParams{Text: "Flat Helpers"})
	tui.DividerTC(&DividerParams{})

	tui.FieldFU("MenuItem", "MenuItemFU(index int, label string, selected bool)"); nl()
	tui.MenuItemFU(1, "example", true); nl()
	tui.MenuItemFU(2, "example", false); nl()

	tui.FieldFU("Field", "FieldFU(label string, value any)"); nl()
	tui.FieldFU("key", "value"); nl()

	tui.FieldFU("KeyHint", "KeyHintFU(key, desc string)"); nl()
	tui.KeyHintFU("q", "quit"); nl()

	tui.FieldFU("InputLine", "InputLineFU(prefix, value string, active bool)"); nl()
	tui.InputLineFU("prompt> ", "text", true); nl()
	tui.InputLineFU("prompt> ", "text", false); nl()

	tui.FieldFU("StatusInfo", "StatusInfoFU(msg string)"); nl()
	tui.StatusInfoFU("info message"); nl()
	tui.FieldFU("StatusWarn", "StatusWarnFU(msg string)"); nl()
	tui.StatusWarnFU("warning message"); nl()
	tui.FieldFU("StatusError", "StatusErrorFU(msg string)"); nl()
	tui.StatusErrorFU("error message"); nl()

	tui.FieldFU("Divider", "DividerTC(&DividerParams{Width: 20})"); nl()
	tui.DividerTC(&DividerParams{Width: 20})
	tui.FieldFU("DividerRune", "DividerTC(&DividerParams{Width: 20, Rune: '*'})"); nl()
	tui.DividerTC(&DividerParams{Width: 20, Rune: '*'})

	tui.DividerTC(&DividerParams{Rune: '='})

	// Scene 3: config-driven (from tui.config.toml) — rendered last
	t.Log("Scene 3: Config-driven layout (from tui.config.toml)")
	demoConfig(t)

	tui.DividerTC(&DividerParams{Rune: '='})
	tui.MenuTitleTC(&TitleParams{Text: "Main Menu"})
	tui.DividerTC(&DividerParams{})
	tui.MenuTC(&MenuParams{Items: []MenuEntry{
		{Label: "Status", Selected: true},
		{Label: "Settings", Selected: false},
		{Label: "Logs", Selected: false},
		{Label: "Quit", Selected: false},
	}})
	tui.SelectorTC(&SelectorParams{
		Label:   "mode",
		Items:   []string{"debug", "verbose", "silent"},
		Current: 1,
	})
	tui.InputTC(&InputParams{Label: "filter", Value: "error", Active: false})
	tui.InputTC(&InputParams{Label: "output", Value: "stdout", Active: true})
	tui.DividerTC(&DividerParams{Rune: '='})
}

// TestVisualTreeView renders a tree view for visual inspection.
// Run with:
//
//	go test -v -run TestVisualTreeView ./...
func TestVisualTreeView(t *testing.T) {
	demoConfig(t)
	tui := NewTUI(os.Stdout)

	tui.DividerTC(&DividerParams{Rune: '='})
	tui.MenuTitleTC(&TitleParams{Text: "TreeView"})
	tui.DividerTC(&DividerParams{})

	nodes := []TreeNode{
		testNode{key: "src", label: "src/", parent: ""},
		testNode{key: "src/cmd", label: "cmd/", parent: "src"},
		testNode{key: "src/cmd/main.go", label: "main.go", parent: "src/cmd"},
		testNode{key: "src/lib", label: "lib/", parent: "src"},
		testNode{key: "src/lib/util.go", label: "util.go", parent: "src/lib"},
		testNode{key: "src/lib/tree.go", label: "tree.go", parent: "src/lib"},
		testNode{key: "docs", label: "docs/", parent: ""},
		testNode{key: "docs/readme.md", label: "readme.md", parent: "docs"},
	}

	tui.FieldFU("TreeView", "TreeViewParams{Nodes []TreeNode, Width, ShowIndex}"); smplog.Fprintln(os.Stdout, "")
	tui.TreeViewTC(&TreeViewParams{Nodes: nodes})
	tui.DividerTC(&DividerParams{})

	tui.FieldFU("TreeView (ShowIndex)", "ShowIndex: true"); smplog.Fprintln(os.Stdout, "")
	tui.TreeViewTC(&TreeViewParams{Nodes: nodes, ShowIndex: true})
	tui.DividerTC(&DividerParams{Rune: '='})
}

// TestVisualProgressBar renders a progress bar for visual inspection.
// Run with:
//
//	go test -v -run TestVisualProgressBar ./...
func TestVisualProgressBar(t *testing.T) {
	demoConfig(t)
	tui := NewTUI(os.Stdout)

	tui.DividerTC(&DividerParams{Rune: '='})
	tui.MenuTitleTC(&TitleParams{Text: "ProgressBar"})
	tui.DividerTC(&DividerParams{})

	tui.FieldFU("ProgressBar", "NewProgressBar(dst, ProgressBarParams{...})"); smplog.Fprintln(os.Stdout, "")
	pb := NewProgressBar(&bytes.Buffer{}, ProgressBarParams{
		Label:     "upload",
		Total:     1024 * 1024,
		Width:     30,
		MinRender: 0,
		Out:       os.Stdout,
	})
	// Simulate progress
	chunk := make([]byte, 256*1024)
	for i := 0; i < 4; i++ {
		pb.Write(chunk)
	}
	pb.Done()

	tui.DividerTC(&DividerParams{Rune: '='})
}

// TestVisualOperationSummary renders operation summaries for visual inspection.
// Run with:
//
//	go test -v -run TestVisualOperationSummary ./...
func TestVisualOperationSummary(t *testing.T) {
	demoConfig(t)
	tui := NewTUI(os.Stdout)

	tui.DividerTC(&DividerParams{Rune: '='})
	tui.MenuTitleTC(&TitleParams{Text: "OperationSummary"})
	tui.DividerTC(&DividerParams{})

	// OK case
	tui.FieldFU("OperationSummary", "OK=true"); smplog.Fprintln(os.Stdout, "")
	pt := NewPhaseTimer()
	pt.Begin("connect")
	time.Sleep(10 * time.Millisecond)
	pt.Begin("transfer")
	time.Sleep(15 * time.Millisecond)
	pt.Begin("verify")
	time.Sleep(5 * time.Millisecond)
	pt.End()

	tui.OperationSummaryTC(&OperationSummaryParams{
		Title: "File Upload",
		OK:    true,
		Fields: []SummaryField{
			{Label: "Target", Value: "s3://bucket/key"},
			{Label: "Size", Value: "10.5 MiB"},
			{Label: "Files", Value: "42"},
		},
		Timer: pt,
	})

	tui.DividerTC(&DividerParams{})

	// FAILED case
	tui.FieldFU("OperationSummary", "OK=false"); smplog.Fprintln(os.Stdout, "")
	tui.OperationSummaryTC(&OperationSummaryParams{
		Title: "Database Migration",
		OK:    false,
		Fields: []SummaryField{
			{Label: "Schema", Value: "v3 → v4"},
			{Label: "Error", Value: "foreign key constraint"},
		},
	})

	tui.DividerTC(&DividerParams{Rune: '='})
}
