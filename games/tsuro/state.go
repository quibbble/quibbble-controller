package tsuro

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

type state struct {
	seed int64
	rand *rand.Rand

	variant string

	turn    string
	teams   []string
	winners []string

	board           *board
	deck            *deck
	tokens          map[string]*token
	hands           map[string]*hand
	dragon          string
	playedFirstTurn map[string]bool // teams that have placed and still alive
	alive           map[string]bool // teams that are alive
	points          map[string]int
}

func newState(variant string, seed int64, teams []string) (*state, error) {
	r := rand.New(rand.NewSource(seed))
	hands := make(map[string]*hand)
	tokens := make(map[string]*token)
	alive := make(map[string]bool)
	deck := newDeck(r)
	points := make(map[string]int)
	switch variant {
	case ClassicVariant, SoloVariant:
		for _, team := range teams {
			hand := newHand()
			for i := 0; i < 3; i++ {
				tile, err := deck.Draw()
				if err != nil {
					return nil, err
				}
				hand.Add(tile)
			}
			hands[team] = hand
			token := uniqueRandomToken(tokens, r)
			tokens[team] = token
			alive[team] = true
		}
	case LongestPathVariant, MostCrossingsVariant:
		for _, team := range teams {
			hand := newHand()
			for i := 0; i < 3; i++ {
				tile, err := deck.Draw()
				if err != nil {
					return nil, err
				}
				hand.Add(tile)
			}
			hands[team] = hand
			token := uniqueRandomToken(tokens, r)
			tokens[team] = token
			alive[team] = true
			points[team] = 0
		}
	case OpenTilesVariant:
		hand := newHand()
		for i := 0; i < 3; i++ {
			tile, err := deck.Draw()
			if err != nil {
				return nil, err
			}
			hand.Add(tile)
		}
		for _, team := range teams {
			hands[team] = hand
			token := uniqueRandomToken(tokens, r)
			tokens[team] = token
			alive[team] = true
		}
	default:
		return nil, fmt.Errorf("invalid variant %s", variant)
	}
	if len(teams) != len(tokens) {
		return nil, fmt.Errorf("failed to build new state likely due to duplicate teams")
	}
	return &state{
		seed:            seed,
		rand:            r,
		variant:         variant,
		turn:            teams[0],
		teams:           teams,
		winners:         make([]string, 0),
		board:           newBoard(),
		deck:            deck,
		tokens:          tokens,
		hands:           hands,
		dragon:          "",
		playedFirstTurn: make(map[string]bool),
		alive:           alive,
		points:          points,
	}, nil
}

// Rotate rotates a tile in hand clockwise
func (s *state) Rotate(team, tile string) error {
	if s.variant == OpenTilesVariant && team != s.turn {
		return fmt.Errorf("%s cannot rotate tile on %s turn", team, s.turn)
	}
	if !slices.Contains(s.teams, team) {
		return fmt.Errorf("%s not a valid team", team)
	}
	t, err := newTile(tile)
	if err != nil {
		return err
	}
	if !t.in(s.hands[team].hand) {
		return fmt.Errorf("%s's hand does not contain %s", team, tile)
	}
	s.hands[team].hand[s.hands[team].IndexOf(t)].RotateRight()
	return nil
}

// place places a tile on the board
func (s *state) Place(team, tile string, row, column int) error {
	if team != s.turn {
		return fmt.Errorf("%s cannot play on %s turn", team, s.turn)
	}
	if !s.playedFirstTurn[s.turn] && (s.tokens[team].Row != row || s.tokens[team].Col != column) {
		return fmt.Errorf("%s cannot place in row %d column %d", team, row, column)
	} else if s.playedFirstTurn[s.turn] {
		adj, err := s.tokens[team].getAdjacent()
		if err != nil {
			return err
		}
		if row != adj.Row || column != adj.Col {
			return fmt.Errorf("%s cannot place in row %d column %d", team, row, column)
		}
	}
	t, err := newTile(tile)
	if err != nil {
		return err
	}
	if !t.in(s.hands[team].hand) {
		return fmt.Errorf("%s's hand does not contain %s", team, tile)
	}
	if err := s.hands[team].Remove(t); err != nil {
		return err
	}
	if err := s.board.Place(t, row, column); err != nil {
		return err
	}
	if !s.playedFirstTurn[s.turn] {
		s.playedFirstTurn[s.turn] = true
	}
	s.moveTokens()
	if s.variant == LongestPathVariant || s.variant == MostCrossingsVariant {
		s.score()
	}
	s.updateAlive()
	s.handleDraws()
	s.nextTurn()
	return nil
}

