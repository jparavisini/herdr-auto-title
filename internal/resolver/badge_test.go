package resolver

import (
	"strings"
	"testing"

	"github.com/rivo/uniseg"

	"github.com/kryptamine/herdr-auto-title/internal/state"
)

func badgedTab(badges ...string) state.TabState {
	panes := make([]*state.PaneState, 0, len(badges))
	for i, badge := range badges {
		panes = append(
			panes,
			&state.PaneState{
				ID:      "wE:p" + string(rune('1'+i)),
				Dir:     dashboard,
				Focused: i == 0,
				Badge:   badge,
			},
		)
	}

	tab := tabOf(panes)
	tab.Position = 4

	return tab
}

func TestBadgeLeadsThePosition(t *testing.T) {
	t.Parallel()

	r := NewBadged(numberedCWD(DefaultMaxLength), DefaultMaxLength)

	tests := []struct {
		name   string
		badges []string
		want   string
	}{
		{"no pane has a badge", []string{"", ""}, "4 · dashboard"},
		{"one pane has a badge", []string{"", "✓"}, "✓ 4 · dashboard"},
		{"the first badge wins", []string{"✓", "!"}, "✓ 4 · dashboard"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := r.Resolve(badgedTab(tc.badges...)); got.Name != tc.want {
				t.Errorf("name = %q, want %q", got.Name, tc.want)
			}
		})
	}
}

func TestBadgeCountsAgainstTheBound(t *testing.T) {
	t.Parallel()

	const maxLength = 10

	got := NewBadged(numberedCWD(maxLength), maxLength).Resolve(badgedTab("✓"))
	if !strings.HasPrefix(got.Name, "✓ 4 · ") || uniseg.StringWidth(got.Name) > maxLength {
		t.Errorf("name = %q, want the badge and position within %d columns", got.Name, maxLength)
	}
}
