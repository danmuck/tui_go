package tui

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode/utf8"

	smplog "github.com/danmuck/smplog"
)

// ---------- Menu ----------

func TestTUIMenuRendersAllItems(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: false,
		Colors: ColorConfig{
			Menu:  smplog.StyleColor256(14),
			Title: smplog.StyleColor256(15),
		},
		TUI: TUIConfig{
			MenuSelectedPrefix:   ">",
			MenuUnselectedPrefix: " ",
			MenuIndexWidth:       2,
		},
	})

	out := captureStdout(t, func() {
		NewTUI().Menu(&MenuParams{
			Items: []MenuEntry{
				{Label: "alpha", Selected: true},
				{Label: "beta", Selected: false},
			},
		})
	})

	// Both items should appear
	if !strings.Contains(out, "alpha") {
		t.Fatalf("expected 'alpha' in output: %q", out)
	}
	if !strings.Contains(out, "beta") {
		t.Fatalf("expected 'beta' in output: %q", out)
	}
	// Selected item uses title color
	if !strings.Contains(out, "\x1b[38;5;15m") {
		t.Fatalf("expected title color for selected item: %q", out)
	}
	// Unselected item uses menu color
	if !strings.Contains(out, "\x1b[38;5;14m") {
		t.Fatalf("expected menu color for unselected item: %q", out)
	}
	// Item numbering should be present
	plain := smplog.StripANSI(out)
	if !strings.Contains(plain, "1)") {
		t.Fatalf("expected '1)' in plain output: %q", plain)
	}
	if !strings.Contains(plain, "2)") {
		t.Fatalf("expected '2)' in plain output: %q", plain)
	}
}

func TestTUIMenuNoColor(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: true,
		Colors: ColorConfig{
			Menu:  smplog.StyleColor256(14),
			Title: smplog.StyleColor256(15),
		},
	})

	out := captureStdout(t, func() {
		NewTUI().Menu(&MenuParams{
			Items: []MenuEntry{
				{Label: "item", Selected: false},
			},
		})
	})

	if strings.Contains(out, "\x1b[") {
		t.Fatalf("expected no ANSI escapes with NoColor=true: %q", out)
	}
	if !strings.Contains(out, "item") {
		t.Fatalf("expected 'item' in output: %q", out)
	}
}

func TestTUIMenuTitlePrefixPostfix(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: true,
		TUI: TUIConfig{
			MenuTitlePrefix:  "[ ",
			MenuTitlePostfix: " ]",
		},
	})

	out := captureStdout(t, func() {
		NewTUI().MenuTitle(&TitleParams{Text: "Main Menu"})
	})

	plain := smplog.StripANSI(out)
	if !strings.Contains(plain, "[ Main Menu ]") {
		t.Fatalf("expected '[ Main Menu ]' in plain output: %q", plain)
	}
}

func TestTUIMenuTitleUsesTitleColor(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: false,
		Colors: ColorConfig{
			Title: smplog.StyleColor256(15),
		},
	})

	out := captureStdout(t, func() {
		NewTUI().MenuTitle(&TitleParams{Text: "Main Menu"})
	})

	if !strings.Contains(out, "\x1b[38;5;15m") {
		t.Fatalf("expected title color escape in output: %q", out)
	}
	if !strings.Contains(out, "Main Menu") {
		t.Fatalf("expected title text in output: %q", out)
	}
}

// ---------- Selector ----------

func TestTUISelectorRendersLabelAndCurrentItem(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: false,
		Colors: ColorConfig{
			Prompt: smplog.StyleColor256(10),
			Data:   smplog.StyleColor256(7),
		},
	})

	out := captureStdout(t, func() {
		NewTUI().Selector(&SelectorParams{
			Label:   "mode",
			Items:   []string{"a", "b", "c"},
			Current: 1, // "b"
		})
	})

	plain := smplog.StripANSI(out)
	if !strings.Contains(plain, "< b >") {
		t.Fatalf("expected '< b >' in plain output: %q", plain)
	}
	if !strings.Contains(plain, "mode") {
		t.Fatalf("expected label in plain output: %q", plain)
	}
	// data color should be present for current item
	if !strings.Contains(out, "\x1b[38;5;7m") {
		t.Fatalf("expected data color in output: %q", out)
	}
}

