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

	t.Log("── flat helpers ──")
	DividerRune(48, '=')
	MenuItem(1, "selected item", true)
	MenuItem(2, "normal item", false)
	DividerRune(48, '-')
	Field("host", "localhost:8080")
	KeyHint("q", "quit")
	KeyHint("r", "refresh")
	DividerRune(48, '-')
	InputLine("search> ", "foo", true)
	InputLine("filter> ", "bar", false)
	DividerRune(48, '-')
	StatusInfo("everything is fine")
	StatusWarn("something looks off")
	StatusError("something broke")
	DividerRune(48, '=')
}

// TestVisualComponents renders three full scenes to stdout for visual inspection.
// Scene 1 uses hardcoded overrides; Scene 2 is a component reference showing all
// components and their signatures; Scene 3 uses tui.config.toml so it renders
// last and you can immediately see the effect of TOML tweaks.
// Run with:
//
//	go test -v -run TestVisualComponents ./...
func TestVisualComponents(t *testing.T) {
	tui := NewTUI()

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

	tui.Divider(&DividerParams{Rune: '='})
	tui.MenuTitle(&TitleParams{Text: "Settings"})
	tui.Divider(&DividerParams{Rune: '='})
	tui.Menu(&MenuParams{Items: []MenuEntry{
		{Label: "Network", Selected: false},
		{Label: "Storage", Selected: true},
		{Label: "Security", Selected: false},
	}})
	tui.Divider(&DividerParams{Rune: '='})
	tui.Selector(&SelectorParams{
		Label:   "theme",
		Items:   []string{"dark", "light", "system"},
		Current: 0,
	})
	tui.Input(&InputParams{Label: "alias", Value: "dev-box", Active: true})
	tui.Divider(&DividerParams{Rune: '='})

	// Scene 2: Component reference — all components with their signatures
	t.Log("Scene 2: Component Reference (all components & signatures)")
	demoConfig(t)

	tui.Divider(&DividerParams{Rune: '='})
	tui.MenuTitle(&TitleParams{Text: "Component Reference"})
	tui.Divider(&DividerParams{Rune: '='})

	// TUI.MenuTitle(p *TitleParams)
	//   TitleParams{Text string, Width int}
	Field("MenuTitle", "TitleParams{Text, Width}")
	tui.MenuTitle(&TitleParams{Text: "Example Title"})
	tui.Divider(&DividerParams{})

	// TUI.Menu(p *MenuParams)
	//   MenuParams{Items []MenuEntry, Width int}
	//   MenuEntry{Label string, Selected bool}
	Field("Menu", "MenuParams{Items []MenuEntry, Width}")
	Field("MenuEntry", "{Label, Selected}")
	tui.Menu(&MenuParams{Items: []MenuEntry{
		{Label: "Selected item", Selected: true},
		{Label: "Normal item", Selected: false},
	}})
	tui.Divider(&DividerParams{})

	// TUI.Selector(p *SelectorParams)
	//   SelectorParams{Label string, Items []string, Current int, Width int}
	Field("Selector", "SelectorParams{Label, Items []string, Current, Width}")
	tui.Selector(&SelectorParams{
		Label:   "example",
		Items:   []string{"option-a", "option-b", "option-c"},
		Current: 1,
	})
	tui.Divider(&DividerParams{})

	// TUI.Input(p *InputParams)
	//   InputParams{Label string, Value string, Active bool, Width int}
	Field("Input", "InputParams{Label, Value, Active, Width}")
	tui.Input(&InputParams{Label: "active", Value: "typing", Active: true})
	tui.Input(&InputParams{Label: "inactive", Value: "static", Active: false})
	tui.Divider(&DividerParams{})

	// TUI.Divider(p *DividerParams)
	//   DividerParams{Rune rune, Width int}
	Field("Divider", "DividerParams{Rune, Width}")
	tui.Divider(&DividerParams{Rune: '~', Width: 20})
	tui.Divider(&DividerParams{})

	// Flat helpers (no TUI receiver)
	tui.MenuTitle(&TitleParams{Text: "Flat Helpers"})
	tui.Divider(&DividerParams{})

	Field("MenuItem", "MenuItem(index int, label string, selected bool)")
	MenuItem(1, "example", true)
	MenuItem(2, "example", false)

	Field("Field", "Field(label string, value any)")
	Field("key", "value")

	Field("KeyHint", "KeyHint(key, desc string)")
	KeyHint("q", "quit")

	Field("InputLine", "InputLine(prefix, value string, active bool)")
	InputLine("prompt> ", "text", true)
	InputLine("prompt> ", "text", false)

	Field("StatusInfo", "StatusInfo(msg string)")
	StatusInfo("info message")
	Field("StatusWarn", "StatusWarn(msg string)")
	StatusWarn("warning message")
	Field("StatusError", "StatusError(msg string)")
	StatusError("error message")

	Field("Divider", "Divider(width int)")
	Divider(20)
	Field("DividerRune", "DividerRune(width int, r rune)")
	DividerRune(20, '*')

	tui.Divider(&DividerParams{Rune: '='})

	// Scene 3: config-driven (from tui.config.toml) — rendered last
	t.Log("Scene 3: Config-driven layout (from tui.config.toml)")
	demoConfig(t)

	tui.Divider(&DividerParams{Rune: '='})
	tui.MenuTitle(&TitleParams{Text: "Main Menu"})
	tui.Divider(&DividerParams{})
	tui.Menu(&MenuParams{Items: []MenuEntry{
		{Label: "Status", Selected: true},
		{Label: "Settings", Selected: false},
		{Label: "Logs", Selected: false},
		{Label: "Quit", Selected: false},
	}})
	tui.Selector(&SelectorParams{
		Label:   "mode",
		Items:   []string{"debug", "verbose", "silent"},
		Current: 1,
	})
	tui.Input(&InputParams{Label: "filter", Value: "error", Active: false})
	tui.Input(&InputParams{Label: "output", Value: "stdout", Active: true})
	tui.Divider(&DividerParams{Rune: '='})
}

