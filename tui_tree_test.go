package tui

import (
	"strings"
	"sync"
	"testing"
	"time"

	smplog "github.com/danmuck/smplog"
)

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
	// Content-level indentation (e.g. "  target:") means leading-space checks won't work,
	// but equal line lengths prove they share the same centering block.
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
