package indigo

import (
	"fmt"
	"slices"
	"strings"

	cl "github.com/quibbble/quibbble-controller/pkg/collection"
	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type state struct {
	turn                  string
	teams                 []string
	winners               []string
	board                 *board
	deck                  *cl.Collection[tile]
	hands                 map[string]*cl.Collection[tile]
	variant               string
	points                map[string]int
	gemsCount             map[string]int
	round, roundsUntilEnd int
}

func newState(variant string, random int64, teams []string) (*state, error) {
	hands := make(map[string]*cl.Collection[tile])
	points := make(map[string]int)
	gemsCount := make(map[string]int)
	deck := cl.NewCollection[tile](random)
	for idx, numCopies := range numCopiesByUniquePathsIndex {
		for i := 0; i < numCopies; i++ {
			deck.Add(tile{Paths: uniquePaths[idx]})
		}
	}
	deck.Shuffle()

	switch variant {
	case ClassicVariant:
		for _, team := range teams {
			hand := cl.NewCollection[tile](0)
			for i := 0; i < 1; i++ {
				tile, err := deck.Draw()
				if err != nil {
					return nil, err
				}
				hand.Add(*tile)
			}
			points[team] = 0
			gemsCount[team] = 0
			hands[team] = hand
		}
	case LargeHandsVariant:
		for _, team := range teams {
			hand := cl.NewCollection[tile](0)
			for i := 0; i < 2; i++ {
				tile, err := deck.Draw()
				if err != nil {
					return nil, err
				}
				hand.Add(*tile)
			}
			points[team] = 0
			gemsCount[team] = 0
			hands[team] = hand
		}
	}

	return &state{
		turn:           teams[0],
		teams:          teams,
		winners:        make([]string, 0),
		board:          newBoard(teams),
		deck:           deck,
		hands:          hands,
		variant:        variant,
		points:         points,
		gemsCount:      gemsCount,
		round:          0,
		roundsUntilEnd: 99999,
	}, nil
}

func (s *state) rotate(team, paths string) error {
	if !slices.Contains(s.teams, team) {
		return fmt.Errorf("%s not a valid team", team)
	}
	t, err := newTile(paths)
	if err != nil {
		return err
	}
	idx := s.hands[team].IndexOf(*t, func(a, b tile) bool { return a.equals(&b) })
	if idx < 0 {
		return fmt.Errorf("%s's hand does not contain %s", team, paths)
	}
	tile, _ := s.hands[team].GetItem(idx)
	tile.RotateClockwise()
	return nil
}

func (s *state) place(team, paths string, row, col int) error {
	if team != s.turn {
		return fmt.Errorf("%s cannot play on %s turn", team, s.turn)
	}
	t, err := newTile(paths)
	if err != nil {
		return err
	}

	// place tile and remove it from your hand
	tileIdx := s.hands[team].IndexOf(*t, func(a, b tile) bool { return a.equals(&b) })
	if tileIdx < 0 {
		return fmt.Errorf("%s's hand does not contain %s", team, paths)
	}
	if err := s.board.place(t, row, col); err != nil {
		return err
	}
	_ = s.hands[team].Remove(tileIdx)

	// update gem locations
	movedGems, err := s.board.moveGems(row, col)
	if err != nil {
		return err
	}

	// update scores based on new gem locations
	for _, gem := range movedGems {
		if gem.gateway != nil {
			for _, team := range gem.gateway.Teams {
				s.points[team] += colorToPoints[gem.Color]
				s.gemsCount[team] += 1
			}
		}
	}

	// draw tile and add to hand if there tiles left in the deck
	if t, err = s.deck.Draw(); err == nil {
		s.hands[team].Add(*t)
	}

	// change turn
	for idx, team := range s.teams {
		if team == s.turn {
			s.turn = s.teams[(idx+1)%len(s.teams)]
			break
		}
	}

	// inc round counter
	if s.turn == s.teams[0] {
		s.round++
	}

	// check if the game is over and set winners if so
	if s.round >= s.roundsUntilEnd || s.board.gemsInPlay() <= 0 {
		winners := make([]string, 0)
		maxPoints := 0
		for team, points := range s.points {
			if points == maxPoints {
				winners = append(winners, team)
			} else if points > maxPoints {
				winners = []string{team}
				maxPoints = points
			}
		}
		// if tied the player with most points AND gems wins
		if len(winners) > 1 {
			possibleWinners := winners
			winners = make([]string, 0)
			maxGemCount := 0
			for _, team := range possibleWinners {
				gemCount := s.gemsCount[team]
				if gemCount == maxGemCount {
					winners = append(winners, team)
				} else if gemCount > maxGemCount {
					winners = []string{team}
					maxGemCount = gemCount
				}
			}
		}
		s.winners = winners
	}

	return nil
}

func (s *state) actions(team ...string) []*qg.Action {
	targets := make([]*qg.Action, 0)
	if len(s.winners) > 0 {
		return targets
	}
	// place tile actions
	if len(team) == 0 || (len(team) == 1 && team[0] == s.turn) {
		for r, row := range s.board.Tiles {
			for c, t := range row {
				if t == nil {
					for _, t1 := range s.hands[s.turn].GetItems() {
						t2, _ := newTile(t1.Paths)
						t2.RotateClockwise()
						t3, _ := newTile(t2.Paths)
						t3.RotateClockwise()
						t4, _ := newTile(t3.Paths)
						t4.RotateClockwise()
						t5, _ := newTile(t4.Paths)
						t5.RotateClockwise()
						t6, _ := newTile(t5.Paths)
						t6.RotateClockwise()

						targets = append(targets, &qg.Action{
							Team:    s.turn,
							Type:    PlaceAction,
							Details: PlaceDetails{t1.Paths, r, c},
						}, &qg.Action{
							Team:    s.turn,
							Type:    PlaceAction,
							Details: PlaceDetails{t2.Paths, r, c},
						}, &qg.Action{
							Team:    s.turn,
							Type:    PlaceAction,
							Details: PlaceDetails{t3.Paths, r, c},
						}, &qg.Action{
							Team:    s.turn,
							Type:    PlaceAction,
							Details: PlaceDetails{t4.Paths, r, c},
						}, &qg.Action{
							Team:    s.turn,
							Type:    PlaceAction,
							Details: PlaceDetails{t5.Paths, r, c},
						}, &qg.Action{
							Team:    s.turn,
							Type:    PlaceAction,
							Details: PlaceDetails{t6.Paths, r, c},
						})
					}
				}
			}
		}
	}
	return targets
}

func (s *state) message() string {
	message := fmt.Sprintf("%s must place a tile", s.turn)
	if len(s.winners) > 0 {
		message = fmt.Sprintf("%s tie", strings.Join(s.winners, ", "))
		if len(s.winners) == 1 {
			message = fmt.Sprintf("%s wins", s.winners[0])
		}
	}
	return message
}
