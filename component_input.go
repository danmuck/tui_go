package tui

import (
	"fmt"
	"unicode/utf8"

	smplog "github.com/danmuck/smplog"
)

// InputParams configures TUI.Input.
type InputParams struct {
	Label  string
	Value  string
	Active bool // if true, appends configured cursor character
	Width  int
}

// Input renders a "label: value[cursor]" input row to stdout.
// It uses prompt color for the label and data color for the value.
// When Active is true, the configured cursor character is appended.
func (TUI) Input(p *InputParams) {
	cfg := Configured()

	value := p.Value
	width := effectiveWidth(p.Width, cfg)
	if width > 0 {
		labelRunes := utf8.RuneCountInString(p.Label) + 2 // ": " = 2
		cursorRunes := 0
		if p.Active {
			cursorRunes = utf8.RuneCountInString(cfg.TUI.InputCursor)
		}
		value = Clip(max(width-labelRunes-cursorRunes, 0), value)
	}

	cursorStr := ""
	if p.Active {
		cursorStr = cfg.TUI.InputCursor
	}
	plainWidth := utf8.RuneCountInString(p.Label) + 2 + utf8.RuneCountInString(value) + utf8.RuneCountInString(cursorStr)

	labelText := smplog.Colorize(cfg.Colors.Prompt, p.Label, cfg.NoColor)
	valueText := smplog.Colorize(cfg.Colors.Data, value, cfg.NoColor)
	var line string
	if p.Active {
		cursorText := smplog.Colorize(cfg.Colors.Prompt, cfg.TUI.InputCursor, cfg.NoColor)
		line = fmt.Sprintf("%s: %s%s", labelText, valueText, cursorText)
	} else {
		line = fmt.Sprintf("%s: %s", labelText, valueText)
	}
	writeComposite(cfg, line, plainWidth) //nolint:errcheck
}
