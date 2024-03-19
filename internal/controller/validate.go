package controller

import (
	"fmt"

	"github.com/quibbble/quibbble-controller/games/tictactoe"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

var builders = map[string]qg.GameBuilder{
	tictactoe.Builder{}.GetInformation().Key: tictactoe.Builder{},
}

func validateSnapshot(snapshot *qgn.Snapshot) error {
	key := snapshot.Tags[qgn.KeyTag]
	if err := validateKey(key); err != nil {
		return err
	}
	teams, err := snapshot.Tags.Teams()
	if err != nil {
		return err
	}
	if err := validateTeams(key, teams); err != nil {
		return err
	}
	return nil
}

func validateKey(key string) error {
	for k := range builders {
		if k == key {
			return nil
		}
	}
	return fmt.Errorf("key %s is invalid", key)
}

func validateTeams(key string, teams []string) error {
	builder := builders[key]
	if len(teams) < builder.GetInformation().Min {
		return fmt.Errorf("too few teams for key %s", key)
	} else if len(teams) > builder.GetInformation().Max {
		return fmt.Errorf("too many teams for key %s", key)
	}
	return nil
}
