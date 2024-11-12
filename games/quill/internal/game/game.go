package game

import (
	"context"
	"math/rand"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	tr "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card/trait"
	hk "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	cn "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/condition"
	ev "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/event"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

var (
	ErrWrongTurn = func(player uuid.UUID) error { return errors.Errorf("'%s' cannot play on other player's turn", player) }
)

type Game struct {
	*en.Engine
	*st.State
	*uuid.Gen
}

func NewGame(seed int64, player1, player2 uuid.UUID, deck1, deck2 []string) (*Game, error) {
	if player1 == player2 {
		return nil, errors.Errorf("player uuids must not match")
	}
	gen := uuid.NewGen(rand.New(rand.NewSource(seed)))
	engineBuilders := en.Builders{
		BuildCondition: cn.NewCondition,
		BuildEvent:     ev.NewEvent,
		BuildHook:      hk.NewHook,
		BuildChoose:    ch.NewChoose,
	}
	cardBuilders := cd.Builders{
		Builders:   engineBuilders,
		BuildTrait: tr.NewTrait,
		Gen:        gen,
	}
	buildCard := func(id string, player uuid.UUID, token bool) (st.ICard, error) {
		return cd.NewCard(&cardBuilders, id, player, token)
	}

	engine := en.NewEngine()
	state, err := st.NewState(seed, buildCard, &engineBuilders, player1, player2, deck1, deck2)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &Game{
		Engine: engine,
		State:  state,
		Gen:    gen,
	}, nil
}

func (g *Game) PlayCard(player, card uuid.UUID, targets ...uuid.UUID) error {
	if player != g.State.GetTurn() {
		return ErrWrongTurn(player)
	}
	event, err := ev.NewEvent(g.Gen.New(en.EventUUID), ev.PlayCardEvent, ev.PlayCardArgs{
		ChoosePlayer: parse.Choose{
			Type: ch.CurrentPlayerChoice,
			Args: ch.CurrentPlayerArgs{},
		},
		ChooseCard: parse.Choose{
			Type: ch.UUIDChoice,
			Args: ch.UUIDArgs{
				UUID: card,
			},
		},
	})
	if err != nil {
		return errors.Wrap(err)
	}
	ctx := context.WithValue(context.WithValue(context.Background(), en.CardCtx, card), en.TargetsCtx, targets)
	if err := g.Engine.Do(ctx, event, g.State); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (g *Game) MoveUnit(player, unit, tile uuid.UUID) error {
	if player != g.State.GetTurn() {
		return ErrWrongTurn(player)
	}
	event, err := ev.NewEvent(g.Gen.New(en.EventUUID), ev.MoveUnitEvent, ev.MoveUnitArgs{
		ChooseUnit: parse.Choose{
			Type: ch.UUIDChoice,
			Args: ch.UUIDArgs{
				UUID: unit,
			},
		},
		ChooseTile: parse.Choose{
			Type: ch.UUIDChoice,
			Args: ch.UUIDArgs{
				UUID: tile,
			},
		},
		UnitMovement: true,
	})
	if err != nil {
		return errors.Wrap(err)
	}
	if err := g.Engine.Do(context.Background(), event, g.State); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (g *Game) MoveUnitXY(player, unit uuid.UUID, x, y int) error {
	if x < 0 || x >= st.Cols || y < 0 || y > st.Rows {
		return errors.ErrIndexOutOfBounds
	}
	return g.MoveUnit(player, unit, g.Board.XYs[x][y].UUID)
}

func (g *Game) AttackUnit(player, attacker, defender uuid.UUID) error {
	if player != g.State.GetTurn() {
		return ErrWrongTurn(player)
	}
	event, err := ev.NewEvent(g.Gen.New(en.EventUUID), ev.AttackUnitEvent, ev.AttackUnitArgs{
		ChooseUnit: parse.Choose{
			Type: ch.UUIDChoice,
			Args: ch.UUIDArgs{
				UUID: attacker,
			},
		},
		ChooseDefender: parse.Choose{
			Type: ch.UUIDChoice,
			Args: ch.UUIDArgs{
				UUID: defender,
			},
		},
	})
	if err != nil {
		return errors.Wrap(err)
	}
	if err := g.Engine.Do(context.Background(), event, g.State); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (g *Game) SackCard(player, card uuid.UUID, option string) error {
	if player != g.State.GetTurn() {
		return ErrWrongTurn(player)
	}
	event, err := ev.NewEvent(g.Gen.New(en.EventUUID), ev.SackCardEvent, ev.SackCardArgs{
		ChoosePlayer: parse.Choose{
			Type: ch.CurrentPlayerChoice,
			Args: ch.CurrentPlayerArgs{},
		},
		ChooseCard: parse.Choose{
			Type: ch.UUIDChoice,
			Args: ch.UUIDArgs{
				UUID: card,
			},
		},
		SackOption: option,
	})
	if err != nil {
		return errors.Wrap(err)
	}
	if err := g.Engine.Do(context.Background(), event, g.State); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (g *Game) EndTurn(player uuid.UUID) error {
	if player != g.State.GetTurn() {
		return ErrWrongTurn(player)
	}
	event, err := ev.NewEvent(g.Gen.New(en.EventUUID), ev.EndTurnEvent, ev.EndTurnArgs{})
	if err != nil {
		return errors.Wrap(err)
	}
	if err := g.Engine.Do(context.Background(), event, g.State); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// GetNextTargets given a list of past targets, return a list of next valid targets in the chain.
// When returned list is empty you know that input list is a complete chain that can be used to
// perform an action on the game state.
func (g *Game) GetNextTargets(player uuid.UUID, targets ...uuid.UUID) ([]uuid.UUID, error) {
	if player != g.State.GetTurn() {
		return nil, ErrWrongTurn(player)
	}
	switch len(targets) {
	case 0:
		choices := make([]uuid.UUID, 0)
		choose1, err := ch.NewChoose(g.Gen.New(en.ChooseUUID), ch.UnitsChoice, &ch.UnitsArgs{
			Types: []string{cd.CreatureUnit},
		})
		if err != nil {
			return nil, errors.Wrap(err)
		}
		choose2, err := ch.NewChoose(g.Gen.New(en.ChooseUUID), ch.OwnedUnitsChoice, &ch.OwnedTilesArgs{
			ChoosePlayer: parse.Choose{
				Type: ch.CurrentPlayerChoice,
				Args: ch.CurrentPlayerArgs{},
			},
		})
		if err != nil {
			return nil, errors.Wrap(err)
		}
		c, err := ch.NewChooseChain(ch.SetIntersect, choose1, choose2).Retrieve(context.Background(), g.Engine, g.State)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		for _, choice := range c {
			x, y, err := g.Board.GetUnitXY(choice)
			if err != nil {
				return nil, errors.Wrap(err)
			}
			unit := g.Board.XYs[x][y].Unit.(*cd.UnitCard)

			canMove := false
			if unit.Movement > 0 {
				moveChoose1, err := ch.NewChoose(g.State.Gen.New(en.ChooseUUID), ch.CodexChoice, &ch.CodexArgs{
					Types: []string{"Tile"},
					Codex: unit.Codex,
					ChooseUnitOrTile: parse.Choose{
						Type: ch.UUIDChoice,
						Args: ch.UUIDArgs{
							UUID: unit.UUID,
						},
					},
				})
				if err != nil {
					return nil, errors.Wrap(err)
				}
				moveChoose2, err := ch.NewChoose(g.State.Gen.New(en.ChooseUUID), ch.TilesChoice, &ch.TilesArgs{
					Empty: true,
				})
				if err != nil {
					return nil, errors.Wrap(err)
				}
				ctx := context.WithValue(context.Background(), en.CardCtx, unit.GetUUID())
				moveChoices, err := ch.NewChooseChain(ch.SetIntersect, moveChoose1, moveChoose2).Retrieve(ctx, g.Engine, g.State)
				if err != nil {
					return nil, errors.Wrap(err)
				}
				if len(moveChoices) > 0 {
					canMove = true
				}
			}
			canAttack := false
			if unit.Cooldown == 0 {
				var choose1 en.IChoose
				if len(unit.GetTraits(tr.RangedTrait)) > 0 {
					choose1, err = ch.NewChoose(g.State.Gen.New(en.ChooseUUID), ch.RangedChoice, &ch.RangedArgs{
						Types: []string{"Unit"},
						Range: unit.GetTraits(tr.RangedTrait)[0].GetArgs().(*tr.RangedArgs).Amount,
						ChooseUnitOrTile: parse.Choose{
							Type: ch.UUIDChoice,
							Args: ch.UUIDArgs{
								UUID: unit.UUID,
							},
						},
					})
					if err != nil {
						return nil, errors.Wrap(err)
					}
				} else {
					choose1, err = ch.NewChoose(g.State.Gen.New(en.ChooseUUID), ch.CodexChoice, &ch.CodexArgs{
						Types: []string{"Unit"},
						Codex: unit.Codex,
						ChooseUnitOrTile: parse.Choose{
							Type: ch.UUIDChoice,
							Args: ch.UUIDArgs{
								UUID: unit.UUID,
							},
						},
					})
					if err != nil {
						return nil, errors.Wrap(err)
					}
				}
				choose2, err := ch.NewChoose(g.Gen.New(en.ChooseUUID), ch.OwnedUnitsChoice, &ch.OwnedUnitsArgs{
					ChoosePlayer: parse.Choose{
						Type: ch.OpposingPlayerChoice,
						Args: ch.OpposingPlayerArgs{},
					},
				})
				if err != nil {
					return nil, errors.Wrap(err)
				}
				ctx := context.WithValue(context.Background(), en.CardCtx, unit.GetUUID())
				choices, err := ch.NewChooseChain(ch.SetIntersect, choose1, choose2).Retrieve(ctx, g.Engine, g.State)
				if err != nil {
					return nil, errors.Wrap(err)
				}
				if len(choices) > 0 {
					canAttack = true
				}
			}
			if canMove || canAttack {
				choices = append(choices, unit.GetUUID())
			}
		}

		for _, card := range g.Hand[player].GetItems() {
			if card.GetCost() > g.Mana[player].Amount {
				continue
			}
			playable, err := card.Playable(g.Engine, g.State)
			if err != nil {
				return nil, errors.Wrap(err)
			}
			if !playable {
				continue
			}
			targets, err := card.NextTargets(context.WithValue(context.Background(), en.TargetsCtx, []uuid.UUID{}), g.Engine, g.State)
			if err != nil {
				return nil, errors.Wrap(err)
			}
			if len(targets) == 0 && len(card.GetTargets()) != 0 {
				continue
			}
			choices = append(choices, card.GetUUID())
		}
		if !g.Sacked[player] {
			choices = append(choices, uuid.UUID(en.SackUUID))
		}
		return choices, nil
	default:
		if c, err := g.Hand[player].GetCard(targets[0]); err == nil {
			if playable, err := c.Playable(g.Engine, g.State); err != nil || !playable || c.GetCost() > g.Mana[player].Amount {
				return nil, errors.Errorf("'%s' not playable", targets[0])
			}
			ctx := context.WithValue(context.Background(), en.TargetsCtx, targets[1:])
			return c.NextTargets(ctx, g.Engine, g.State)
		} else if x1, y1, err := g.Board.GetUnitXY(targets[0]); err == nil {
			unit := g.Board.XYs[x1][y1].Unit.(*cd.UnitCard)
			if unit.Type != cd.CreatureUnit {
				return nil, errors.Errorf("'%s' unit may not move or attack", unit.Type)
			}
			switch len(targets) {
			case 1:
				var moveChoices, attackChoices []uuid.UUID
				if unit.Movement > 0 {
					moveChoose1, err := ch.NewChoose(g.State.Gen.New(en.ChooseUUID), ch.CodexChoice, &ch.CodexArgs{
						Types: []string{"Tile"},
						Codex: unit.Codex,
						ChooseUnitOrTile: parse.Choose{
							Type: ch.UUIDChoice,
							Args: ch.UUIDArgs{
								UUID: unit.UUID,
							},
						},
					})
					if err != nil {
						return nil, errors.Wrap(err)
					}
					moveChoose2, err := ch.NewChoose(g.State.Gen.New(en.ChooseUUID), ch.TilesChoice, &ch.TilesArgs{
						Empty: true,
					})
					if err != nil {
						return nil, errors.Wrap(err)
					}
					ctx := context.WithValue(context.Background(), en.CardCtx, unit.GetUUID())
					moveChoices, err = ch.NewChooseChain(ch.SetIntersect, moveChoose1, moveChoose2).Retrieve(ctx, g.Engine, g.State)
					if err != nil {
						return nil, errors.Wrap(err)
					}
				}

				if unit.Cooldown == 0 {
					var attackChoose1 en.IChoose
					if len(unit.GetTraits(tr.RangedTrait)) > 0 {
						attackChoose1, err = ch.NewChoose(g.State.Gen.New(en.ChooseUUID), ch.RangedChoice, &ch.RangedArgs{
							Types: []string{"Unit"},
							Range: unit.GetTraits(tr.RangedTrait)[0].GetArgs().(*tr.RangedArgs).Amount,
							ChooseUnitOrTile: parse.Choose{
								Type: ch.UUIDChoice,
								Args: ch.UUIDArgs{
									UUID: unit.UUID,
								},
							},
						})
						if err != nil {
							return nil, errors.Wrap(err)
						}
					} else {
						attackChoose1, err = ch.NewChoose(g.State.Gen.New(en.ChooseUUID), ch.CodexChoice, &ch.CodexArgs{
							Types: []string{"Unit"},
							Codex: unit.Codex,
							ChooseUnitOrTile: parse.Choose{
								Type: ch.UUIDChoice,
								Args: ch.UUIDArgs{
									UUID: unit.UUID,
								},
							},
						})
						if err != nil {
							return nil, errors.Wrap(err)
						}
					}
					attackChoose2, err := ch.NewChoose(g.State.Gen.New(en.ChooseUUID), ch.OwnedUnitsChoice, &ch.OwnedUnitsArgs{
						ChoosePlayer: parse.Choose{
							Type: ch.OpposingPlayerChoice,
							Args: ch.OpposingPlayerArgs{},
						},
					})
					if err != nil {
						return nil, errors.Wrap(err)
					}
					ctx := context.WithValue(context.Background(), en.CardCtx, unit.GetUUID())
					attackChoices, err = ch.NewChooseChain(ch.SetIntersect, attackChoose1, attackChoose2).Retrieve(ctx, g.Engine, g.State)
					if err != nil {
						return nil, errors.Wrap(err)
					}
				}

				choices := append(moveChoices, attackChoices...)
				if len(choices) == 0 {
					return nil, errors.Errorf("'%s' cannot move or attack", unit.GetUUID())
				}
				return choices, nil
			case 2:
				switch targets[1].Type() {
				case en.UnitUUID:
					x2, y2, err := g.Board.GetUnitXY(targets[1])
					if err != nil {
						return nil, errors.Wrap(err)
					}
					if unit.Cooldown != 0 {
						return nil, errors.Errorf("'%s' may not attack due to cooldown stat", unit.GetUUID())
					}
					ranged := unit.GetTraits(tr.RangedTrait)
					if len(ranged) == 0 && !unit.CheckCodex(x1, y1, x2, y2) {
						return nil, errors.Errorf("invalid attack for unit '%s'", unit.GetUUID())
					}
					if len(ranged) > 0 && !ranged[0].GetArgs().(*tr.RangedArgs).CheckRange(x1, y1, x2, y2) {
						return nil, errors.Errorf("invalid ranged attack for unit '%s'", unit.GetUUID())
					}
					return make([]uuid.UUID, 0), nil
				case en.TileUUID:
					x2, y2, err := g.Board.GetTileXY(targets[1])
					if err != nil {
						return nil, errors.Wrap(err)
					}
					if unit.Movement <= 0 {
						return nil, errors.Errorf("'%s' may not move due to movement stat", unit.GetUUID())
					}
					if !unit.CheckCodex(x1, y1, x2, y2) {
						return nil, errors.Errorf("invalid move for unit '%s'", unit.GetUUID())
					}
					return make([]uuid.UUID, 0), nil
				default:
					return nil, en.ErrInvalidUUIDType(targets[1], en.UnitUUID, en.TileUUID)
				}
			default:
				return nil, errors.ErrInvalidSliceLength
			}
		}
		return nil, errors.Errorf("invalid target list")
	}
}
