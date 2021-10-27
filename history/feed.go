package history

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/bouncepaw/mycorrhiza/cfg"

	"github.com/gorilla/feeds"
)

const changeGroupMaxSize = 100

func recentChangesFeed(opts FeedOptions) *feeds.Feed {
	feed := &feeds.Feed{
		Title:       cfg.WikiName + " (recent changes)",
		Link:        &feeds.Link{Href: cfg.URL},
		Description: fmt.Sprintf("List of %d recent changes on the wiki", changeGroupMaxSize),
		Updated:     time.Now(),
	}
	revs := newRecentChangesStream()
	groups := opts.grouping.Group(revs)
	for _, grp := range groups {
		item := grp.feedItem(opts)
		feed.Add(&item)
	}
	return feed
}

// RecentChangesRSS creates recent changes feed in RSS format.
func RecentChangesRSS(opts FeedOptions) (string, error) {
	return recentChangesFeed(opts).ToRss()
}

// RecentChangesAtom creates recent changes feed in Atom format.
func RecentChangesAtom(opts FeedOptions) (string, error) {
	return recentChangesFeed(opts).ToAtom()
}

// RecentChangesJSON creates recent changes feed in JSON format.
func RecentChangesJSON(opts FeedOptions) (string, error) {
	return recentChangesFeed(opts).ToJSON()
}

// revisionGroup is a slice of revisions, ordered most recent first.
type revisionGroup []Revision

func newRevisionGroup(rev Revision) revisionGroup {
	return revisionGroup([]Revision{rev})
}

func (grp *revisionGroup) addRevision(rev Revision) {
	*grp = append(*grp, rev)
}

func groupRevisionsByMonth(revs []Revision) (res []revisionGroup) {
	var (
		currentYear  int
		currentMonth time.Month
	)
	for _, rev := range revs {
		if rev.Time.Month() != currentMonth || rev.Time.Year() != currentYear {
			currentYear = rev.Time.Year()
			currentMonth = rev.Time.Month()
			res = append(res, newRevisionGroup(rev))
		} else {
			res[len(res)-1].addRevision(rev)
		}
	}
	return res
}

// groupRevisionsByPeriodFromNow groups close-together revisions and returns the first changeGroupMaxSize (30) groups.
// If two revisions happened within period of each other, they are put in the same group.
func groupRevisionsByPeriod(revs recentChangesStream, period time.Duration) (res []revisionGroup) {
	nextRev := revs.iterator()
	rev, empty := nextRev()
	if empty {
		return res
	}

	currTime := rev.Time
	currGroup := newRevisionGroup(rev)
	for {
		rev, done := nextRev()
		if done {
			return append(res, currGroup)
		}

		if currTime.Sub(rev.Time) < period && currGroup[0].Username == rev.Username {
			currGroup.addRevision(rev)
		} else {
			res = append(res, currGroup)
			if len(res) == changeGroupMaxSize {
				return res
			}
			currGroup = newRevisionGroup(rev)
		}
		currTime = rev.Time
	}
}

func (grp revisionGroup) feedItem(opts FeedOptions) feeds.Item {
	return feeds.Item{
		Title: grp.title(opts.groupOrder),
		// groups for feeds should have the same author for all revisions
		Author:      &feeds.Author{Name: grp[0].Username},
		Id:          grp[len(grp)-1].Hash,
		Description: grp.descriptionForFeed(opts.groupOrder),
		Created:     grp[len(grp)-1].Time, // earliest revision
		Updated:     grp[0].Time,          // latest revision
		Link:        &feeds.Link{Href: cfg.URL + grp[0].bestLink()},
	}
}

func (grp revisionGroup) title(order FeedGroupOrder) string {
	var message string
	switch order {
	case NewToOld:
		message = grp[0].Message
	case OldToNew:
		message = grp[len(grp)-1].Message
	}

	author := grp[0].Username
	if len(grp) == 1 {
		return fmt.Sprintf("%s by %s", message, author)
	} else {
		return fmt.Sprintf("%d edits by %s (%s, ...)", len(grp), author, message)
	}
}

func (grp revisionGroup) descriptionForFeed(order FeedGroupOrder) string {
	builder := strings.Builder{}
	switch order {
	case NewToOld:
		for _, rev := range grp {
			builder.WriteString(rev.descriptionForFeed())
		}
	case OldToNew:
		for i := len(grp) - 1; i >= 0; i-- {
			builder.WriteString(grp[i].descriptionForFeed())
		}
	}
	return builder.String()
}

type FeedOptions struct {
	grouping   FeedGrouping
	groupOrder FeedGroupOrder
}

func ParseFeedOptions(query url.Values) (FeedOptions, error) {
	grouping, err := parseFeedGrouping(query)
	if err != nil {
		return FeedOptions{}, err
	}
	groupOrder, err := parseFeedGroupOrder(query)
	if err != nil {
		return FeedOptions{}, err
	}
	return FeedOptions{grouping, groupOrder}, nil
}

type FeedGrouping interface {
	Group(recentChangesStream) []revisionGroup
}

func parseFeedGrouping(query url.Values) (FeedGrouping, error) {
	if query.Get("period") == "" {
		return NormalFeedGrouping{}, nil
	} else {
		period, err := time.ParseDuration(query.Get("period"))
		if err != nil {
			return nil, err
		}
		return PeriodFeedGrouping{Period: period}, nil
	}
}

type NormalFeedGrouping struct{}

func (NormalFeedGrouping) Group(revs recentChangesStream) (res []revisionGroup) {
	for _, rev := range revs.next(changeGroupMaxSize) {
		res = append(res, newRevisionGroup(rev))
	}
	return res
}

type PeriodFeedGrouping struct {
	Period time.Duration
}

func (g PeriodFeedGrouping) Group(revs recentChangesStream) (res []revisionGroup) {
	return groupRevisionsByPeriod(revs, g.Period)
}

type FeedGroupOrder int

const (
	NewToOld FeedGroupOrder = iota
	OldToNew FeedGroupOrder = iota
)

func parseFeedGroupOrder(query url.Values) (FeedGroupOrder, error) {
	switch query.Get("order") {
	case "", "old-to-new":
		return OldToNew, nil
	case "new-to-old":
		return NewToOld, nil
	}
	return 0, errors.New("unknown order")
}
