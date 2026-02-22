package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	smplog "github.com/danmuck/smplog"
)

// --- Terminal control (TERM suffix) ---

func (t TUI) writeANSI(seq string) (int, error) {
	return smplog.Fprint(t.w, seq)
}

// EnterAltScreenTERM switches the terminal to an alternate screen buffer.
func (t TUI) EnterAltScreenTERM() (int, error) { return t.writeANSI("\033[?1049h") }

// ExitAltScreenTERM returns the terminal to the main screen buffer.
func (t TUI) ExitAltScreenTERM() (int, error) { return t.writeANSI("\033[?1049l") }

// HideCursorTERM hides the terminal cursor.
func (t TUI) HideCursorTERM() (int, error) { return t.writeANSI("\033[?25l") }

// ShowCursorTERM shows the terminal cursor.
func (t TUI) ShowCursorTERM() (int, error) { return t.writeANSI("\033[?25h") }

// MoveToTERM moves the cursor to a 1-based row/column position.
func (t TUI) MoveToTERM(row, col int) (int, error) {
	return t.writeANSI(fmt.Sprintf("\033[%d;%dH", maxOne(row), maxOne(col)))
}

// ClearScreenTERM clears the full terminal viewport.
func (t TUI) ClearScreenTERM() (int, error) { return t.writeANSI("\033[2J") }

// ClearLineTERM clears the current line and returns the cursor to column 1.
func (t TUI) ClearLineTERM() (int, error) { return t.writeANSI("\033[2K\r") }

// WriteAtTERM moves to row/col and writes a formatted message.
func (t TUI) WriteAtTERM(row, col int, color, format string, v ...any) (int, error) {
	n, err := t.MoveToTERM(row, col)
	if err != nil {
		return n, err
	}
	m, err := smplog.Fcolorf(t.w, color, format, v...)
	return n + m, err
}

// BeginFrameTERM switches to alt-screen, hides the cursor, clears, and positions at 1,1.
func (t TUI) BeginFrameTERM() error {
	if _, err := t.EnterAltScreenTERM(); err != nil {
		return err
	}
	if _, err := t.HideCursorTERM(); err != nil {
		return err
	}
	if _, err := t.ClearScreenTERM(); err != nil {
		return err
	}
	_, err := t.MoveToTERM(1, 1)
	return err
}

// EndFrameTERM restores the cursor and returns to the main screen.
func (t TUI) EndFrameTERM() error {
	if _, err := t.ShowCursorTERM(); err != nil {
		return err
	}
	_, err := t.ExitAltScreenTERM()
	return err
}

func maxOne(n int) int {
	if n < 1 {
		return 1
	}
	return n
}

// --- Text utilities (package-level) ---

// Clip truncates s to width runes.
func Clip(width int, s string) string {
	if width <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= width {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	n := 0
	for _, r := range s {
		if n >= width {
			break
		}
		b.WriteRune(r)
		n++
	}
	return b.String()
}

// PadLeft left-pads s with spaces up to width runes.
func PadLeft(width int, s string) string {
	s = Clip(width, s)
	return strings.Repeat(" ", max(width-utf8.RuneCountInString(s), 0)) + s
}

// PadRight right-pads s with spaces up to width runes.
func PadRight(width int, s string) string {
	s = Clip(width, s)
	return s + strings.Repeat(" ", max(width-utf8.RuneCountInString(s), 0))
}

// Center centers s within width runes.
func Center(width int, s string) string {
	s = Clip(width, s)
	pad := max(width-utf8.RuneCountInString(s), 0)
	left := pad / 2
	right := pad - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

// --- Flat utility helpers (FU suffix) ---

// MenuItemFU writes a compact menu entry.
// Selected entries are rendered with title color; others use menu color.
func (t TUI) MenuItemFU(index int, label string, selected bool) (int, error) {
	cfg := Configured()
	color := cfg.Colors.Menu
	prefix := cfg.TUI.MenuUnselectedPrefix
	if selected {
		color = cfg.Colors.Title
		prefix = cfg.TUI.MenuSelectedPrefix
	}
	text := fmt.Sprintf("%s %*d) %s", prefix, cfg.TUI.MenuIndexWidth, index, label)
	return smplog.Fcolorf(t.w, color, "%s", text)
}

// KeyHintFU writes a keyboard hint using prompt and data colors.
func (t TUI) KeyHintFU(key, desc string) (int, error) {
	cfg := Configured()
	keyText := smplog.Colorize(cfg.Colors.Prompt, key, cfg.NoColor)
	descText := smplog.Colorize(cfg.Colors.Data, desc, cfg.NoColor)
	return smplog.Fprintf(t.w, "[%s] %s", keyText, descText)
}

// FieldFU writes a key/value pair using prompt and data colors.
func (t TUI) FieldFU(label string, value any) (int, error) {
	cfg := Configured()
	labelText := smplog.Colorize(cfg.Colors.Prompt, label, cfg.NoColor)
	valueText := smplog.Colorize(cfg.Colors.Data, fmt.Sprint(value), cfg.NoColor)
	return smplog.Fprintf(t.w, "%s: %s", labelText, valueText)
}

// StatusInfoFU writes an info-status message using the data color.
func (t TUI) StatusInfoFU(msg string) (int, error) {
	cfg := Configured()
	return smplog.Fcolorf(t.w, cfg.Colors.Data, "%s", msg)
}

// StatusWarnFU writes a warning-status message using the prompt color.
func (t TUI) StatusWarnFU(msg string) (int, error) {
	cfg := Configured()
	return smplog.Fcolorf(t.w, cfg.Colors.Prompt, "%s", msg)
}

// StatusErrorFU writes an error-status message using the configured error color.
func (t TUI) StatusErrorFU(msg string) (int, error) {
	cfg := Configured()
	return smplog.Fcolorf(t.w, cfg.Colors.Error, "%s", msg)
}

// InputLineFU writes a compact prompt/value input row.
// If active, the configured cursor character is appended.
func (t TUI) InputLineFU(prefix, value string, active bool) (int, error) {
	cfg := Configured()
	prefixText := smplog.Colorize(cfg.Colors.Prompt, prefix, cfg.NoColor)
	valueText := smplog.Colorize(cfg.Colors.Data, value, cfg.NoColor)
	if !active {
		return smplog.Fprintf(t.w, "%s%s", prefixText, valueText)
	}
	cursor := smplog.Colorize(cfg.Colors.Prompt, cfg.TUI.InputCursor, cfg.NoColor)
	return smplog.Fprintf(t.w, "%s%s%s", prefixText, valueText, cursor)
}
