package carcassonne

import (
	"fmt"
	"math/rand"
	"slices"
	"sort"
	"strings"

	qg "github.com/quibbble/quibbble-controller/pkg/game"
)

// state holds all necessary game objects and high level game logic
type state struct {
	seed            int64
	turn            string
	teams           []string
	winners         []string
	playTiles       map[string]*tile // teams to the tiles to place onto the board at the start of any given turn
	lastPlacedTiles map[string]*tile // the tiles that were last placed by each team
	board           *board
	boardTokens     []*token       // a list of tokens currently on the board
	tokens          map[string]int // number of tokens each team can play
	scores          map[string]int // points of each team
	deck            *deck
}

func newState(seed int64, teams []string) *state {
	r := rand.New(rand.NewSource(seed))
	tokens := make(map[string]int)
	scores := make(map[string]int)
	playTiles := make(map[string]*tile)
	lastPlacedTiles := make(map[string]*tile)
	for _, team := range teams {
		tokens[team] = 7
		scores[team] = 0
	}
	deck := newDeck(r)
	for _, team := range teams {
		tile, _ := deck.Draw()
		playTiles[team] = tile
		lastPlacedTiles[team] = nil
	}
	return &state{
		seed:            seed,
		turn:            teams[0],
		teams:           teams,
		winners:         make([]string, 0),
		playTiles:       playTiles,
		lastPlacedTiles: lastPlacedTiles,
		board:           newBoard(),
		boardTokens:     make([]*token, 0),
		tokens:          tokens,
		scores:          scores,
		deck:            deck,
	}
}

func (s *state) RotateTileRight(team string) error {
	if len(s.winners) > 0 {
		return fmt.Errorf("game already over")
	}
	if s.playTiles[team] == nil {
		return fmt.Errorf("cannot rotate tile")
	}
	s.playTiles[team].RotateRight()
	return nil
}

func (s *state) RotateTileLeft(team string) error {
	if len(s.winners) > 0 {
		return fmt.Errorf("game already over")
	}
	if s.playTiles[team] == nil {
		return fmt.Errorf("cannot rotate tile")
	}
	s.playTiles[team].RotateLeft()
	return nil
}

func (s *state) PlaceTile(team string, tile *tile, x, y int) error {
	if len(s.winners) > 0 {
		return fmt.Errorf("game already over")
	}
	if team != s.turn {
		return fmt.Errorf("%s cannot play on %s turn", team, s.turn)
	}
	if s.playTiles[team] == nil || !tile.equals(s.playTiles[team]) {
		return fmt.Errorf("tile %+v doesn't match %+v", tile, s.playTiles[team])
	}
	if err := s.board.Place(tile, x, y); err != nil {
		return err
	}
	s.lastPlacedTiles[team] = tile
	s.playTiles[team] = nil

	// if there are no tokens to place or cannot place token anywhere skip place token action here
	if s.tokens[s.turn] == 0 || len(s.actions()) <= 1 {
		if err := s.PlaceToken(s.turn, true, 0, 0, "", ""); err != nil {
			return err
		}
	}
	return nil
}

