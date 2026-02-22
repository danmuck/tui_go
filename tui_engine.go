package tui

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	smplog "github.com/danmuck/smplog"
)

// EnterAltScreen switches the terminal to an alternate screen buffer.
func EnterAltScreen() (int, error) {
	return writeANSI("\033[?1049h")
}

// ExitAltScreen returns the terminal to the main screen buffer.
func ExitAltScreen() (int, error) {
	return writeANSI("\033[?1049l")
}

// HideCursor hides the terminal cursor.
func HideCursor() (int, error) {
	return writeANSI("\033[?25l")
}

// ShowCursor shows the terminal cursor.
func ShowCursor() (int, error) {
	return writeANSI("\033[?25h")
}

// MoveTo moves the cursor to a 1-based row/column position.
func MoveTo(row, col int) (int, error) {
	return writeANSI(fmt.Sprintf("\033[%d;%dH", maxOne(row), maxOne(col)))
}

// ClearScreen clears the full terminal viewport.
func ClearScreen() (int, error) {
	return writeANSI("\033[2J")
}

// ClearLine clears the current line and returns the cursor to column 1.
func ClearLine() (int, error) {
	return writeANSI("\033[2K\r")
}

// WriteAt moves to row/col and writes a formatted message.
// Color output is controlled by Config.NoColor.
// The returned byte count includes the ANSI cursor-position sequence and
// should not be used as a visible-character-width measurement.
func WriteAt(row, col int, color, format string, v ...any) (int, error) {
	n, err := MoveTo(row, col)
	if err != nil {
		return n, err
	}
	m, err := smplog.Fcolorf(os.Stdout, color, format, v...)
	return n + m, err
}

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

// MenuItem writes a compact menu entry.
// Selected entries are rendered with title color; others use menu color.
func MenuItem(index int, label string, selected bool) (int, error) {
	cfg := Configured()
	color := cfg.Colors.Menu
	prefix := cfg.TUI.MenuUnselectedPrefix
	if selected {
		color = cfg.Colors.Title
		prefix = cfg.TUI.MenuSelectedPrefix
	}
	text := fmt.Sprintf("%s %*d) %s", prefix, cfg.TUI.MenuIndexWidth, index, label)
	return smplog.Fcolorf(os.Stdout, color, "%s", text)
}

// KeyHint writes a keyboard hint using prompt and data colors.
func KeyHint(key, desc string) (int, error) {
	cfg := Configured()
	keyText := smplog.Colorize(cfg.Colors.Prompt, key, cfg.NoColor)
	descText := smplog.Colorize(cfg.Colors.Data, desc, cfg.NoColor)
	return smplog.Fprintf(os.Stdout, "[%s] %s", keyText, descText)
}

// Field writes a key/value pair using prompt and data colors.
func Field(label string, value any) (int, error) {
	cfg := Configured()
	labelText := smplog.Colorize(cfg.Colors.Prompt, label, cfg.NoColor)
	valueText := smplog.Colorize(cfg.Colors.Data, fmt.Sprint(value), cfg.NoColor)
	return smplog.Fprintf(os.Stdout, "%s: %s", labelText, valueText)
}

// StatusInfo writes an info-status message using the data color.
func StatusInfo(msg string) (int, error) {
	cfg := Configured()
	return smplog.Fcolorf(os.Stdout, cfg.Colors.Data, "%s", msg)
}

// StatusWarn writes a warning-status message using the prompt color.
func StatusWarn(msg string) (int, error) {
	cfg := Configured()
	return smplog.Fcolorf(os.Stdout, cfg.Colors.Prompt, "%s", msg)
}

// StatusError writes an error-status message using the configured error color.
func StatusError(msg string) (int, error) {
	cfg := Configured()
	return smplog.Fcolorf(os.Stdout, cfg.Colors.Error, "%s", msg)
}

// InputLine writes a compact prompt/value input row.
// If active, the configured cursor character is appended.
func InputLine(prefix, value string, active bool) (int, error) {
	cfg := Configured()
	prefixText := smplog.Colorize(cfg.Colors.Prompt, prefix, cfg.NoColor)
	valueText := smplog.Colorize(cfg.Colors.Data, value, cfg.NoColor)
	if !active {
		return smplog.Fprintf(os.Stdout, "%s%s", prefixText, valueText)
	}
	cursor := smplog.Colorize(cfg.Colors.Prompt, cfg.TUI.InputCursor, cfg.NoColor)
	return smplog.Fprintf(os.Stdout, "%s%s%s", prefixText, valueText, cursor)
}

// Divider writes a horizontal divider using '-' and the divider color.
func Divider(width int) (int, error) {
	return DividerRune(width, '-')
}

// DividerRune writes a horizontal divider using r and the divider color.
func DividerRune(width int, r rune) (int, error) {
	cfg := Configured()
	if width <= 0 {
		width = cfg.TUI.DividerWidth
	}
	if r == 0 {
		r = '-'
	}
	line := strings.Repeat(string(r), width)
	return smplog.Fcolorf(os.Stdout, cfg.Colors.Divider, "%s", line)
}

// BeginFrame switches to alt-screen, hides the cursor, clears, and positions at 1,1.
func BeginFrame() error {
	if _, err := EnterAltScreen(); err != nil {
		return err
	}
	if _, err := HideCursor(); err != nil {
		return err
	}
	if _, err := ClearScreen(); err != nil {
		return err
	}
	_, err := MoveTo(1, 1)
	return err
}

// EndFrame restores the cursor and returns to the main screen.
func EndFrame() error {
	if _, err := ShowCursor(); err != nil {
		return err
	}
	_, err := ExitAltScreen()
	return err
}

func writeANSI(seq string) (int, error) {
	return smplog.Fprint(os.Stdout, seq)
}

func maxOne(n int) int {
	if n < 1 {
		return 1
	}
	return n
}
