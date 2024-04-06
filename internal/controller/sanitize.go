package controller

import (
	"strings"

	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

var colors = []string{
	"red", "blue", "green", "yellow", "orange", "pink", "purple", "teal",
}

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
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(key), " ", "-"), ".", "-")
}

func sanitizeID(id string) string {
	max := 16
	if len(id) > max {
		id = id[:max]
	}
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(id), " ", "-"), ".", "-")
}
