package tui

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	smplog "github.com/danmuck/smplog"
)

// TUI is a stateless component renderer. Construct with NewTUI.
type TUI struct{}

// NewTUI returns a new TUI renderer.
func NewTUI() TUI { return TUI{} }

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
// It clips, colorizes, optionally centers, then writes a line to stdout.
func writeComponent(cfg Config, color, plainContent string, width int) (int, error) {
	if width > 0 {
		plainContent = Clip(width, plainContent)
	}
	colored := smplog.Colorize(color, plainContent, cfg.NoColor)
	return writeComposite(cfg, colored, utf8.RuneCountInString(plainContent))
}

// writeComposite is the centering choke-point for multi-color component lines.
// line must already be fully colorized. plainWidth is the visible rune count of
// line (without ANSI escape bytes) and is used for centering math.
func writeComposite(cfg Config, line string, plainWidth int) (int, error) {
	var output string
	if cfg.TUI.Centered && cfg.TUI.MaxWidth > 0 {
		pad := max(cfg.TUI.MaxWidth-plainWidth, 0)
		left := pad / 2
		right := pad - left
		output = strings.Repeat(" ", left) + line + strings.Repeat(" ", right)
	} else {
		output = line
	}
	return fmt.Fprintln(os.Stdout, output)
}

// blockLine holds a pre-colorized line and its visible rune count.
type blockLine struct {
	colored    string
	plainWidth int
}

// writeBlock renders lines as a left-aligned block. When centering is active,
// all lines are padded to the widest line's width so they share the same left margin.
func writeBlock(cfg Config, lines []blockLine) {
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
		writeComposite(cfg, l.colored, blockWidth) //nolint:errcheck
	}
}

// Refresh clears the screen and moves the cursor to position (1, 1).
func (TUI) Refresh() error {
	if _, err := ClearScreen(); err != nil {
		return err
	}
	_, err := MoveTo(1, 1)
	return err
}