func (s *state) SetWinners(winners []string) error {
	for _, winner := range winners {
		if !slices.Contains(s.teams, winner) {
			return fmt.Errorf("winner not in teams")
		}
	}
	s.winners = winners
	return nil
}

func (s *state) moveTokens() {
	moved := 0
	move := map[string]string{"A": "F", "B": "E", "C": "H", "D": "G", "E": "B", "F": "A", "G": "D", "H": "C"}
	for team, token := range s.tokens {
		if s.playedFirstTurn[team] {
			t := s.board.board[token.Row][token.Col]
			if !mapContainsVal(t.Paths, team) {
				// first placement so move through the just placed tile
				destination := t.GetDestination(token.Notch)
				t.Paths[token.Notch+destination] = team
				token.Notch = destination
				// token was moved
				moved++
			} else if s.collided(s.tokens, team, token) {
				// token collided with other token
				continue
			} else {
				// normal case
				var nextTile *tile
				if strings.Contains("AB", token.Notch) && token.Row-1 >= 0 && s.board.board[token.Row-1][token.Col] != nil {
					nextTile = s.board.board[token.Row-1][token.Col]
					token.Row -= 1
				} else if strings.Contains("CD", token.Notch) && token.Col+1 < columns && s.board.board[token.Row][token.Col+1] != nil {
					nextTile = s.board.board[token.Row][token.Col+1]
					token.Col += 1
				} else if strings.Contains("EF", token.Notch) && token.Row+1 < rows && s.board.board[token.Row+1][token.Col] != nil {
					nextTile = s.board.board[token.Row+1][token.Col]
					token.Row += 1
				} else if strings.Contains("GH", token.Notch) && token.Col-1 >= 0 && s.board.board[token.Row][token.Col-1] != nil {
					nextTile = s.board.board[token.Row][token.Col-1]
					token.Col -= 1
				} else {
					continue
				}
				// move the token to the notch on the next tile
				startNotch := move[token.Notch]
				// where the token ends up on the next tile
				endNotch := nextTile.GetDestination(startNotch)
				// update token location
				nextTile.Paths[startNotch+endNotch] = team
				token.Notch = endNotch
				// token was moved
				moved++
			}
		}
	}
	if moved > 0 {
		s.moveTokens()
	}
}

func (s *state) collided(tokens map[string]*token, team string, token *token) bool {
	for team2, token2 := range tokens {
		if team != team2 && (token.collided(token2) || token.equals(token2)) {
			return true
		}
	}
	return false
}

func (s *state) score() {
	switch s.variant {
	case LongestPathVariant:
		points := make(map[string]int)
		for _, team := range s.teams {
			points[team] = 0
		}
		for _, row := range s.board.board {
			for _, tile := range row {
				if tile != nil {
					for _, team := range tile.Paths {
						points[team]++
					}
				}
			}
		}
		s.points = points
	case MostCrossingsVariant:
		points := make(map[string]int)
		for _, team := range s.teams {
			points[team] = 0
		}
		for _, row := range s.board.board {
			for _, tile := range row {
				if tile != nil {
					for _, team := range s.teams {
						points[team] += tile.countCrossings(team)
					}
				}
			}
		}
		s.points = points
	}
}

func (s *state) updateAlive() {
	if len(s.winners) > 0 {
		return
	}
	// alive before checking
	initialAlive := make([]string, 0)
	for _, team := range s.teams {
		if s.alive[team] {
			initialAlive = append(initialAlive, team)
		}
	}
	// update who is still alive
	for team, token := range s.tokens {
		if s.playedFirstTurn[team] {
			if (token.Row == 0 && strings.Contains("AB", token.Notch)) ||
				(token.Row == rows-1 && strings.Contains("EF", token.Notch)) ||
				(token.Col == 0 && strings.Contains("GH", token.Notch)) ||
				(token.Col == columns-1 && strings.Contains("CD", token.Notch)) {
				// check on board edge
				s.setLost(team)
			} else if s.collided(s.tokens, team, token) {
				// check if collided with another token
				s.setLost(team)
			}
		}
	}
	// who is still alive
	stillAlive := make([]string, 0)
	for _, team := range s.teams {
		if s.alive[team] {
			stillAlive = append(stillAlive, team)
		}
	}
	switch s.variant {
	case ClassicVariant, OpenTilesVariant:
		if len(stillAlive) == 0 { // no more alive so initial alive all win
			s.winners = initialAlive
		} else if len(stillAlive) == 1 { // one alive so they win
			s.winners = stillAlive
		} else if s.board.getTileCount() == len(tiles) { // all tiles have been placed remaining alive are winners
			s.winners = stillAlive
		}
	case LongestPathVariant, MostCrossingsVariant:
		max := max(s.points)
		if len(stillAlive) == 0 { // no more alive
			s.winners = max
		} else if s.board.getTileCount() == len(tiles) { // all tiles have been placed
			s.winners = max
		} else if len(stillAlive) == 1 && len(max) == 1 && max[0] == stillAlive[0] { // last remaining has the most points to wins
			s.winners = max
		}
	case SoloVariant:
		if len(stillAlive) == 0 {
			s.winners = []string{"FAIL"}
		} else if s.board.getTileCount() == len(tiles) { // win if all tokens are still on board and all tiles have been placed
			s.winners = stillAlive
		}
	}
}

