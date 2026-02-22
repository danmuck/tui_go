package tui

import (
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"
	smplog "github.com/danmuck/smplog"
)

const (
	defaultMenuSelectedPrefix   = ">"
	defaultMenuUnselectedPrefix = " "
	defaultMenuIndexWidth       = 2
	defaultInputCursor          = "_"
	defaultDividerWidth         = 64
)

// TUIConfig controls layout and input rendering for tui_go components.
type TUIConfig struct {
	MenuSelectedPrefix   string
	MenuUnselectedPrefix string
	MenuIndexWidth       int
	InputCursor          string
	DividerWidth         int
	MaxWidth             int  // 0 = unconstrained
	Centered             bool // center content within MaxWidth when true
}

// ColorConfig holds raw ANSI escape strings for each tui color role.
// Empty fields fall back to smplog.DefaultColors() at normalize time.
type ColorConfig struct {
	Menu    string
	Title   string
	Prompt  string
	Data    string
	Divider string
	Error   string
}

// Config is the runtime configuration for tui_go.
type Config struct {
	TUI     TUIConfig
	Colors  ColorConfig
	NoColor bool
}

var (
	stateMu       sync.RWMutex
	currentConfig Config
)

func init() {
	currentConfig = DefaultConfig()
}

// DefaultConfig returns the default Config populated from smplog's default palette.
func DefaultConfig() Config {
	d := smplog.DefaultColors()
	return Config{
		TUI: TUIConfig{
			MenuSelectedPrefix:   defaultMenuSelectedPrefix,
			MenuUnselectedPrefix: defaultMenuUnselectedPrefix,
			MenuIndexWidth:       defaultMenuIndexWidth,
			InputCursor:          defaultInputCursor,
			DividerWidth:         defaultDividerWidth,
		},
		Colors: ColorConfig{
			Menu:    d.Menu,
			Title:   d.Title,
			Prompt:  d.Prompt,
			Data:    d.Data,
			Divider: d.Divider,
			Error:   smplog.StyleColor256(smplog.Red),
		},
	}
}

// Configure sets the package-global config. Zero-valued TUI fields and empty
// color strings are replaced with defaults from DefaultConfig before storing.
func Configure(cfg Config) {
	cfg = normalizeConfig(cfg)
	stateMu.Lock()
	currentConfig = cfg
	stateMu.Unlock()
}

// Configured returns the current package-global config.
func Configured() Config {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return currentConfig
}

func normalizeConfig(cfg Config) Config {
	def := DefaultConfig()
	t := &cfg.TUI
	dt := def.TUI
	if t.MenuSelectedPrefix == "" {
		t.MenuSelectedPrefix = dt.MenuSelectedPrefix
	}
	if t.MenuUnselectedPrefix == "" {
		t.MenuUnselectedPrefix = dt.MenuUnselectedPrefix
	}
	if t.MenuIndexWidth <= 0 {
		t.MenuIndexWidth = dt.MenuIndexWidth
	}
	if t.InputCursor == "" {
		t.InputCursor = dt.InputCursor
	}
	if t.DividerWidth <= 0 {
		t.DividerWidth = dt.DividerWidth
	}
	c := &cfg.Colors
	dc := def.Colors
	if c.Menu == "" {
		c.Menu = dc.Menu
	}
	if c.Title == "" {
		c.Title = dc.Title
	}
	if c.Prompt == "" {
		c.Prompt = dc.Prompt
	}
	if c.Data == "" {
		c.Data = dc.Data
	}
	if c.Divider == "" {
		c.Divider = dc.Divider
	}
	if c.Error == "" {
		c.Error = dc.Error
	}
	return cfg
}

// fileConfig is the TOML-decodable shape of Config.
type fileConfig struct {
	TUI     tuiFileConfig   `toml:"tui"`
	Colors  colorFileConfig `toml:"colors"`
	NoColor bool            `toml:"no_color"`
}

// tuiFileConfig is the on-disk TOML shape of TUIConfig.
type tuiFileConfig struct {
	MenuSelectedPrefix   string `toml:"menu_selected_prefix"`
	MenuUnselectedPrefix string `toml:"menu_unselected_prefix"`
	MenuIndexWidth       int    `toml:"menu_index_width"`
	InputCursor          string `toml:"input_cursor"`
	DividerWidth         int    `toml:"divider_width"`
	MaxWidth             int    `toml:"max_width"`
	Centered             bool   `toml:"centered"`
}

// colorFileConfig is the on-disk TOML shape of ColorConfig.
// Color values are 256-color palette indices (0–255); absent fields are nil.
type colorFileConfig struct {
	Menu    *int `toml:"menu"`
	Title   *int `toml:"title"`
	Prompt  *int `toml:"prompt"`
	Data    *int `toml:"data"`
	Divider *int `toml:"divider"`
}

func color256(p *int) string {
	if p == nil {
		return ""
	}
	v := *p
	if v < 0 || v > 255 {
		return ""
	}
	return smplog.StyleColor256(v)
}

// ConfigFromFile parses a TOML file at path and returns a Config.
// Absent fields keep zero values; Configure fills them with defaults.
func ConfigFromFile(path string) (Config, error) {
	var fc fileConfig
	if _, err := toml.DecodeFile(path, &fc); err != nil {
		return Config{}, fmt.Errorf("tui: parse config %q: %w", path, err)
	}
	return Config{
		NoColor: fc.NoColor,
		TUI: TUIConfig{
			MenuSelectedPrefix:   fc.TUI.MenuSelectedPrefix,
			MenuUnselectedPrefix: fc.TUI.MenuUnselectedPrefix,
			MenuIndexWidth:       fc.TUI.MenuIndexWidth,
			InputCursor:          fc.TUI.InputCursor,
			DividerWidth:         fc.TUI.DividerWidth,
			MaxWidth:             fc.TUI.MaxWidth,
			Centered:             fc.TUI.Centered,
		},
		Colors: ColorConfig{
			Menu:    color256(fc.Colors.Menu),
			Title:   color256(fc.Colors.Title),
			Prompt:  color256(fc.Colors.Prompt),
			Data:    color256(fc.Colors.Data),
			Divider: color256(fc.Colors.Divider),
		},
	}, nil
}