func (s *state) PlaceToken(team string, pass bool, x, y int, typ, side string) error {
	if len(s.winners) > 0 {
		return fmt.Errorf("game already over")
	}
	if team != s.turn {
		return fmt.Errorf("currently %s's turn", s.turn)
	}
	if s.playTiles[team] != nil {
		return fmt.Errorf("cannot place token")
	}
	// try placing token
	if !pass {
		if s.lastPlacedTiles[team].X != x || s.lastPlacedTiles[team].Y != y {
			return fmt.Errorf("cannot place token on tile at %d,%d", s.lastPlacedTiles[team].X, s.lastPlacedTiles[team].Y)
		}
		if !slices.Contains(TokenTypes, typ) {
			return fmt.Errorf("invalid token type %s", typ)
		}
		if (typ == Thief || typ == Knight) && !slices.Contains(Sides, side) {
			return fmt.Errorf("invalid side %s with token %s", side, typ)
		} else if typ == Farmer && !slices.Contains(FarmSides, side) {
			return fmt.Errorf("invalid farm side %s with token %s", side, typ)
		} else if typ == Monk && s.lastPlacedTiles[team].Center != Cloister {
			return fmt.Errorf("cannot place %s on tile that does not contain %s", Monk, Cloister)
		}
		if s.tokens[team] <= 0 {
			return fmt.Errorf("not enough tokens to place for team %s", team)
		}
		// check to ensure token does not connect to pre-existing tokens in given structure
		switch typ {
		case Thief:
			road, err := s.board.generateRoad(x, y, side)
			if err != nil {
				return err
			}
			tokens := tokensInStructure(s.boardTokens, road)
			if len(tokens) > 0 {
				return fmt.Errorf("cannot place token on road that is already claimed")
			}
		case Knight:
			city, err := s.board.generateCity(x, y, side)
			if err != nil {
				return err
			}
			tokens := tokensInStructure(s.boardTokens, city)
			if len(tokens) > 0 {
				return fmt.Errorf("cannot place token on city that is already claimed")
			}
		case Farmer:
			farm, err := s.board.generateFarm(x, y, side)
			if err != nil {
				return err
			}
			tokens := tokensInStructure(s.boardTokens, farm)
			if len(tokens) > 0 {
				return fmt.Errorf("cannot place token on farm that is already claimed")
			}
		}
		// add the token
		s.tokens[team]--
		token := newToken(x, y, team, typ, side)
		s.boardTokens = append(s.boardTokens, token)
	}
	// score completed cities
	citySides := make([]string, 0)
	for _, side := range Sides {
		if s.lastPlacedTiles[team].Sides[side] == City {
			citySides = append(citySides, side)
		}
	}
	if len(citySides) > 0 {
		if s.lastPlacedTiles[team].ConnectedCitySides {
			citySides = citySides[:1]
		}
		for _, citySide := range citySides {
			city, err := s.board.generateCity(s.lastPlacedTiles[team].X, s.lastPlacedTiles[team].Y, citySide)
			if err != nil {
				return err
			}
			if city.complete {
				// add to completed list in board
				s.board.completeCities = append(s.board.completeCities, city)

				// check if token inside city
				inside := tokensInStructure(s.boardTokens, city)
				if len(inside) > 0 {
					// score and add points
					points, err := scoreCity(city)
					if err != nil {
						return err
					}
					winners := pointsWinners(inside)
					for _, winner := range winners {
						s.scores[winner] += points
					}
					// remove inside from board and add back to tokens pile
					for _, token := range inside {
						s.tokens[token.Team]++
					}
					s.boardTokens = removeTokens(s.boardTokens, inside...)
					// set color of completed
					for _, n := range city.nodes {
						for _, side := range n.sides {
							n.tile.Teams[side] = winners
						}
					}
				}
			}
		}
	}
	// score completed roads
	roadSides := make([]string, 0)
	for _, side := range Sides {
		if s.lastPlacedTiles[team].Sides[side] == Road {
			roadSides = append(roadSides, side)
		}
	}
	if len(roadSides) > 0 {
		if len(roadSides) <= 2 {
			roadSides = roadSides[:1]
		}
		for _, roadSide := range roadSides {
			road, err := s.board.generateRoad(s.lastPlacedTiles[team].X, s.lastPlacedTiles[team].Y, roadSide)
			if err != nil {
				return err
			}
			if road.complete {
				inside := tokensInStructure(s.boardTokens, road)
				if len(inside) > 0 {
					// score and add points
					points, err := scoreRoad(road)
					if err != nil {
						return err
					}
					winners := pointsWinners(inside)
					for _, winner := range winners {
						s.scores[winner] += points
					}
					// remove inside from board and add back to tokens pile
					for _, token := range inside {
						s.tokens[token.Team]++
					}
					s.boardTokens = removeTokens(s.boardTokens, inside...)
					// add to completed list in board
					s.board.completeRoads = append(s.board.completeRoads, road)
					// set color of completed
					for _, n := range road.nodes {
						for _, side := range n.sides {
							n.tile.Teams[side] = winners
						}
					}
				}
			}
		}
	}
	// score completed cloister
	cloisterLocationsToCheck := [][]int{
		{s.lastPlacedTiles[team].X, s.lastPlacedTiles[team].Y},
		{s.lastPlacedTiles[team].X + 1, s.lastPlacedTiles[team].Y},
		{s.lastPlacedTiles[team].X - 1, s.lastPlacedTiles[team].Y},
		{s.lastPlacedTiles[team].X, s.lastPlacedTiles[team].Y + 1},
		{s.lastPlacedTiles[team].X, s.lastPlacedTiles[team].Y - 1},
		{s.lastPlacedTiles[team].X + 1, s.lastPlacedTiles[team].Y + 1},
		{s.lastPlacedTiles[team].X + 1, s.lastPlacedTiles[team].Y - 1},
		{s.lastPlacedTiles[team].X - 1, s.lastPlacedTiles[team].Y + 1},
		{s.lastPlacedTiles[team].X - 1, s.lastPlacedTiles[team].Y - 1}}
	for _, location := range cloisterLocationsToCheck {
		tile := s.board.tile(location[0], location[1])
		if tile != nil && tile.Center == Cloister {
			count, err := s.board.tilesSurroundingCloister(location[0], location[1])
			if err != nil {
				return err
			}
			if count == 8 {
				for _, token := range s.boardTokens {
					if token.Type == Monk && token.X == location[0] && token.Y == location[1] {
						// add to score
						s.scores[token.Team] += count + 1
						// remove inside from board and add back to tokens pile
						s.tokens[token.Team]++
						s.boardTokens = removeTokens(s.boardTokens, token)
						// set color of completed
						tile.CenterTeam = token.Team
						break
					}
				}
			}
		}
	}

	// draw tile for player
	if !s.deck.Empty() {
		tile, _ := s.deck.Draw()
		s.playTiles[s.turn] = tile
	}

	tilesInHands := 0
	for _, tile := range s.playTiles {
		if tile != nil {
			tilesInHands++
		}
	}

	if tilesInHands > 0 {
		// next turn
		for idx, team := range s.teams {
			if team == s.turn {
				s.turn = s.teams[(idx+1)%len(s.teams)]
				break
			}
		}

		// edge case where play tile isn't playable so re-draw
		if !s.board.playable(s.playTiles[s.turn]) {
			if !s.deck.Empty() {
				tried := []*tile{s.playTiles[s.turn]}
				retryLimit := s.deck.Size()
				for i := 0; i < retryLimit; i++ {
					tile, _ := s.deck.Draw()
					if s.board.playable(tile) {
						s.playTiles[s.turn] = tile
						s.deck.Add(tried...)
						tried = nil
						break
					} else {
						tried = append(tried, tile)
					}
				}
				if tried != nil {
					// edge case where no tile in the deck is playable so end the game instead
					if err := s.scoreAndClean(); err != nil {
						return err
					}
				}
			} else {
				// edge case where tiles still remain but cannot be played so end the game instead
				if err := s.scoreAndClean(); err != nil {
					return err
				}
			}
		}
	} else {
		// all tiles have been played so score
		if err := s.scoreAndClean(); err != nil {
			return err
		}
	}
	return nil
}

