package history

import (
	"errors"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/feeds"
)

const changeGroupMaxSize = 30

func recentChangesFeed(opts FeedOptions) *feeds.Feed {
	feed := &feeds.Feed{
		Title:       cfg.WikiName + " (recent changes)",
		Link:        &feeds.Link{Href: cfg.URL},
		Description: fmt.Sprintf("List of %d recent changes on the wiki", changeGroupMaxSize),
		Updated:     time.Now(),
	}
	revs := newRecentChangesStream()
	groups := groupRevisions(revs, opts)
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
	return []Revision{rev}
}

func (grp *revisionGroup) addRevision(rev Revision) {
	*grp = append(*grp, rev)
}

// orderedIndex returns the ith revision in the group following the given order.
func (grp *revisionGroup) orderedIndex(i int, order feedGroupOrder) *Revision {
	switch order {
	case newToOld:
		return &(*grp)[i]
	case oldToNew:
		return &(*grp)[len(*grp)-1-i]
	}
	// unreachable
	return nil
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

// groupRevisions groups revisions for a feed.
// It returns the first changeGroupMaxSize (30) groups.
// The grouping parameter determines when two revisions will be grouped.
func groupRevisions(revs recentChangesStream, opts FeedOptions) (res []revisionGroup) {
	nextRev := revs.iterator()
	rev, empty := nextRev()
	if empty {
		return res
	}

	currGroup := newRevisionGroup(rev)
	for rev, done := nextRev(); !done; rev, done = nextRev() {
		if opts.canGroup(currGroup, rev) {
			currGroup.addRevision(rev)
		} else {
			res = append(res, currGroup)
			if len(res) == changeGroupMaxSize {
				return res
			}
			currGroup = newRevisionGroup(rev)
		}
	}
	// no more revisions, haven't added the last group yet
	return append(res, currGroup)
}

func (grp revisionGroup) feedItem(opts FeedOptions) feeds.Item {
	title, author := grp.titleAndAuthor(opts.order)
	return feeds.Item{
		Title:       title,
		Author:      author,
		Id:          grp[len(grp)-1].Hash,
		Description: grp.descriptionForFeed(opts.order),
		Created:     grp[len(grp)-1].Time, // earliest revision
		Updated:     grp[0].Time,          // latest revision
		Link:        &feeds.Link{Href: cfg.URL + grp[0].bestLink()},
		Content:     grp.descriptionForFeed(opts.order),
	}
}

// titleAndAuthor creates a title and author for a feed item.
// If all messages and authors are the same (or there's just one rev), "message by author"
// If all authors are the same, "num edits (first message, ...) by author"
// Else (even if all messages are the same), "num edits (first message, ...)"
func (grp revisionGroup) titleAndAuthor(order feedGroupOrder) (title string, author *feeds.Author) {
	allMessagesSame := true
	allAuthorsSame := true
	for _, rev := range grp[1:] {
		if rev.Message != grp[0].Message {
			allMessagesSame = false
		}
		if rev.Username != grp[0].Username {
			allAuthorsSame = false
		}
		if !allMessagesSame && !allAuthorsSame {
			break
		}
	}

	if allMessagesSame && allAuthorsSame {
		title = grp[0].Message
	} else {
		title = fmt.Sprintf("%d edits (%s, ...)", len(grp), grp.orderedIndex(0, order).Message)
	}

	if allAuthorsSame {
		title += fmt.Sprintf(" by %s", grp[0].Username)
		author = &feeds.Author{Name: grp[0].Username}
	} else {
		author = nil
	}

	return title, author
}

func (grp revisionGroup) descriptionForFeed(order feedGroupOrder) string {
	builder := strings.Builder{}
	for i := 0; i < len(grp); i++ {
		desc := grp.orderedIndex(i, order).descriptionForFeed()
		builder.WriteString(desc)
	}
	return builder.String()
}

type feedOptionParserState struct {
	isAnythingSet bool
	conds         []groupingCondition
	order         feedGroupOrder
}

// feedGrouping represents a set of conditions that must all be satisfied for revisions to be grouped.
// If there are no conditions, revisions will never be grouped.
type FeedOptions struct {
	conds []groupingCondition
	order feedGroupOrder
}

func ParseFeedOptions(query url.Values) (FeedOptions, error) {
	parser := feedOptionParserState{}

	err := parser.parseFeedGroupingPeriod(query)
	if err != nil {
		return FeedOptions{}, err
	}
	err = parser.parseFeedGroupingSame(query)
	if err != nil {
		return FeedOptions{}, err
	}
	err = parser.parseFeedGroupingOrder(query)
	if err != nil {
		return FeedOptions{}, err
	}

	var conds []groupingCondition
	if parser.isAnythingSet {
		conds = parser.conds
	} else {
		// if no options are applied, do no grouping instead of using the default options
		conds = nil
	}
	return FeedOptions{conds: conds, order: parser.order}, nil
}

func (parser *feedOptionParserState) parseFeedGroupingPeriod(query url.Values) error {
	if query["period"] != nil {
		parser.isAnythingSet = true
		period, err := time.ParseDuration(query.Get("period"))
		if err != nil {
			return err
		}
		parser.conds = append(parser.conds, periodGroupingCondition{period})
	}
	return nil
}

func (parser *feedOptionParserState) parseFeedGroupingSame(query url.Values) error {
	if same := query["same"]; same != nil {
		parser.isAnythingSet = true
		if len(same) == 1 && same[0] == "none" {
			// same=none adds no condition
			parser.conds = append(parser.conds, sameGroupingCondition{})
			return nil
		} else {
			// handle same=author, same=author&same=message, etc.
			cond := sameGroupingCondition{}
			for _, sameCond := range same {
				switch sameCond {
				case "author":
					if cond.author {
						return errors.New("set same=author twice")
					}
					cond.author = true
				case "message":
					if cond.message {
						return errors.New("set same=message twice")
					}
					cond.message = true
				default:
					return errors.New("unknown same option " + sameCond)
				}
			}
			parser.conds = append(parser.conds, cond)
			return nil
		}
	} else {
		// same defaults to both author and message
		// but this won't be applied if no grouping options are set
		parser.conds = append(parser.conds, sameGroupingCondition{author: true, message: true})
		return nil
	}
}

type feedGroupOrder int

const (
	newToOld feedGroupOrder = iota
	oldToNew feedGroupOrder = iota
)

func (parser *feedOptionParserState) parseFeedGroupingOrder(query url.Values) error {
	if order := query["order"]; order != nil {
		parser.isAnythingSet = true
		switch query.Get("order") {
		case "old-to-new":
			parser.order = oldToNew
		case "new-to-old":
			parser.order = newToOld
		default:
			return errors.New("unknown order option " + query.Get("order"))
		}
	} else {
		parser.order = oldToNew
	}
	return nil
}

// canGroup determines whether a revision can be added to a group.
func (opts FeedOptions) canGroup(grp revisionGroup, rev Revision) bool {
	if len(opts.conds) == 0 {
		return false
	}

	for _, cond := range opts.conds {
		if !cond.canGroup(grp, rev) {
			return false
		}
	}
	return true
}

type groupingCondition interface {
	canGroup(grp revisionGroup, rev Revision) bool
}

// periodGroupingCondition will group two revisions if they are within period of each other.
type periodGroupingCondition struct {
	period time.Duration
}

func (cond periodGroupingCondition) canGroup(grp revisionGroup, rev Revision) bool {
	return grp[len(grp)-1].Time.Sub(rev.Time) < cond.period
}

type sameGroupingCondition struct {
	author  bool
	message bool
}

func (c sameGroupingCondition) canGroup(grp revisionGroup, rev Revision) bool {
	return (!c.author || grp[0].Username == rev.Username) &&
		(!c.message || grp[0].Message == rev.Message)
}
