package history

import (
	"fmt"
	"strings"
	"time"

	"github.com/bouncepaw/mycorrhiza/cfg"

	"github.com/gorilla/feeds"
)

var groupPeriod, _ = time.ParseDuration("30m")

func recentChangesFeed() *feeds.Feed {
	feed := &feeds.Feed{
		Title:       "Recent changes",
		Link:        &feeds.Link{Href: cfg.URL},
		Description: "List of 30 recent changes on the wiki",
		Author:      &feeds.Author{Name: "Wikimind", Email: "wikimind@mycorrhiza"},
		Updated:     time.Now(),
	}
	revs := RecentChanges(30)
	groups := groupRevisionsByPeriod(revs, groupPeriod)
	for _, grp := range groups {
		item := grp.feedItem()
		feed.Add(&item)
	}
	return feed
}

// RecentChangesRSS creates recent changes feed in RSS format.
func RecentChangesRSS() (string, error) {
	return recentChangesFeed().ToRss()
}

// RecentChangesAtom creates recent changes feed in Atom format.
func RecentChangesAtom() (string, error) {
	return recentChangesFeed().ToAtom()
}

// RecentChangesJSON creates recent changes feed in JSON format.
func RecentChangesJSON() (string, error) {
	return recentChangesFeed().ToJSON()
}

func (grp revisionGroup) feedItem() feeds.Item {
	return feeds.Item{
		Title:       grp.title(),
		Author:      grp.author(),
		Id:          grp[0].Hash,
		Description: grp.descriptionForFeed(),
		Created:     grp[len(grp)-1].Time, // earliest revision
		Updated:     grp[0].Time,          // latest revision
		Link:        &feeds.Link{Href: cfg.URL + grp[0].bestLink()},
	}
}

func (grp revisionGroup) title() string {
	if len(grp) == 1 {
		return grp[0].Message
	} else {
		return fmt.Sprintf("%d edits (%s, ...)", len(grp), grp[0].Message)
	}
}

func (grp revisionGroup) author() *feeds.Author {
	author := grp[0].Username
	for _, rev := range grp[1:] {
		// if they don't all have the same author, return nil
		if rev.Username != author {
			return nil
		}
	}
	return &feeds.Author{Name: author}
}

func (grp revisionGroup) descriptionForFeed() string {
	builder := strings.Builder{}
	for _, rev := range grp {
		builder.WriteString(rev.descriptionForFeed())
	}
	return builder.String()
}
