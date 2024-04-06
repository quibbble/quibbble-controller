package controller

import (
	"regexp"
	"strings"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

var colors = []string{
	"red", "blue", "green", "yellow", "orange", "pink", "purple", "teal",
}

var allowed = regexp.MustCompile(`[^a-z0-9-]+`) // a-z, 0-9, and -

func sanitizeSnapshot(snapshot *qgn.Snapshot) {
	snapshot.Tags[qgn.KeyTag] = sanitizeKey(snapshot.Tags[qgn.KeyTag])
	snapshot.Tags[qgn.IDTag] = sanitizeID(snapshot.Tags[qgn.IDTag])
	teams := make([]string, 0)
	t, _ := snapshot.Tags.Teams()
	for i := range len(t) {
		teams = append(teams, colors[i])
	}
	snapshot.Tags[qgn.TeamsTag] = strings.Join(teams, ", ")
}

func sanitizeKey(key string) string {
	return allowed.ReplaceAllString(key, "")
}

func sanitizeID(id string) string {
	max := 16
	if len(id) > max {
		id = id[:max]
	}
	return allowed.ReplaceAllString(id, "")
}
