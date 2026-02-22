package tui

import (
	"os"
	"path/filepath"
	"testing"

	smplog "github.com/danmuck/smplog"
)

func TestDefaultConfigMatchesSmplogColors(t *testing.T) {
	cfg := DefaultConfig()
	d := smplog.DefaultColors()
	if cfg.Colors.Menu != d.Menu {
		t.Fatalf("Menu color mismatch: got %q want %q", cfg.Colors.Menu, d.Menu)
	}
	if cfg.Colors.Title != d.Title {
		t.Fatalf("Title color mismatch: got %q want %q", cfg.Colors.Title, d.Title)
	}
	if cfg.Colors.Prompt != d.Prompt {
		t.Fatalf("Prompt color mismatch: got %q want %q", cfg.Colors.Prompt, d.Prompt)
	}
	if cfg.Colors.Data != d.Data {
		t.Fatalf("Data color mismatch: got %q want %q", cfg.Colors.Data, d.Data)
	}
	if cfg.Colors.Divider != d.Divider {
		t.Fatalf("Divider color mismatch: got %q want %q", cfg.Colors.Divider, d.Divider)
	}
}

func TestDefaultConfigTUIDefaults(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.TUI.MenuSelectedPrefix != ">" {
		t.Fatalf("MenuSelectedPrefix: got %q want %q", cfg.TUI.MenuSelectedPrefix, ">")
	}
	if cfg.TUI.MenuIndexWidth != 2 {
		t.Fatalf("MenuIndexWidth: got %d want 2", cfg.TUI.MenuIndexWidth)
	}
	if cfg.TUI.DividerWidth != 64 {
		t.Fatalf("DividerWidth: got %d want 64", cfg.TUI.DividerWidth)
	}
	if cfg.TUI.InputCursor != "_" {
		t.Fatalf("InputCursor: got %q want %q", cfg.TUI.InputCursor, "_")
	}
}

func TestConfigureRoundTrip(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })

	custom := Config{
		NoColor: true,
		TUI:     TUIConfig{DividerWidth: 99, InputCursor: "|"},
	}
	Configure(custom)

	got := Configured()
	if !got.NoColor {
		t.Fatal("expected NoColor=true after Configure")
	}
	if got.TUI.DividerWidth != 99 {
		t.Fatalf("DividerWidth: got %d want 99", got.TUI.DividerWidth)
	}
	if got.TUI.InputCursor != "|" {
		t.Fatalf("InputCursor: got %q want |", got.TUI.InputCursor)
	}
}

func TestNormalizeConfigFillsZeroTUIFields(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })

	Configure(Config{}) // all zero values
	got := Configured()
	def := DefaultConfig()

	if got.TUI.MenuSelectedPrefix != def.TUI.MenuSelectedPrefix {
		t.Fatalf("expected default MenuSelectedPrefix, got %q", got.TUI.MenuSelectedPrefix)
	}
	if got.TUI.DividerWidth != def.TUI.DividerWidth {
		t.Fatalf("expected default DividerWidth, got %d", got.TUI.DividerWidth)
	}
}

func TestNormalizeConfigFillsZeroColors(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })

	Configure(Config{}) // empty colors
	got := Configured()
	def := DefaultConfig()

	if got.Colors.Menu != def.Colors.Menu {
		t.Fatalf("expected default Menu color, got %q", got.Colors.Menu)
	}
}

func TestConfigFromFileLoadsToml(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tui.config.toml")
	content := `
[tui]
divider_width = 32
input_cursor  = "|"
centered      = true
max_width     = 80

[colors]
menu  = 14
title = 15
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := ConfigFromFile(path)
	if err != nil {
		t.Fatalf("ConfigFromFile: %v", err)
	}
	if cfg.TUI.DividerWidth != 32 {
		t.Fatalf("DividerWidth: got %d want 32", cfg.TUI.DividerWidth)
	}
	if cfg.TUI.InputCursor != "|" {
		t.Fatalf("InputCursor: got %q want |", cfg.TUI.InputCursor)
	}
	if !cfg.TUI.Centered {
		t.Fatal("expected Centered=true")
	}
	if cfg.TUI.MaxWidth != 80 {
		t.Fatalf("MaxWidth: got %d want 80", cfg.TUI.MaxWidth)
	}
	if cfg.Colors.Menu == "" {
		t.Fatal("expected non-empty Menu color from TOML")
	}
}

func TestConfigFromFileMissingReturnsError(t *testing.T) {
	_, err := ConfigFromFile("/nonexistent/tui.config.toml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestConfigFromFileInvalidColorIndexReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tui.config.toml")
	content := `
[colors]
menu = -1
title = 999
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := ConfigFromFile(path)
	if err != nil {
		t.Fatalf("ConfigFromFile: %v", err)
	}
	if cfg.Colors.Menu != "" {
		t.Fatalf("expected empty Menu color for out-of-range index, got %q", cfg.Colors.Menu)
	}
	if cfg.Colors.Title != "" {
		t.Fatalf("expected empty Title color for out-of-range index, got %q", cfg.Colors.Title)
	}
}
