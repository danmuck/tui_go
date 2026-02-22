package tui

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	smplog "github.com/danmuck/smplog"
)

// ProgressBarParams configures a new ProgressBar.
type ProgressBarParams struct {
	Label     string        // label printed before the bar
	Total     int64         // expected total bytes; 0 = unknown (no percentage)
	Width     int           // bar character width; default 30
	MinRender time.Duration // minimum interval between renders; default 100ms
	Out       io.Writer     // destination for progress output; default os.Stderr
}

// ProgressBar is a stateful io.Writer wrapper that renders a progress bar
// to Out (default stderr) using carriage-return overwrite.
type ProgressBar struct {
	dst io.Writer // wrapped writer that receives Write data

	mu         sync.Mutex
	label      string
	total      int64
	width      int
	minRender  time.Duration
	out        io.Writer
	written    int64
	lastRender time.Time
	startTime  time.Time
	noColor    bool
}

// NewProgressBar creates a ProgressBar wrapping dst. Writes to the returned
// ProgressBar are forwarded to dst while rendering progress to p.Out.
func NewProgressBar(dst io.Writer, p ProgressBarParams) *ProgressBar {
	width := p.Width
	if width <= 0 {
		width = 30
	}
	minRender := p.MinRender
	if minRender <= 0 {
		minRender = 100 * time.Millisecond
	}
	out := p.Out
	if out == nil {
		out = io.Discard
	}
	cfg := Configured()
	return &ProgressBar{
		dst:       dst,
		label:     p.Label,
		total:     p.Total,
		width:     width,
		minRender: minRender,
		out:       out,
		startTime: time.Now(),
		noColor:   cfg.NoColor,
	}
}

// Write forwards data to the underlying writer and updates the progress bar.
func (pb *ProgressBar) Write(data []byte) (int, error) {
	n, err := pb.dst.Write(data)
	pb.mu.Lock()
	pb.written += int64(n)
	now := time.Now()
	shouldRender := now.Sub(pb.lastRender) >= pb.minRender
	pb.mu.Unlock()
	if shouldRender {
		pb.render(now)
	}
	return n, err
}

func (pb *ProgressBar) render(now time.Time) {
	pb.mu.Lock()
	pb.lastRender = now
	written := pb.written
	total := pb.total
	elapsed := now.Sub(pb.startTime)
	label := pb.label
	width := pb.width
	noColor := pb.noColor
	pb.mu.Unlock()

	cfg := Configured()

	var bar string
	var pctStr string
	if total > 0 {
		frac := float64(written) / float64(total)
		if frac > 1 {
			frac = 1
		}
		filled := int(frac * float64(width))
		bar = strings.Repeat("=", filled) + strings.Repeat("-", width-filled)
		pctStr = fmt.Sprintf("%5.1f%%", frac*100)
	} else {
		// unknown total — show spinner-style fill
		filled := int(written) % width
		bar = strings.Repeat("=", filled) + strings.Repeat("-", width-filled)
		pctStr = "   ?%"
	}

	var speed string
	if elapsed > 0 {
		bps := float64(written) / elapsed.Seconds()
		speed = formatBytes(int64(bps)) + "/s"
	}

	var sizeStr string
	if total > 0 {
		sizeStr = fmt.Sprintf("%s / %s", formatBytes(written), formatBytes(total))
	} else {
		sizeStr = formatBytes(written)
	}

	plain := fmt.Sprintf("\r  %s  [%s]  %s  %s  %s", label, bar, pctStr, sizeStr, speed)
	colored := smplog.Colorize(cfg.Colors.Data, plain, noColor)
	fmt.Fprint(pb.out, colored) //nolint:errcheck
}

// Done flushes the final render with a trailing newline.
func (pb *ProgressBar) Done() {
	pb.render(time.Now())
	fmt.Fprintln(pb.out) //nolint:errcheck
}

// Written returns the total number of bytes written so far.
func (pb *ProgressBar) Written() int64 {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	return pb.written
}

// formatBytes formats a byte count into a human-readable string.
func formatBytes(b int64) string {
	const (
		kib = 1024
		mib = 1024 * kib
		gib = 1024 * mib
	)
	switch {
	case b >= gib:
		return fmt.Sprintf("%.1f GiB", float64(b)/float64(gib))
	case b >= mib:
		return fmt.Sprintf("%.1f MiB", float64(b)/float64(mib))
	case b >= kib:
		return fmt.Sprintf("%.1f KiB", float64(b)/float64(kib))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
