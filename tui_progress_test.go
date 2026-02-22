package tui

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

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
