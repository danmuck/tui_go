package tui

import (
	"fmt"
	"sync"
	"time"
	"unicode/utf8"

	smplog "github.com/danmuck/smplog"
)

// PhaseRecord stores timing data for a single phase.
type PhaseRecord struct {
	Label   string
	Elapsed time.Duration
}

// PhaseTimer tracks sequential phases of an operation with thread safety.
type PhaseTimer struct {
	mu         sync.Mutex
	phases     []PhaseRecord
	current    string
	phaseStart time.Time
	startTime  time.Time
}

// NewPhaseTimer creates a new PhaseTimer. The clock starts now.
func NewPhaseTimer() *PhaseTimer {
	return &PhaseTimer{startTime: time.Now()}
}

// Begin ends the previous phase (if any) and starts a new one.
func (pt *PhaseTimer) Begin(label string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	now := time.Now()
	if pt.current != "" {
		pt.phases = append(pt.phases, PhaseRecord{
			Label:   pt.current,
			Elapsed: now.Sub(pt.phaseStart),
		})
	}
	pt.current = label
	pt.phaseStart = now
}

// End closes the current phase without starting a new one.
func (pt *PhaseTimer) End() {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	now := time.Now()
	if pt.current != "" {
		pt.phases = append(pt.phases, PhaseRecord{
			Label:   pt.current,
			Elapsed: now.Sub(pt.phaseStart),
		})
		pt.current = ""
	}
}

// Phases returns a copy of all completed phase records.
func (pt *PhaseTimer) Phases() []PhaseRecord {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	out := make([]PhaseRecord, len(pt.phases))
	copy(out, pt.phases)
	return out
}

// Elapsed returns the total time since the timer was created.
func (pt *PhaseTimer) Elapsed() time.Duration {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	return time.Since(pt.startTime)
}

// SummaryField is a key/value pair rendered in an OperationSummary.
type SummaryField struct {
	Label string
	Value string
}

// OperationSummaryParams configures TUI.OperationSummaryTC.
type OperationSummaryParams struct {
	Title  string
	OK     bool
	Fields []SummaryField
	Timer  *PhaseTimer // optional
	Width  int
}

// OperationSummaryTC renders a titled summary block with status tag, fields,
// and optional phase timing breakdown.
func (t TUI) OperationSummaryTC(p *OperationSummaryParams) {
	cfg := Configured()

	// Title line: "  Title  [OK]" or "  Title  [FAILED]"
	var tag, tagColor string
	if p.OK {
		tag = "[OK]"
		tagColor = cfg.Colors.Success
	} else {
		tag = "[FAILED]"
		tagColor = cfg.Colors.Error
	}

	// Collect all lines into a block for consistent centering.
	var lines []blockLine

	// Title line
	titlePlain := p.Title + "  " + tag
	titleColored := smplog.Colorize(cfg.Colors.Title, p.Title+"  ", cfg.NoColor) +
		smplog.Colorize(tagColor, tag, cfg.NoColor)
	titleWidth := utf8.RuneCountInString(titlePlain)
	lines = append(lines, blockLine{colored: titleColored, plainWidth: titleWidth})

	// Fields
	for _, f := range p.Fields {
		plain := fmt.Sprintf("  %s: %s", f.Label, f.Value)
		labelText := smplog.Colorize(cfg.Colors.Prompt, "  "+f.Label, cfg.NoColor)
		valueText := smplog.Colorize(cfg.Colors.Data, f.Value, cfg.NoColor)
		line := fmt.Sprintf("%s: %s", labelText, valueText)
		plainWidth := utf8.RuneCountInString(plain)
		lines = append(lines, blockLine{colored: line, plainWidth: plainWidth})
	}

	// Phase breakdown
	if p.Timer != nil {
		phases := p.Timer.Phases()
		if len(phases) > 0 {
			// Blank separator line
			lines = append(lines, blockLine{colored: "", plainWidth: 0})
			// Header
			headerPlain := "  Phases:"
			headerColored := smplog.Colorize(cfg.Colors.Title, headerPlain, cfg.NoColor)
			lines = append(lines, blockLine{colored: headerColored, plainWidth: utf8.RuneCountInString(headerPlain)})
			// Phase rows
			for _, ph := range phases {
				plain := fmt.Sprintf("    %s  %s", ph.Label, ph.Elapsed.Round(time.Millisecond))
				labelText := smplog.Colorize(cfg.Colors.Prompt, "    "+ph.Label, cfg.NoColor)
				durText := smplog.Colorize(cfg.Colors.Data, "  "+ph.Elapsed.Round(time.Millisecond).String(), cfg.NoColor)
				line := labelText + durText
				plainWidth := utf8.RuneCountInString(plain)
				lines = append(lines, blockLine{colored: line, plainWidth: plainWidth})
			}
			// Total
			totalPlain := fmt.Sprintf("    Total  %s", p.Timer.Elapsed().Round(time.Millisecond))
			totalColored := smplog.Colorize(cfg.Colors.Data, totalPlain, cfg.NoColor)
			lines = append(lines, blockLine{colored: totalColored, plainWidth: utf8.RuneCountInString(totalPlain)})
		}
	}

	t.writeBlock(cfg, lines)
}
