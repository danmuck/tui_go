package tui

import (
	"fmt"
	"unicode/utf8"

	smplog "github.com/danmuck/smplog"
)

// MenuEntry is a single item in a Menu component.
type MenuEntry struct {
	Label    string
	Selected bool
}

// MenuParams configures TUI.Menu.
type MenuParams struct {
	Items []MenuEntry
	Width int // 0 = TUIConfig.MaxWidth
}

// TitleParams configures TUI.MenuTitle.
type TitleParams struct {
	Text  string
	Width int
}

// Menu renders a list of MenuEntry items to stdout.
// Selected entries use title color; others use menu color.
// When centered, all items share the same left margin so their prefix markers
// stay visually aligned as a block.
func (TUI) Menu(p *MenuParams) {
	cfg := Configured()
	width := effectiveWidth(p.Width, cfg)

	lines := make([]blockLine, len(p.Items))
	for i, entry := range p.Items {
		color := cfg.Colors.Menu
		prefix := cfg.TUI.MenuUnselectedPrefix
		if entry.Selected {
			color = cfg.Colors.Title
			prefix = cfg.TUI.MenuSelectedPrefix
		}
		plain := fmt.Sprintf("%s %*d) %s", prefix, cfg.TUI.MenuIndexWidth, i+1, entry.Label)
		if width > 0 {
			plain = Clip(width, plain)
		}
		colored := smplog.Colorize(color, plain, cfg.NoColor)
		lines[i] = blockLine{colored: colored, plainWidth: utf8.RuneCountInString(plain)}
	}

	writeBlock(cfg, lines)
}

// MenuTitle renders a title string to stdout using the title color.
func (TUI) MenuTitle(p *TitleParams) {
	cfg := Configured()
	width := effectiveWidth(p.Width, cfg)
	text := cfg.TUI.MenuTitlePrefix + p.Text + cfg.TUI.MenuTitlePostfix
	writeComponent(cfg, cfg.Colors.Title, text, width) //nolint:errcheck
}