func (s *state) scoreAndClean() error {
	// score incomplete roads, cities, and cloister and score farms
	for len(s.boardTokens) > 0 {
		token := s.boardTokens[0]
		switch token.Type {
		case Knight:
			city, err := s.board.generateCity(token.X, token.Y, token.Side)
			if err != nil {
				return err
			}
			// score and add points
			points, err := scoreCity(city)
			if err != nil {
				return err
			}
			inside := tokensInStructure(s.boardTokens, city)
			winners := pointsWinners(inside)
			for _, winner := range winners {
				s.scores[winner] += points
			}
			// remove inside from board and add back to tokens pile
			for _, token := range inside {
				s.tokens[token.Team]++
			}
			s.boardTokens = removeTokens(s.boardTokens, inside...)
			// set color of incomplete
			for _, n := range city.nodes {
				for _, side := range n.sides {
					n.tile.Teams[side] = winners
				}
			}
		case Thief:
			road, err := s.board.generateRoad(token.X, token.Y, token.Side)
			if err != nil {
				return err
			}
			// score and add points
			points, err := scoreRoad(road)
			if err != nil {
				return err
			}
			inside := tokensInStructure(s.boardTokens, road)
			winners := pointsWinners(inside)
			for _, winner := range winners {
				s.scores[winner] += points
			}
			// remove inside from board and add back to tokens pile
			for _, token := range inside {
				s.tokens[token.Team]++
			}
			s.boardTokens = removeTokens(s.boardTokens, inside...)
			// set color of incomplete
			for _, n := range road.nodes {
				for _, side := range n.sides {
					n.tile.Teams[side] = winners
				}
			}
		case Monk:
			tile := s.board.tile(token.X, token.Y)
			if tile != nil && tile.Center == Cloister {
				count, err := s.board.tilesSurroundingCloister(token.X, token.Y)
				if err != nil {
					return err
				}
				s.scores[token.Team] += count + 1
				// remove inside from board and add back to tokens pile
				s.tokens[token.Team]++
				s.boardTokens = removeTokens(s.boardTokens, token)
				// set color of incomplete
				tile.CenterTeam = token.Team
			}
		case Farmer:
			farm, err := s.board.generateFarm(token.X, token.Y, token.Side)
			if err != nil {
				return err
			}
			// score and add points
			points, err := scoreFarm(farm, s.board.completeCities)
			if err != nil {
				return err
			}
			inside := tokensInStructure(s.boardTokens, farm)
			winners := pointsWinners(inside)
			for _, winner := range winners {
				s.scores[winner] += points
			}
			// remove inside from board and add back to tokens pile
			for _, token := range inside {
				s.tokens[token.Team]++
			}
			s.boardTokens = removeTokens(s.boardTokens, inside...)
			// set color of farmland
			for _, n := range farm.nodes {
				// get number of city sides
				numCities := 0
				for _, section := range n.tile.Sides {
					if section == City {
						numCities++
					}
				}
				// edge case where two adjacent disconnected city sections leads to uncolored farmland between them
				if !n.tile.ConnectedCitySides && numCities == 2 {
					for _, farmSide := range FarmSides {
						n.tile.FarmTeams[farmSide] = winners
					}
				} else {
					// otherwise, do normal coloring
					for _, farmSide := range n.sides {
						n.tile.FarmTeams[farmSide] = winners
					}
				}
			}
		}
	}
	// winner is team with the highest score
	max := 0
	winners := make([]string, 0)
	for p, score := range s.scores {
		if score > max {
			max = score
			winners = []string{p}
		} else if score == max {
			winners = append(winners, p)
		}
	}
	s.winners = winners
	return nil
}

