package tui

import (
	"fmt"
	"sort"
	"unicode/utf8"

	smplog "github.com/danmuck/smplog"
)

// TreeNode is the interface for items that can be displayed in a tree.
type TreeNode interface {
	TreeLabel() string
	TreeKey() string
	TreeParent() string // empty = root
}

// TreeViewEntry is a flattened tree row with computed prefix and depth.
type TreeViewEntry struct {
	Node   TreeNode
	Prefix string // e.g. "├─ " or "└─ "
	Depth  int
	Index  int
}

// TreeViewParams configures TUI.TreeViewTC.
type TreeViewParams struct {
	Nodes     []TreeNode
	Width     int
	ShowIndex bool
}

// TreeViewTC renders a tree and returns the flattened entry list.
// Nodes are grouped by TreeParent(); siblings are sorted by TreeKey().
func (t TUI) TreeViewTC(p *TreeViewParams) []TreeViewEntry {
	cfg := Configured()
	width := effectiveWidth(p.Width, cfg)

	// Build parent → children map
	children := make(map[string][]TreeNode)
	for _, n := range p.Nodes {
		parent := n.TreeParent()
		children[parent] = append(children[parent], n)
	}
	// Sort each group by key
	for k := range children {
		sort.Slice(children[k], func(i, j int) bool {
			return children[k][i].TreeKey() < children[k][j].TreeKey()
		})
	}

	var entries []TreeViewEntry
	var walk func(parentKey string, depth int, parentPrefix string)
	walk = func(parentKey string, depth int, parentPrefix string) {
		kids := children[parentKey]
		for i, node := range kids {
			isLast := i == len(kids)-1
			var connector, childPrefix string
			if depth == 0 {
				connector = ""
				childPrefix = ""
			} else {
				if isLast {
					connector = parentPrefix + "└─ "
					childPrefix = parentPrefix + "   "
				} else {
					connector = parentPrefix + "├─ "
					childPrefix = parentPrefix + "│  "
				}
			}
			entry := TreeViewEntry{
				Node:   node,
				Prefix: connector,
				Depth:  depth,
				Index:  len(entries),
			}
			entries = append(entries, entry)
			walk(node.TreeKey(), depth+1, childPrefix)
		}
	}
	walk("", 0, "")

	// Render as a block so all lines share the same left margin when centered.
	lines := make([]blockLine, len(entries))
	for i, e := range entries {
		var plain string
		if p.ShowIndex {
			plain = fmt.Sprintf("%3d  %s%s", e.Index, e.Prefix, e.Node.TreeLabel())
		} else {
			plain = e.Prefix + e.Node.TreeLabel()
		}
		if width > 0 {
			plain = Clip(width, plain)
		}
		colored := smplog.Colorize(cfg.Colors.Data, plain, cfg.NoColor)
		lines[i] = blockLine{colored: colored, plainWidth: utf8.RuneCountInString(plain)}
	}
	t.writeBlock(cfg, lines)

	return entries
}
