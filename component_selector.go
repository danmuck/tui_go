package tui

import (
	"fmt"
	"unicode/utf8"

	smplog "github.com/danmuck/smplog"
)

// SelectorParams configures TUI.SelectorTC.
type SelectorParams struct {
	Label   string
	Items   []string
	Current int // 0-based index of current selection
	Width   int
}

// SelectorTC renders a "label: < current >" selector row.
// It uses prompt color for the label and data color for the current item.
func (t TUI) SelectorTC(p *SelectorParams) {
	cfg := Configured()

	var current string
	if p.Current >= 0 && p.Current < len(p.Items) {
		current = p.Items[p.Current]
	}

	label := p.Label

	plainWidth := utf8.RuneCountInString(label) + utf8.RuneCountInString(current) + 6 // ": < " + " >"
	labelText := smplog.Colorize(cfg.Colors.Prompt, label, cfg.NoColor)
	currentText := smplog.Colorize(cfg.Colors.Data, current, cfg.NoColor)
	line := fmt.Sprintf("%s: < %s >", labelText, currentText)
	t.writeComposite(cfg, line, plainWidth) //nolint:errcheck
}