// scoreLastTile finds how the points from the last time and token placement
func (s *state) scoreLastTile(team string) (float64, error) {
	if len(s.board.board) == 0 {
		return 0, fmt.Errorf("cannot score when no tile has been placed")
	}
	tile := s.board.board[len(s.board.board)-1]

	// only score farms if farmer was placed on tile
	scoreFarms := false
	if len(s.boardTokens) > 0 {
		last := s.boardTokens[len(s.boardTokens)-1]
		if last.X == tile.X && last.Y == tile.Y && last.Type == Farmer {
			scoreFarms = true
		}
	}

	seen := make(map[string]bool)
	for _, side := range Sides {
		seen[side] = false
	}
	for _, side := range FarmSides {
		seen[side] = false
	}

	structures := make([]*structure, 0)

	for _, side := range Sides {
		if !seen[side] {
			switch tile.Sides[side] {
			case City:
				city, err := s.board.generateCity(tile.X, tile.Y, side)
				if err != nil {
					return 0, err
				}
				for _, s := range city.nodes[0].sides {
					seen[s] = true
				}
				structures = append(structures, city)
			case Road:
				road, err := s.board.generateRoad(tile.X, tile.Y, side)
				if err != nil {
					return 0, err
				}
				for _, s := range road.nodes[0].sides {
					seen[s] = true
				}
				structures = append(structures, road)

				if scoreFarms && !seen[side+FarmNotchA] {
					farm, err := s.board.generateFarm(tile.X, tile.Y, side+FarmNotchA)
					if err != nil {
						return 0, err
					}
					for _, s := range farm.nodes[0].sides {
						seen[s] = true
					}
					structures = append(structures, farm)
				}
				if scoreFarms && !seen[side+FarmNotchB] {
					farm, err := s.board.generateFarm(tile.X, tile.Y, side+FarmNotchB)
					if err != nil {
						return 0, err
					}
					for _, s := range farm.nodes[0].sides {
						seen[s] = true
					}
					structures = append(structures, farm)
				}
			case Farm:
				if scoreFarms && !seen[side+FarmNotchA] {
					farm, err := s.board.generateFarm(tile.X, tile.Y, side+FarmNotchA)
					if err != nil {
						return 0, err
					}
					for _, s := range farm.nodes[0].sides {
						seen[s] = true
					}
					structures = append(structures, farm)
				}
			}
		}
	}

	var points float64
	if tile.Center == Cloister {
		found := false
		for _, tok := range s.boardTokens {
			if tile.X == tok.X && tile.Y == tok.Y && tok.Type == Monk {
				found = true
				break
			}
		}
		if found {
			count, err := s.board.tilesSurroundingCloister(tile.X, tile.Y)
			if err != nil {
				return 0, err
			}
			// prioritize cloisters
			points += float64((count + 1) * 2)
		}
	}

	for _, structure := range structures {
		inside := tokensInStructure(s.boardTokens, structure)
		winners := pointsWinners(inside)
		if slices.Contains(winners, team) {
			switch structure.typ {
			case City:
				pts, err := scoreCity(structure)
				if err != nil {
					return 0, err
				}
				points += float64(pts)
			case Road:
				pts, err := scoreRoad(structure)
				if err != nil {
					return 0, err
				}
				points += float64(pts)
			case Farm:
				pts, err := scoreFarm(structure, s.board.completeCities)
				if err != nil {
					return 0, err
				}
				// halve the score for farms to prevent over farming
				points += float64(pts) / 2.0
			}
		}
	}
	return points, nil
}

