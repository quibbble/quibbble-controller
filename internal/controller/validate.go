package controller

import (
	"fmt"

	"github.com/quibbble/quibbble-controller/games"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
	qgn "github.com/quibbble/quibbble-controller/pkg/gamenotation"
)

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
	for _, builder := range games.Builders {
		if builder.GetInformation().Key == key {
			return nil
		}
	}
	return fmt.Errorf("key %s is invalid", key)
}

func validateTeams(key string, teams []string) error {
	var builder qg.GameBuilder
	for _, b := range games.Builders {
		if b.GetInformation().Key == key {
			builder = b
			break
		}
	}
	if builder == nil {
		return fmt.Errorf("key %s is invalid", key)
	}
	if len(teams) < builder.GetInformation().Min {
		return fmt.Errorf("too few teams for key %s", key)
	} else if len(teams) > builder.GetInformation().Max {
		return fmt.Errorf("too many teams for key %s", key)
	}
	return nil
}
