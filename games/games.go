package games

import (
	"github.com/quibbble/quibbble-controller/games/carcassonne"
	"github.com/quibbble/quibbble-controller/games/connect4"
	"github.com/quibbble/quibbble-controller/games/indigo"
	"github.com/quibbble/quibbble-controller/games/stratego"
	"github.com/quibbble/quibbble-controller/games/tictactoe"
	"github.com/quibbble/quibbble-controller/games/tsuro"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

var Builders = []qg.GameBuilder{
	carcassonne.Builder{},
	connect4.Builder{},
	indigo.Builder{},
	stratego.Builder{},
	tictactoe.Builder{},
	tsuro.Builder{},
}