func (s *state) score(team string) (float64, error) {
	var points float64

	tokens := make([]*token, 0)
	for _, token := range s.boardTokens {
		if token.Team == team {
			tokens = append(tokens, token)
		}
	}

	seen := make(map[string]bool)
	for len(seen) < len(tokens) {
		var token *token
		for _, t := range tokens {
			id := fmt.Sprintf("%d,%d", t.X, t.Y)
			if !seen[id] {
				token = t
				seen[id] = true
				break
			}
		}
		switch token.Type {
		case Knight:
			city, err := s.board.generateCity(token.X, token.Y, token.Side)
			if err != nil {
				return 0, err
			}
			inside := tokensInStructure(s.boardTokens, city)
			winners := pointsWinners(inside)
			if slices.Contains(winners, team) {
				pts, err := scoreCity(city)
				if err != nil {
					return 0, err
				}
				points += float64(pts)
			}
		case Thief:
			road, err := s.board.generateRoad(token.X, token.Y, token.Side)
			if err != nil {
				return 0, err
			}
			inside := tokensInStructure(s.boardTokens, road)
			winners := pointsWinners(inside)
			if slices.Contains(winners, team) {
				pts, err := scoreRoad(road)
				if err != nil {
					return 0, err
				}
				points += float64(pts)
			}
		case Monk:
			tile := s.board.tile(token.X, token.Y)
			if tile != nil && tile.Center == Cloister {
				count, err := s.board.tilesSurroundingCloister(token.X, token.Y)
				if err != nil {
					return 0, err
				}
				points += float64(count + 1)
			}
		case Farmer:
			farm, err := s.board.generateFarm(token.X, token.Y, token.Side)
			if err != nil {
				return 0, err
			}
			inside := tokensInStructure(s.boardTokens, farm)
			winners := pointsWinners(inside)
			if slices.Contains(winners, team) {
				pts, err := scoreFarm(farm, s.board.completeCities)
				if err != nil {
					return 0, err
				}
				points += float64(pts)
			}
		}
	}
	return points, nil
}

