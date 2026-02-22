package tui

import (
	"io"
	"strings"
	"unicode/utf8"

	smplog "github.com/danmuck/smplog"
)

// TUI is a stateless component renderer. Construct with NewTUI.
type TUI struct {
	w io.Writer
}

// NewTUI returns a new TUI renderer.
func NewTUI(w io.Writer) TUI { return TUI{w: w} }

// effectiveWidth resolves the layout width to apply.
// Priority: paramWidth → cfg.TUI.MaxWidth → 0 (unconstrained).
func effectiveWidth(paramWidth int, cfg Config) int {
	if paramWidth > 0 {
		return paramWidth
	}
	if cfg.TUI.MaxWidth > 0 {
		return cfg.TUI.MaxWidth
	}
	return 0
}

// writeComponent is the single choke-point for all single-color component output.
// It colorizes, optionally centers, then writes a line.
func (t TUI) writeComponent(cfg Config, color, plainContent string, width int) (int, error) {
	colored := smplog.Colorize(color, plainContent, cfg.NoColor)
	return t.writeComposite(cfg, colored, utf8.RuneCountInString(plainContent))
}

// writeComposite is the centering choke-point for multi-color component lines.
// line must already be fully colorized. plainWidth is the visible rune count of
// line (without ANSI escape bytes) and is used for centering math.
func (t TUI) writeComposite(cfg Config, line string, plainWidth int) (int, error) {
	var output string
	if cfg.TUI.Centered && cfg.TUI.MaxWidth > 0 {
		pad := max(cfg.TUI.MaxWidth-plainWidth, 0)
		left := pad / 2
		right := pad - left
		output = strings.Repeat(" ", left) + line + strings.Repeat(" ", right)
	} else {
		output = line
	}
	return smplog.Fprintln(t.w, output)
}

// blockLine holds a pre-colorized line and its visible rune count.
type blockLine struct {
	colored    string
	plainWidth int
}

// writeBlock renders lines as a left-aligned block. When centering is active,
// all lines are padded to the widest line's width so they share the same left margin.
func (t TUI) writeBlock(cfg Config, lines []blockLine) {
	blockWidth := 0
	for _, l := range lines {
		if l.plainWidth > blockWidth {
			blockWidth = l.plainWidth
		}
	}
	for _, l := range lines {
		if l.plainWidth < blockWidth {
			l.colored += strings.Repeat(" ", blockWidth-l.plainWidth)
		}
		t.writeComposite(cfg, l.colored, blockWidth) //nolint:errcheck
	}
}

// RefreshTERM clears the screen and moves the cursor to position (1, 1).
func (t TUI) RefreshTERM() error {
	if _, err := t.ClearScreenTERM(); err != nil {
		return err
	}
	_, err := t.MoveToTERM(1, 1)
	return err
}
