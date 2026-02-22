package tui

import (
	"strings"

	smplog "github.com/danmuck/smplog"
)

// DividerParams configures TUI.DividerTC.
type DividerParams struct {
	Rune  rune // 0 = '-'
	Width int  // 0 = TUIConfig.DividerWidth
}

// DividerTC renders a horizontal divider line using the divider color.
// The divider width comes from DividerParams.Width, then TUIConfig.DividerWidth.
// MaxWidth is only used as the centering axis, not as the divider length.
// Blank lines are printed before and after the divider for visual padding.
func (t TUI) DividerTC(p *DividerParams) {
	cfg := Configured()

	r := p.Rune
	if r == 0 {
		r = '-'
	}

	width := p.Width
	if width <= 0 {
		width = cfg.TUI.DividerWidth
	}

	smplog.Fprintln(t.w, "")  // blank line before
	plain := strings.Repeat(string(r), width)
	colored := smplog.Colorize(cfg.Colors.Divider, plain, cfg.NoColor)
	t.writeComposite(cfg, colored, width) //nolint:errcheck
	smplog.Fprintln(t.w, "")  // blank line after
}