func (s *state) actions() []*qg.Action {
	targets := make([]*qg.Action, 0)
	if s.playTiles[s.turn] != nil {
		// find all valid places to play tile in every rotation
		for i := 0; i < 4; i++ {
			tile := s.playTiles[s.turn].copy()
			for j := 0; j < i; j++ {
				tile.RotateRight()
			}
			emptySpaces := s.board.getEmptySpaces()
			for _, emptySpace := range emptySpaces {
				valid := true
				for _, side := range Sides {
					if emptySpace.adjacent[side] != nil && emptySpace.adjacent[side].Sides[AcrossSide[side]] != tile.Sides[side] {
						valid = false
					}
				}
				if valid {
					targets = append(targets, &qg.Action{
						Team: s.turn,
						Type: PlaceTileAction,
						Details: PlaceTileDetails{
							X: emptySpace.X,
							Y: emptySpace.Y,

							Tile: Tile{
								Top:             tile.Sides[SideTop],
								Right:           tile.Sides[SideRight],
								Bottom:          tile.Sides[SideBottom],
								Left:            tile.Sides[SideLeft],
								Center:          tile.Center,
								ConnectedCities: tile.ConnectedCitySides,
								Banner:          tile.Banner,
							},
						},
					})
				}
			}
		}
	} else {
		// find all valid places to play token
		targets = append(targets, &qg.Action{
			Team: s.turn,
			Type: PlaceTokenAction,
			Details: PlaceTokenDetails{
				Pass: true,
			},
		})
		if s.lastPlacedTiles[s.turn].Center == Cloister {
			targets = append(targets, &qg.Action{
				Team: s.turn,
				Type: PlaceTokenAction,
				Details: PlaceTokenDetails{
					X:    s.lastPlacedTiles[s.turn].X,
					Y:    s.lastPlacedTiles[s.turn].Y,
					Type: Monk,
					Side: SideCenter,
				},
			})
		}
		for _, side := range Sides {
			switch s.lastPlacedTiles[s.turn].Sides[side] {
			case Road:
				// check if road is already claimed
				road, _ := s.board.generateRoad(s.lastPlacedTiles[s.turn].X, s.lastPlacedTiles[s.turn].Y, side)
				if len(tokensInStructure(s.boardTokens, road)) == 0 {
					targets = append(targets, &qg.Action{
						Team: s.turn,
						Type: PlaceTokenAction,
						Details: PlaceTokenDetails{
							X:    s.lastPlacedTiles[s.turn].X,
							Y:    s.lastPlacedTiles[s.turn].Y,
							Type: Thief,
							Side: side,
						},
					})
				}
				// check if farmland A is claimed
				farm, _ := s.board.generateFarm(s.lastPlacedTiles[s.turn].X, s.lastPlacedTiles[s.turn].Y, sideToFarmSide(side, FarmNotchA))
				if len(tokensInStructure(s.boardTokens, farm)) == 0 {
					targets = append(targets, &qg.Action{
						Team: s.turn,
						Type: PlaceTokenAction,
						Details: PlaceTokenDetails{
							X:    s.lastPlacedTiles[s.turn].X,
							Y:    s.lastPlacedTiles[s.turn].Y,
							Type: Farmer,
							Side: sideToFarmSide(side, FarmNotchA),
						},
					})
				}
				// check if farmland B is claimed
				farm, _ = s.board.generateFarm(s.lastPlacedTiles[s.turn].X, s.lastPlacedTiles[s.turn].Y, sideToFarmSide(side, FarmNotchB))
				if len(tokensInStructure(s.boardTokens, farm)) == 0 {
					targets = append(targets, &qg.Action{
						Team: s.turn,
						Type: PlaceTokenAction,
						Details: PlaceTokenDetails{
							X:    s.lastPlacedTiles[s.turn].X,
							Y:    s.lastPlacedTiles[s.turn].Y,
							Type: Farmer,
							Side: sideToFarmSide(side, FarmNotchB),
						},
					})
				}
			case City:
				// check if city is already claimed
				city, _ := s.board.generateCity(s.lastPlacedTiles[s.turn].X, s.lastPlacedTiles[s.turn].Y, side)
				if len(tokensInStructure(s.boardTokens, city)) == 0 {
					targets = append(targets, &qg.Action{
						Team: s.turn,
						Type: PlaceTokenAction,
						Details: PlaceTokenDetails{
							X:    s.lastPlacedTiles[s.turn].X,
							Y:    s.lastPlacedTiles[s.turn].Y,
							Type: Knight,
							Side: side,
						},
					})
				}
			case Farm:
				// check if farmland is claimed
				farm, _ := s.board.generateFarm(s.lastPlacedTiles[s.turn].X, s.lastPlacedTiles[s.turn].Y, sideToFarmSide(side, FarmNotchA))
				if len(tokensInStructure(s.boardTokens, farm)) == 0 {
					targets = append(targets, &qg.Action{
						Team: s.turn,
						Type: PlaceTokenAction,
						Details: PlaceTokenDetails{
							X:    s.lastPlacedTiles[s.turn].X,
							Y:    s.lastPlacedTiles[s.turn].Y,
							Type: Farmer,
							Side: sideToFarmSide(side, FarmNotchA),
						},
					}, &qg.Action{
						Team: s.turn,
						Type: PlaceTokenAction,
						Details: PlaceTokenDetails{
							X:    s.lastPlacedTiles[s.turn].X,
							Y:    s.lastPlacedTiles[s.turn].Y,
							Type: Farmer,
							Side: sideToFarmSide(side, FarmNotchB),
						},
					})
				}
			}
		}
	}
	return targets
}