func TestTUISelectorOutOfBoundsCurrentIsEmpty(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{NoColor: true})

	// Should not panic; current="" when index is out of bounds
	out := captureStdout(t, func() {
		NewTUI().Selector(&SelectorParams{
			Label:   "opt",
			Items:   []string{"x"},
			Current: 99,
		})
	})

	// "< %s >" with empty string gives "< >"-style output (space on each side)
	if !strings.Contains(out, "<") || !strings.Contains(out, ">") {
		t.Fatalf("expected selector brackets in out-of-bounds output: %q", out)
	}
}

// ---------- Input ----------

func TestTUIInputActiveRendersLabelValueCursor(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: false,
		Colors: ColorConfig{
			Prompt: smplog.StyleColor256(10),
			Data:   smplog.StyleColor256(7),
		},
		TUI: TUIConfig{InputCursor: "|"},
	})

	out := captureStdout(t, func() {
		NewTUI().Input(&InputParams{
			Label:  "name",
			Value:  "dan",
			Active: true,
		})
	})

	plain := smplog.StripANSI(out)
	if !strings.Contains(plain, "|") {
		t.Fatalf("expected cursor '|' in plain output: %q", plain)
	}
	if !strings.Contains(plain, "name") {
		t.Fatalf("expected label in plain output: %q", plain)
	}
	if !strings.Contains(plain, "dan") {
		t.Fatalf("expected value in plain output: %q", plain)
	}
}

func TestTUIInputInactiveOmitsCursor(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: true,
		TUI:     TUIConfig{InputCursor: "|"},
	})

	out := captureStdout(t, func() {
		NewTUI().Input(&InputParams{
			Label:  "name",
			Value:  "dan",
			Active: false,
		})
	})

	plain := smplog.StripANSI(out)
	if strings.Contains(plain, "|") {
		t.Fatalf("expected no cursor in inactive input: %q", plain)
	}
}

// ---------- Divider ----------

func TestTUIDividerUsesConfigWidth(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: false,
		Colors: ColorConfig{
			Divider: smplog.StyleColor256(8),
		},
		TUI: TUIConfig{DividerWidth: 40},
	})

	out := captureStdout(t, func() {
		NewTUI().Divider(&DividerParams{})
	})

	plain := strings.TrimRight(smplog.StripANSI(out), "\n")
	if utf8.RuneCountInString(plain) != 40 {
		t.Fatalf("expected divider rune count 40, got %d (%q)", utf8.RuneCountInString(plain), plain)
	}
}

func TestTUIDividerParamWidthOverrides(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: true,
		TUI:     TUIConfig{DividerWidth: 40},
	})

	out := captureStdout(t, func() {
		NewTUI().Divider(&DividerParams{Width: 20})
	})

	plain := strings.TrimRight(smplog.StripANSI(out), "\n")
	if utf8.RuneCountInString(plain) != 20 {
		t.Fatalf("expected divider rune count 20, got %d (%q)", utf8.RuneCountInString(plain), plain)
	}
}

func TestTUIDividerCustomRune(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: true,
		TUI:     TUIConfig{DividerWidth: 10},
	})

	out := captureStdout(t, func() {
		NewTUI().Divider(&DividerParams{Rune: '='})
	})

	plain := strings.TrimRight(smplog.StripANSI(out), "\n")
	if plain != strings.Repeat("=", 10) {
		t.Fatalf("expected '=' repeated 10 times, got %q", plain)
	}
}

// ---------- TreeView ----------

// testNode implements TreeNode for testing.
type testNode struct {
	key    string
	label  string
	parent string
}

func (n testNode) TreeKey() string    { return n.key }
func (n testNode) TreeLabel() string  { return n.label }
func (n testNode) TreeParent() string { return n.parent }

func TestTreeViewFlatList(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{NoColor: true})

	tui := NewTUI()
	out := captureStdout(t, func() {
		entries := tui.TreeView(&TreeViewParams{
			Nodes: []TreeNode{
				testNode{key: "b", label: "Bravo", parent: ""},
				testNode{key: "a", label: "Alpha", parent: ""},
				testNode{key: "c", label: "Charlie", parent: ""},
			},
		})
		if len(entries) != 3 {
			t.Fatalf("expected 3 entries, got %d", len(entries))
		}
		// Sorted by key
		if entries[0].Node.TreeLabel() != "Alpha" {
			t.Fatalf("expected Alpha first, got %s", entries[0].Node.TreeLabel())
		}
	})
	if !strings.Contains(out, "Alpha") || !strings.Contains(out, "Bravo") {
		t.Fatalf("expected labels in output: %q", out)
	}
}

