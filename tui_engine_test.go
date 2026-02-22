package tui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	smplog "github.com/danmuck/smplog"
)

// captureStdout redirects os.Stdout for the duration of fn and returns
// everything written to it as a string.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = orig
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	return buf.String()
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

	out := captureStdout(t, func() {
		if _, err := WriteAt(3, 5, Configured().Colors.Menu, "node:%d", 7); err != nil {
			t.Fatalf("WriteAt: %v", err)
		}
	})

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

	out := captureStdout(t, func() {
		WriteAt(1, 1, Configured().Colors.Menu, "plain") //nolint:errcheck
	})

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

	out := captureStdout(t, func() {
		MenuItem(2, "services", true) //nolint:errcheck
	})

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
	out := captureStdout(t, func() {
		if err := BeginFrame(); err != nil {
			t.Fatalf("BeginFrame: %v", err)
		}
		if err := EndFrame(); err != nil {
			t.Fatalf("EndFrame: %v", err)
		}
	})

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

	out := captureStdout(t, func() {
		KeyHint("q", "quit")            //nolint:errcheck
		Field("mode", "diag")           //nolint:errcheck
		InputLine("select> ", "2", true) //nolint:errcheck
	})

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

// TestVisualFlatHelpers renders all flat output helpers to stdout so you can
// inspect them visually during development. Run with:
//
//	go test -v -run TestVisualFlatHelpers ./...
func TestVisualFlatHelpers(t *testing.T) {
	Configure(Config{
		NoColor: false,
		Colors: ColorConfig{
			Title:   smplog.StyleColor256(15),
			Menu:    smplog.StyleColor256(14),
			Prompt:  smplog.StyleColor256(10),
			Data:    smplog.StyleColor256(7),
			Divider: smplog.StyleColor256(8),
			Error:   smplog.StyleColor256(9),
		},
		TUI: TUIConfig{
			MenuSelectedPrefix:   ">",
			MenuUnselectedPrefix: " ",
			MenuIndexWidth:       2,
			InputCursor:          "_",
			DividerWidth:         48,
		},
	})
	t.Cleanup(func() { Configure(DefaultConfig()) })

	t.Log("── flat helpers ──")
	DividerRune(48, '=')
	MenuItem(1, "selected item", true)
	MenuItem(2, "normal item", false)
	DividerRune(48, '-')
	Field("host", "localhost:8080")
	KeyHint("q", "quit")
	KeyHint("r", "refresh")
	DividerRune(48, '-')
	InputLine("search> ", "foo", true)
	InputLine("filter> ", "bar", false)
	DividerRune(48, '-')
	StatusInfo("everything is fine")
	StatusWarn("something looks off")
	StatusError("something broke")
	DividerRune(48, '=')
}
