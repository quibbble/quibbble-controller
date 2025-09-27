package controller

import (
	"errors"

	"github.com/quibbble/quibbble-controller/pkg/persist"
	"github.com/quibbble/quibbble-controller/pkg/proto/controller"
	"github.com/upper/db/v4"
)

func (c *Controller) lookupActiveGame(gk *controller.GameKey) (*persist.Game, error) {
	if c.config.Persistence == nil {
		return nil, nil
	}
	activeGames := c.db.Collection(persist.ActiveGamesTable)
	var game *persist.Game
	err := activeGames.Find(db.Cond{
		persist.RepositoryColumn: gk.Repository,
		persist.TagColumn:        gk.Tag,
		persist.NameColumn:       gk.Name,
	}).One(game)
	if err != nil && !errors.Is(err, db.ErrNoMoreRows) {
		return nil, err
	}
	return game, nil
}