func TestTreeViewNested(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{NoColor: true})

	tui := NewTUI()
	out := captureStdout(t, func() {
		entries := tui.TreeView(&TreeViewParams{
			Nodes: []TreeNode{
				testNode{key: "root", label: "Root", parent: ""},
				testNode{key: "child1", label: "Child 1", parent: "root"},
				testNode{key: "child2", label: "Child 2", parent: "root"},
				testNode{key: "grandchild", label: "Grandchild", parent: "child1"},
			},
		})
		if len(entries) != 4 {
			t.Fatalf("expected 4 entries, got %d", len(entries))
		}
		if entries[0].Depth != 0 {
			t.Fatalf("root depth should be 0, got %d", entries[0].Depth)
		}
		if entries[1].Depth != 1 {
			t.Fatalf("child depth should be 1, got %d", entries[1].Depth)
		}
	})
	if !strings.Contains(out, "├─") || !strings.Contains(out, "└─") {
		t.Fatalf("expected tree connectors in output: %q", out)
	}
}

func TestTreeViewShowIndex(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{NoColor: true})

	tui := NewTUI()
	out := captureStdout(t, func() {
		tui.TreeView(&TreeViewParams{
			Nodes: []TreeNode{
				testNode{key: "a", label: "Alpha", parent: ""},
			},
			ShowIndex: true,
		})
	})
	if !strings.Contains(out, "  0") {
		t.Fatalf("expected index 0 in output: %q", out)
	}
}

func TestTreeViewCenteredBlockAlignment(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: true,
		TUI: TUIConfig{
			MaxWidth: 60,
			Centered: true,
		},
	})

	tui := NewTUI()
	out := captureStdout(t, func() {
		tui.TreeView(&TreeViewParams{
			Nodes: []TreeNode{
				testNode{key: "root", label: "Root", parent: ""},
				testNode{key: "child1", label: "Child 1", parent: "root"},
				testNode{key: "child2", label: "Child 2 has a longer label", parent: "root"},
			},
		})
	})

	// All lines should have the same left margin (leading spaces).
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected multiple lines, got %d", len(lines))
	}
	firstMargin := len(lines[0]) - len(strings.TrimLeft(lines[0], " "))
	for i, line := range lines {
		margin := len(line) - len(strings.TrimLeft(line, " "))
		if margin != firstMargin {
			t.Fatalf("line %d margin %d differs from line 0 margin %d\nlines:\n%s", i, margin, firstMargin, out)
		}
	}
}

// ---------- Summary / PhaseTimer ----------

func TestPhaseTimerOrdering(t *testing.T) {
	pt := NewPhaseTimer()
	pt.Begin("init")
	time.Sleep(5 * time.Millisecond)
	pt.Begin("process")
	time.Sleep(5 * time.Millisecond)
	pt.End()

	phases := pt.Phases()
	if len(phases) != 2 {
		t.Fatalf("expected 2 phases, got %d", len(phases))
	}
	if phases[0].Label != "init" {
		t.Fatalf("expected first phase 'init', got %q", phases[0].Label)
	}
	if phases[1].Label != "process" {
		t.Fatalf("expected second phase 'process', got %q", phases[1].Label)
	}
	if pt.Elapsed() < 10*time.Millisecond {
		t.Fatalf("expected elapsed >= 10ms, got %s", pt.Elapsed())
	}
}

func TestPhaseTimerConcurrency(t *testing.T) {
	pt := NewPhaseTimer()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pt.Begin("phase")
			pt.End()
			_ = pt.Phases()
			_ = pt.Elapsed()
		}()
	}
	wg.Wait()
}

func TestOperationSummaryOK(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: false,
		Colors: ColorConfig{
			Success: smplog.StyleColor256(10),
			Error:   smplog.StyleColor256(9),
			Title:   smplog.StyleColor256(11),
		},
	})

	tui := NewTUI()
	out := captureStdout(t, func() {
		tui.OperationSummary(&OperationSummaryParams{
			Title: "Deploy",
			OK:    true,
			Fields: []SummaryField{
				{Label: "target", Value: "prod"},
			},
		})
	})
	if !strings.Contains(out, "[OK]") {
		t.Fatalf("expected [OK] in output: %q", out)
	}
	// Should use success color (color 10)
	if !strings.Contains(out, "\x1b[38;5;10m") {
		t.Fatalf("expected success color in output: %q", out)
	}
}

