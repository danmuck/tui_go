package tui

import (
	"strings"

	smplog "github.com/danmuck/smplog"
)

// DividerParams configures TUI.Divider.
type DividerParams struct {
	Rune  rune // 0 = '-'
	Width int  // 0 = TUIConfig.DividerWidth
}

// Divider renders a horizontal divider line to stdout using the divider color.
// The divider width comes from DividerParams.Width, then TUIConfig.DividerWidth.
// MaxWidth is only used as the centering axis, not as the divider length.
func (TUI) Divider(p *DividerParams) {
	cfg := Configured()

	r := p.Rune
	if r == 0 {
		r = '-'
	}

	width := p.Width
	if width <= 0 {
		width = cfg.TUI.DividerWidth
	}

	plain := strings.Repeat(string(r), width)
	colored := smplog.Colorize(cfg.Colors.Divider, plain, cfg.NoColor)
	writeComposite(cfg, colored, width) //nolint:errcheck
}
