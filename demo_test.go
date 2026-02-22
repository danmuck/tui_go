package tui

import (
	"testing"

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

// TestVisualComponents renders two full scenes to stdout for visual inspection.
// Scene 1 uses hardcoded overrides; Scene 2 uses tui.config.toml so it renders
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

	// Scene 2: config-driven (from tui.config.toml) — rendered last
	t.Log("Scene 2: Config-driven layout (from tui.config.toml)")
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