func (s *state) handleDraws() {
	if len(s.winners) > 0 {
		return
	}
	current := s.turn
	if s.dragon != "" {
		current = s.dragon
	}
	for s.alive[current] && len(s.deck.deck) > 0 && len(s.hands[current].hand) < 3 {
		tile, err := s.deck.Draw()
		if err != nil {
			return
		}
		s.hands[current].Add(tile)
		current = s.getNextTurn(current)
	}
	if len(s.deck.deck) == 0 && len(s.hands[current].hand) < 3 {
		s.dragon = current
	} else {
		s.dragon = ""
	}
}

func (s *state) nextTurn() {
	if len(s.winners) > 0 {
		return
	}
	s.turn = s.getNextTurn(s.turn)
}

func (s *state) getNextTurn(turn string) string {
	nextTurn := ""
	if len(s.winners) > 0 {
		return nextTurn
	}
	for idx, team := range s.teams {
		if team == turn {
			nextTurn = s.teams[(idx+1)%len(s.teams)]
			if !s.alive[nextTurn] {
				return s.getNextTurn(nextTurn)
			}
			return nextTurn
		}
	}
	return nextTurn
}

func (s *state) setLost(team string) {
	s.alive[team] = false
	s.playedFirstTurn[team] = false
	s.deck.Add(s.hands[team].hand...)
	s.hands[team].Clear()
	if s.aliveCount() <= 0 {
		return
	}
	next := s.getNextTurn(s.turn)
	if s.dragon == team && len(s.hands[next].hand) < 3 {
		s.dragon = next
	}
}

func (s *state) aliveCount() int {
	count := 0
	for _, alive := range s.alive {
		if alive {
			count++
		}
	}
	return count
}

func (s *state) actions(team ...string) []*qg.Action {
	// rotate actions are not stored in qgn so ignore them here as well
	targets := make([]*qg.Action, 0)
	if len(team) == 0 || (len(team) == 1 && team[0] == s.turn) {
		t, _ := s.tokens[s.turn].getAdjacent()
		row := t.Row
		col := t.Col
		for _, t1 := range s.hands[s.turn].hand {
			t2, _ := newTile(t1.Edges)
			t2.RotateRight()
			t3, _ := newTile(t2.Edges)
			t3.RotateRight()
			t4, _ := newTile(t3.Edges)
			t4.RotateRight()
			targets = append(targets, &qg.Action{
				Team:    s.turn,
				Type:    PlaceAction,
				Details: PlaceDetails{Row: row, Col: col, Tile: t1.Edges},
			}, &qg.Action{
				Team:    s.turn,
				Type:    PlaceAction,
				Details: PlaceDetails{Row: row, Col: col, Tile: t2.Edges},
			}, &qg.Action{
				Team:    s.turn,
				Type:    PlaceAction,
				Details: PlaceDetails{Row: row, Col: col, Tile: t3.Edges},
			}, &qg.Action{
				Team:    s.turn,
				Type:    PlaceAction,
				Details: PlaceDetails{Row: row, Col: col, Tile: t4.Edges},
			})
		}
	}
	return targets
}

func (s *state) message() string {
	message := fmt.Sprintf("%s must place a tile", s.turn)
	if len(s.winners) > 0 {
		switch s.variant {
		case ClassicVariant, OpenTilesVariant, LongestPathVariant, MostCrossingsVariant:
			message = fmt.Sprintf("%s tie", strings.Join(s.winners, ", "))
			if len(s.winners) == 1 {
				message = fmt.Sprintf("%s wins", s.winners[0])
			}
		case SoloVariant:
			if len(s.winners) == 1 && s.winners[0] == "FAIL" {
				message = "you saved 0 tokens"
			} else if len(s.winners) < Max {
				message = fmt.Sprintf("you saved %d tokens", len(s.winners))
			} else {
				message = "you saved all the tokens"
			}
		}
	}
	return message
}