func TestOperationSummaryFailed(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: false,
		Colors: ColorConfig{
			Error: smplog.StyleColor256(9),
			Title: smplog.StyleColor256(11),
		},
	})

	tui := NewTUI()
	out := captureStdout(t, func() {
		tui.OperationSummary(&OperationSummaryParams{
			Title: "Deploy",
			OK:    false,
		})
	})
	if !strings.Contains(out, "[FAILED]") {
		t.Fatalf("expected [FAILED] in output: %q", out)
	}
	// Should use error color (color 9)
	if !strings.Contains(out, "\x1b[38;5;9m") {
		t.Fatalf("expected error color in output: %q", out)
	}
}

func TestOperationSummaryCenteredBlockAlignment(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: true,
		TUI: TUIConfig{
			MaxWidth: 60,
			Centered: true,
		},
	})

	tui := NewTUI()
	out := captureStdout(t, func() {
		tui.OperationSummary(&OperationSummaryParams{
			Title: "Deploy",
			OK:    true,
			Fields: []SummaryField{
				{Label: "target", Value: "prod"},
				{Label: "version", Value: "v1.2.3-beta.42"},
			},
		})
	})

	// All lines should have the same total length (padded to blockWidth, then centered).
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected multiple lines, got %d", len(lines))
	}
	firstLen := len(lines[0])
	for i, line := range lines {
		if len(line) != firstLen {
			t.Fatalf("line %d length %d differs from line 0 length %d\nlines:\n%s", i, len(line), firstLen, out)
		}
	}
}

// ---------- Progress ----------

func TestProgressBarRendersOutput(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{NoColor: true})

	var out bytes.Buffer
	var dst bytes.Buffer
	pb := NewProgressBar(&dst, ProgressBarParams{
		Label:     "upload",
		Total:     100,
		Width:     10,
		MinRender: 0, // render every write
		Out:       &out,
	})

	data := make([]byte, 50)
	n, err := pb.Write(data)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	if n != 50 {
		t.Fatalf("Write returned %d, want 50", n)
	}

	pb.Done()

	rendered := out.String()
	if !strings.Contains(rendered, "upload") {
		t.Fatalf("expected label in output: %q", rendered)
	}
	if !strings.Contains(rendered, "[") {
		t.Fatalf("expected bar in output: %q", rendered)
	}
	if dst.Len() != 50 {
		t.Fatalf("dst got %d bytes, want 50", dst.Len())
	}
}

func TestProgressBarThrottling(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{NoColor: true})

	var out bytes.Buffer
	var dst bytes.Buffer
	pb := NewProgressBar(&dst, ProgressBarParams{
		Label:     "dl",
		Total:     1000,
		Width:     10,
		MinRender: 1 * time.Hour, // effectively never re-render
		Out:       &out,
	})

	// First write triggers initial render
	pb.Write(make([]byte, 100))
	firstLen := out.Len()
	if firstLen == 0 {
		t.Fatal("expected initial render")
	}

	// Second write should be throttled
	pb.Write(make([]byte, 100))
	if out.Len() != firstLen {
		t.Fatal("expected throttled render (no new output)")
	}
}

func TestProgressBarNoTotal(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{NoColor: true})

	var out bytes.Buffer
	var dst bytes.Buffer
	pb := NewProgressBar(&dst, ProgressBarParams{
		Label:     "stream",
		Total:     0, // unknown
		Width:     10,
		MinRender: 0,
		Out:       &out,
	})

	pb.Write(make([]byte, 42))
	pb.Done()

	rendered := out.String()
	if !strings.Contains(rendered, "?%") {
		t.Fatalf("expected unknown percentage marker in output: %q", rendered)
	}
}