func (s *state) message() string {
	message := fmt.Sprintf("%s must place a tile", s.turn)
	if s.playTiles[s.turn] == nil {
		message = fmt.Sprintf("%s must place a token", s.turn)
	}
	if len(s.winners) > 0 {
		message = fmt.Sprintf("%s tie", strings.Join(s.winners, ", "))
		if len(s.winners) == 1 {
			message = fmt.Sprintf("%s wins", s.winners[0])
		}
	}
	return message
}

// get the tokens that fall in the structure
func tokensInStructure(tokens []*token, structure *structure) []*token {
	tokensInside := make([]*token, 0)
	for _, token := range tokens {
		for _, n := range structure.nodes {
			// check if token type matches section type and token on node
			if StructureTypeToTokenType[structure.typ] == token.Type &&
				n.tile.X == token.X && n.tile.Y == token.Y && slices.Contains(n.sides, token.Side) {
				tokensInside = append(tokensInside, token)
			}
		}
	}
	return tokensInside
}

// create a new list that has removed toRemove from original
func removeTokens(original []*token, toRemove ...*token) []*token {
	newTokens := make([]*token, 0)
	for _, token := range original {
		found := false
		for _, rm := range toRemove {
			if token.X == rm.X && token.Y == rm.Y && token.Side == rm.Side &&
				token.Type == rm.Type && token.Team == rm.Team {
				found = true
			}
		}
		if !found {
			newTokens = append(newTokens, token)
		}
	}
	return newTokens
}

// given a list of tokens, get the teams(s) with the most tokens
func pointsWinners(tokens []*token) []string {
	max := 0
	winners := make([]string, 0)
	tally := make(map[string]int)
	for _, token := range tokens {
		tally[token.Team]++
	}
	for team, count := range tally {
		if count > max {
			winners = []string{team}
			max = count
		} else if count == max {
			winners = append(winners, team)
		}
	}
	sort.Strings(winners)
	return winners
}
