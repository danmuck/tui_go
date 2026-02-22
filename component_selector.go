package tui

import (
	"fmt"
	"unicode/utf8"

	smplog "github.com/danmuck/smplog"
)

// SelectorParams configures TUI.Selector.
type SelectorParams struct {
	Label   string
	Items   []string
	Current int // 0-based index of current selection
	Width   int
}

// Selector renders a "label: < current >" selector row to stdout.
// It uses prompt color for the label and data color for the current item.
func (TUI) Selector(p *SelectorParams) {
	cfg := Configured()

	var current string
	if p.Current >= 0 && p.Current < len(p.Items) {
		current = p.Items[p.Current]
	}

	label := p.Label
	width := effectiveWidth(p.Width, cfg)
	if width > 0 {
		// Reserve space for ": < " (4) + " >" (2) = 6 chars around current
		labelMax := max(width-utf8.RuneCountInString(current)-6, 0)
		label = Clip(labelMax, label)
		remaining := max(width-utf8.RuneCountInString(label)-4, 0) // ": < " = 4
		current = Clip(remaining-2, current)                        // " >" = 2
	}

	plainWidth := utf8.RuneCountInString(label) + utf8.RuneCountInString(current) + 6 // ": < " + " >"
	labelText := smplog.Colorize(cfg.Colors.Prompt, label, cfg.NoColor)
	currentText := smplog.Colorize(cfg.Colors.Data, current, cfg.NoColor)
	line := fmt.Sprintf("%s: < %s >", labelText, currentText)
	writeComposite(cfg, line, plainWidth) //nolint:errcheck
}