func TestProgressBarNoColor(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{NoColor: true})

	var out bytes.Buffer
	pb := NewProgressBar(&bytes.Buffer{}, ProgressBarParams{
		Label:     "test",
		Total:     100,
		Width:     10,
		MinRender: 0,
		Out:       &out,
	})

	pb.Write(make([]byte, 50))
	pb.Done()

	if strings.Contains(out.String(), "\x1b[") {
		t.Fatalf("expected no ANSI escapes with NoColor=true: %q", out.String())
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		in   int64
		want string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KiB"},
		{1536, "1.5 KiB"},
		{1048576, "1.0 MiB"},
		{1073741824, "1.0 GiB"},
	}
	for _, tt := range tests {
		got := formatBytes(tt.in)
		if got != tt.want {
			t.Errorf("formatBytes(%d) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// ---------- Shared / Width / Centering ----------

func TestTUIWidthClampsTruncates(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: true,
		Colors:  ColorConfig{Title: smplog.StyleColor256(15)},
	})

	out := captureStdout(t, func() {
		NewTUI().MenuTitle(&TitleParams{Text: "Hello World", Width: 5})
	})

	plain := strings.TrimRight(smplog.StripANSI(out), "\n")
	if utf8.RuneCountInString(plain) != 5 {
		t.Fatalf("expected 5 runes after clipping, got %d (%q)", utf8.RuneCountInString(plain), plain)
	}
	if plain != "Hello" {
		t.Fatalf("expected 'Hello', got %q", plain)
	}
}

func TestTUICenteringPadsContent(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: true,
		TUI:     TUIConfig{MaxWidth: 20, Centered: true},
	})

	out := captureStdout(t, func() {
		NewTUI().MenuTitle(&TitleParams{Text: "Hi"})
	})

	// Strip trailing newline for analysis
	line := strings.TrimRight(smplog.StripANSI(out), "\n")
	total := utf8.RuneCountInString(line)
	if total != 20 {
		t.Fatalf("expected total visible width 20, got %d (%q)", total, line)
	}
	if !strings.HasPrefix(line, " ") {
		t.Fatalf("expected leading spaces for centering: %q", line)
	}
	if !strings.HasSuffix(line, " ") {
		t.Fatalf("expected trailing spaces for centering: %q", line)
	}
}

func TestTUICenteringRequiresMaxWidth(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: true,
		TUI:     TUIConfig{MaxWidth: 0, Centered: true},
	})

	out := captureStdout(t, func() {
		NewTUI().MenuTitle(&TitleParams{Text: "Hi"})
	})

	line := strings.TrimRight(smplog.StripANSI(out), "\n")
	// Without MaxWidth, no padding should be added
	if line != "Hi" {
		t.Fatalf("expected 'Hi' without padding when MaxWidth=0, got %q", line)
	}
}

func TestTUISelectorCenteringPadsContent(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: true,
		TUI:     TUIConfig{MaxWidth: 30, Centered: true},
	})

	out := captureStdout(t, func() {
		NewTUI().Selector(&SelectorParams{
			Label:   "x",
			Items:   []string{"y"},
			Current: 0,
		})
	})

	// "x: < y >" = 9 runes; padded to 30 total
	line := strings.TrimRight(smplog.StripANSI(out), "\n")
	total := utf8.RuneCountInString(line)
	if total != 30 {
		t.Fatalf("expected total visible width 30, got %d (%q)", total, line)
	}
	if !strings.HasPrefix(line, " ") {
		t.Fatalf("expected leading spaces for centering: %q", line)
	}
}

func TestTUIInputCenteringPadsContent(t *testing.T) {
	orig := Configured()
	t.Cleanup(func() { Configure(orig) })
	Configure(Config{
		NoColor: true,
		TUI:     TUIConfig{MaxWidth: 30, Centered: true, InputCursor: "|"},
	})

	out := captureStdout(t, func() {
		NewTUI().Input(&InputParams{
			Label:  "name",
			Value:  "dan",
			Active: true,
		})
	})

	// "name: dan|" = 10 runes; padded to 30 total
	line := strings.TrimRight(smplog.StripANSI(out), "\n")
	total := utf8.RuneCountInString(line)
	if total != 30 {
		t.Fatalf("expected total visible width 30, got %d (%q)", total, line)
	}
	if !strings.HasPrefix(line, " ") {
		t.Fatalf("expected leading spaces for centering: %q", line)
	}
}

func TestTUIRefreshWritesClearAndMoveTo(t *testing.T) {
	out := captureStdout(t, func() {
		if err := NewTUI().Refresh(); err != nil {
			t.Fatalf("refresh: %v", err)
		}
	})

	if !strings.Contains(out, "\x1b[2J") {
		t.Fatalf("expected clear screen sequence in output: %q", out)
	}
	if !strings.Contains(out, "\x1b[1;1H") {
		t.Fatalf("expected move-to-1,1 sequence in output: %q", out)
	}
}
