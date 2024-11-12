package tests

import (
	gm "github.com/quibbble/quibbble-controller/games/quill/internal/game"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const (
	Seed    = int64(0)
	Player1 = uuid.UUID("P0000001")
	Player2 = uuid.UUID("P0000002")
)

func NewTestEnv(player uuid.UUID, ids ...string) (*gm.Game, []uuid.UUID, error) {

	deck1 := map[string]int{
		"U0002": 30,
	}
	deck2 := map[string]int{
		"U0002": 30,
	}

	d1 := make([]string, 0)
	d2 := make([]string, 0)

	for id, count := range deck1 {
		for i := 0; i < count; i++ {
			d1 = append(d1, id)
		}
	}
	for id, count := range deck2 {
		for i := 0; i < count; i++ {
			d2 = append(d2, id)
		}
	}

	game, err := gm.NewGame(Seed, Player1, Player2, d1, d2)
	if err != nil {
		return nil, nil, errors.Wrap(err)
	}

	uuids := make([]uuid.UUID, 0)

	for _, id := range ids {
		card, err := game.BuildCard(id, player, false)
		if err != nil {
			return nil, nil, errors.Wrap(err)
		}
		uuids = append(uuids, card.GetUUID())
		game.Hand[player].Add(card)
	}

	game.Mana[Player1].Amount = 8
	game.Mana[Player1].BaseAmount = 8
	game.Mana[Player2].Amount = 8
	game.Mana[Player2].BaseAmount = 8

	return game, uuids, nil
}
