package resolver

import (
	"github.com/kryptamine/herdr-auto-title/internal/state"
)

// Badged puts a mark another plugin reported in front of a tab's title, ahead
// of the position: `✓ 4 · dashboard › claude`. It wraps a resolver for the
// reason Numbered does: the mark says nothing about what the tab holds. The
// mark is part of the label Auto Title sets, so a tab wearing it is never read
// as renamed by hand, and the mark goes when the token does.
type Badged struct {
	inner TitleResolver
	// maxLength is the bound inner was built with; the badge is counted
	// against it, as Numbered counts the position.
	maxLength int
}

var _ TitleResolver = (*Badged)(nil)

// NewBadged wraps inner, which must have been given the same maxLength.
// Zero or less takes the default, as New does.
func NewBadged(inner TitleResolver, maxLength int) *Badged {
	if maxLength <= 0 {
		maxLength = DefaultMaxLength
	}

	return &Badged{inner: inner, maxLength: maxLength}
}

// Resolve names the tab and puts the first badge any of its panes carries in
// front. Panes are ordered by ID, so the same session always picks the same one.
func (b *Badged) Resolve(tab state.TabState) Decision {
	decision := b.inner.Resolve(tab)

	for _, pane := range tab.Panes {
		if pane.Badge != "" {
			decision.Name = withPrefix(pane.Badge+" ", decision.Name, b.maxLength)

			break
		}
	}

	return decision
}
