package gamenotation

import (
	"fmt"
)

// Snapshot is a representation of a game's snapshot in time using Board Game Notation (BGN)
type Snapshot struct {
	Tags    Tags     `json:"tags"`
	Actions []Action `json:"actions"`
}

func (s *Snapshot) String() string {
	qgn := ""
	for key, value := range s.Tags {
		qgn += fmt.Sprintf("[%s \"%s\"]\n", key, value)
	}
	qgn += "\n"
	line := ""
	for _, action := range s.Actions {
		line += fmt.Sprintf("%s ", action.String())
		if len(line) > 70 {
			qgn += fmt.Sprintf("%s\n", line[:len(line)-1])
			line = ""
		}
	}
	if line != "" {
		qgn += line[:len(line)-1]
	}
	return qgn
}

type Action struct {
	Index   int      `json:"index"`   // the index of team in Teams Tag
	Key     string   `json:"key"`     // single character key indicating the action to perform
	Details []string `json:"details"` // additional details that can be optionally used when describing an action
}

func (a *Action) String() string {
	qgn := fmt.Sprintf("%d%s", a.Index, string(a.Key))
	details := ""
	for _, detail := range a.Details {
		details += fmt.Sprintf("%s.", detail)
	}
	if len(details) > 0 {
		details = details[:len(details)-1]
		qgn = fmt.Sprintf("%s&%s", qgn, details)
	}
	return qgn
}