// TestVisualTreeView renders a tree view for visual inspection.
// Run with:
//
//	go test -v -run TestVisualTreeView ./...
func TestVisualTreeView(t *testing.T) {
	demoConfig(t)
	tui := NewTUI()

	tui.Divider(&DividerParams{Rune: '='})
	tui.MenuTitle(&TitleParams{Text: "TreeView"})
	tui.Divider(&DividerParams{})

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

	Field("TreeView", "TreeViewParams{Nodes []TreeNode, Width, ShowIndex}")
	tui.TreeView(&TreeViewParams{Nodes: nodes})
	tui.Divider(&DividerParams{})

	Field("TreeView (ShowIndex)", "ShowIndex: true")
	tui.TreeView(&TreeViewParams{Nodes: nodes, ShowIndex: true})
	tui.Divider(&DividerParams{Rune: '='})
}

// TestVisualProgressBar renders a progress bar for visual inspection.
// Run with:
//
//	go test -v -run TestVisualProgressBar ./...
func TestVisualProgressBar(t *testing.T) {
	demoConfig(t)
	tui := NewTUI()

	tui.Divider(&DividerParams{Rune: '='})
	tui.MenuTitle(&TitleParams{Text: "ProgressBar"})
	tui.Divider(&DividerParams{})

	Field("ProgressBar", "NewProgressBar(dst, ProgressBarParams{...})")
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

	tui.Divider(&DividerParams{Rune: '='})
}

// TestVisualOperationSummary renders operation summaries for visual inspection.
// Run with:
//
//	go test -v -run TestVisualOperationSummary ./...
func TestVisualOperationSummary(t *testing.T) {
	demoConfig(t)
	tui := NewTUI()

	tui.Divider(&DividerParams{Rune: '='})
	tui.MenuTitle(&TitleParams{Text: "OperationSummary"})
	tui.Divider(&DividerParams{})

	// OK case
	Field("OperationSummary", "OK=true")
	pt := NewPhaseTimer()
	pt.Begin("connect")
	time.Sleep(10 * time.Millisecond)
	pt.Begin("transfer")
	time.Sleep(15 * time.Millisecond)
	pt.Begin("verify")
	time.Sleep(5 * time.Millisecond)
	pt.End()

	tui.OperationSummary(&OperationSummaryParams{
		Title: "File Upload",
		OK:    true,
		Fields: []SummaryField{
			{Label: "Target", Value: "s3://bucket/key"},
			{Label: "Size", Value: "10.5 MiB"},
			{Label: "Files", Value: "42"},
		},
		Timer: pt,
	})

	tui.Divider(&DividerParams{})

	// FAILED case
	Field("OperationSummary", "OK=false")
	tui.OperationSummary(&OperationSummaryParams{
		Title: "Database Migration",
		OK:    false,
		Fields: []SummaryField{
			{Label: "Schema", Value: "v3 → v4"},
			{Label: "Error", Value: "foreign key constraint"},
		},
	})

	tui.Divider(&DividerParams{Rune: '='})
}
