package tui

import (
	"bytes"
	"strings"
	"testing"

	smplog "github.com/danmuck/smplog"
)

func newTestTUI() (TUI, *bytes.Buffer) {
	var buf bytes.Buffer
	return NewTUI(&buf), &buf
}

func TestClipPadCenter(t *testing.T) {
	if got := Clip(4, "abcdef"); got != "abcd" {
		t.Fatalf("Clip: got %q want %q", got, "abcd")
	}
	if got := PadLeft(6, "xy"); got != "    xy" {
		t.Fatalf("PadLeft: got %q want %q", got, "    xy")
	}
	if got := PadRight(6, "xy"); got != "xy    " {
		t.Fatalf("PadRight: got %q want %q", got, "xy    ")
	}
	if got := Center(7, "abc"); got != "  abc  " {
		t.Fatalf("Center: got %q want %q", got, "  abc  ")
	}
}

func TestWriteAtMovesCursorAndColorizes(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: false,
		Colors:  ColorConfig{Menu: smplog.StyleColor256(14)},
	})

	tui, buf := newTestTUI()
	if _, err := tui.WriteAtTERM(3, 5, Configured().Colors.Menu, "node:%d", 7); err != nil {
		t.Fatalf("WriteAtTERM: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "\x1b[3;5H") {
		t.Fatalf("expected cursor move in output: %q", out)
	}
	if !strings.Contains(out, "\x1b[38;5;14m") {
		t.Fatalf("expected color sequence in output: %q", out)
	}
	if !strings.Contains(out, "node:7") {
		t.Fatalf("expected payload in output: %q", out)
	}
}

func TestWriteAtRespectsNoColor(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: true,
		Colors:  ColorConfig{Menu: smplog.StyleColor256(14)},
	})

	tui, buf := newTestTUI()
	tui.WriteAtTERM(1, 1, Configured().Colors.Menu, "plain") //nolint:errcheck
	out := buf.String()

	if strings.Contains(out, "\x1b[38;5;14m") {
		t.Fatalf("expected no color with NoColor=true: %q", out)
	}
	if !strings.Contains(out, "plain") {
		t.Fatalf("expected payload in output: %q", out)
	}
}

func TestMenuItemSelectionUsesTitleColor(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: false,
		Colors: ColorConfig{
			Menu:  smplog.StyleColor256(14),
			Title: smplog.StyleColor256(15),
		},
		TUI: TUIConfig{
			MenuSelectedPrefix:   ">>",
			MenuUnselectedPrefix: "..",
			MenuIndexWidth:       4,
			InputCursor:          "|",
			DividerWidth:         72,
		},
	})

	tui, buf := newTestTUI()
	tui.MenuItemFU(2, "services", true) //nolint:errcheck
	out := buf.String()

	if !strings.Contains(out, "\x1b[38;5;15m") {
		t.Fatalf("expected selected title color in output: %q", out)
	}
	if !strings.Contains(out, ">>    2) services") {
		t.Fatalf("expected selected item payload in output: %q", out)
	}
}

func TestBeginEndFrameWritesExpectedSequences(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })

	tui, buf := newTestTUI()
	if err := tui.BeginFrameTERM(); err != nil {
		t.Fatalf("BeginFrameTERM: %v", err)
	}
	if err := tui.EndFrameTERM(); err != nil {
		t.Fatalf("EndFrameTERM: %v", err)
	}
	out := buf.String()

	required := []string{
		"\x1b[?1049h",
		"\x1b[?25l",
		"\x1b[2J",
		"\x1b[1;1H",
		"\x1b[?25h",
		"\x1b[?1049l",
	}
	for _, seq := range required {
		if !strings.Contains(out, seq) {
			t.Fatalf("expected sequence %q in output: %q", seq, out)
		}
	}
}

func TestKeyHintFieldAndInputLineNoColor(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: true,
		Colors: ColorConfig{
			Prompt: smplog.StyleColor256(10),
			Data:   smplog.StyleColor256(7),
		},
		TUI: TUIConfig{InputCursor: "|"},
	})

	tui, buf := newTestTUI()
	tui.KeyHintFU("q", "quit")            //nolint:errcheck
	tui.FieldFU("mode", "diag")           //nolint:errcheck
	tui.InputLineFU("select> ", "2", true) //nolint:errcheck
	out := buf.String()

	if strings.Contains(out, "\x1b[") {
		t.Fatalf("expected no ANSI with NoColor=true: %q", out)
	}
	if !strings.Contains(out, "[q] quit") {
		t.Fatalf("expected key hint in output: %q", out)
	}
	if !strings.Contains(out, "mode: diag") {
		t.Fatalf("expected field in output: %q", out)
	}
	if !strings.Contains(out, "select> 2|") {
		t.Fatalf("expected active input in output: %q", out)
	}
}
