package quill

import (
	"math"

	"github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type AI struct{}

func (ai AI) Score(game qg.Game, team string) float64 {
	g := game.(*Quill)
	if g.state.Winner != nil {
		if string(*g.state.Winner) == team {
			return math.Inf(1)
		} else {
			return math.Inf(-1)
		}
	}
	score := 0.0
	for _, col := range g.state.Board.XYs {
		for _, tile := range col {
			if tile.Unit != nil &&
				tile.Unit.GetID() == "U0001" &&
				string(tile.Unit.GetPlayer()) == team {
				score += float64(tile.Unit.(*card.UnitCard).Health)
			}
		}
	}
	return score
}
