package indigo

import (
	"fmt"
	"slices"
	"strings"
)

const (
	rows       = 9
	minColumns = 5
	maxColumns = 9
)

var (
	edgeToEdge = map[string]string{A: D, B: E, C: F, D: A, E: B, F: C}
)

type board struct {
	Tiles    [][]*tile  `json:"tiles"` // 0,0 is upper left most tile
	Gateways []*gateway `json:"gateways"`
	Gems     []*gem     `json:"gems"`
}

func newBoard(teams []string) *board {
	var b = make([][]*tile, rows)
	columns := minColumns
	for i := 0; i < rows; i++ {
		b[i] = make([]*tile, columns)
		if i < 4 {
			columns++
		} else {
			columns--
		}
	}

	// place treasure tiles
	for edges, location := range initTreasureTiles {
		b[location[0]][location[1]] = newTreasureTile(edges)
	}

	// create gateways
	gateways := make([]*gateway, 0)
	for edges, teamsIdxs := range numTeamsToGatewayOwnership[len(teams)] {
		owners := make([]string, 0)
		for _, idx := range teamsIdxs {
			owners = append(owners, teams[idx])
		}
		gateways = append(gateways, newGateway(initGateways[edges], edges, owners...))
	}

	// create gems
	gems := make([]*gem, 0)
	for _, gem := range initGems {
		gems = append(gems, newGem(gem[0].(string), gem[1].(string), gem[2].(int), gem[3].(int)))
	}

	return &board{
		Tiles:    b,
		Gems:     gems,
		Gateways: gateways,
	}
}

func (b *board) place(tile *tile, row, col int) error {
	if row < 0 || col < 0 || row >= rows || col >= len(b.Tiles[row]) {
		return fmt.Errorf("index out of bounds")
	}
	if b.Tiles[row][col] != nil {
		return fmt.Errorf("tile already exists at (%d, %d)", row, col)
	}
	if len(tile.Paths) != 6 {
		return fmt.Errorf("invalid tile paths")
	}
	paths := []string{tile.Paths[0:2], tile.Paths[2:4], tile.Paths[4:6]}
	for _, gateway := range b.Gateways {
		for _, location := range gateway.Locations {
			if row == location[0] && col == location[1] && slices.Contains(paths, gateway.Edges) {
				return fmt.Errorf("cannot place a tile in a way that blocks a gateway")
			}
		}
	}
	b.Tiles[row][col] = tile
	return nil
}

func (b *board) moveGems(placedRow, placedCol int) ([]*gem, error) {
	moved := []*gem{}
	centerGemMoved := false

nextGem:
	for _, gem := range b.Gems {
		if gem.collided || gem.gateway != nil {
			continue
		}

		var (
			adjRow, adjCol int
			adjEdge        string
		)

		// case where tile placed adj to middle treasure tile and one gem must be moved
		if gem.Edge == Special {
			if !centerGemMoved && placedRow >= 0 && placedCol >= 0 {
				edgeMap := map[string][2]int{
					"a": {-1, -1},
					"b": {-1, 0},
					"c": {0, 1},
					"d": {1, 0},
					"e": {1, -1},
					"f": {0, -1},
				}
				for edge, loc := range edgeMap {
					if gem.Row+loc[0] == placedRow &&
						gem.Column+loc[1] == placedCol {
						adjRow = placedRow
						adjCol = placedCol
						adjEdge = edgeToEdge[edge]
						centerGemMoved = true
						break
					}
				}
				if !centerGemMoved {
					continue nextGem
				}
			} else {
				continue nextGem
			}

		}

		// base case where gem has a adj tile and must be moved
		if adjEdge == "" {
			adjRow, adjCol, adjEdge = b.getAdjacent(gem.Row, gem.Column, gem.Edge)
			if adjRow < 0 || adjRow >= len(b.Tiles) ||
				adjCol < 0 || adjCol >= len(b.Tiles[adjRow]) ||
				b.Tiles[adjRow][adjCol] == nil {
				continue nextGem
			}
		}

		// check for collision
		for _, g := range b.Gems {
			if g.Row == adjRow && g.Column == adjCol && g.Edge == adjEdge {
				gem.collided = true
				g.collided = true
				continue nextGem
			}
		}

		movedEdge, err := b.Tiles[adjRow][adjCol].GetDestination(adjEdge)
		if err != nil {
			return nil, err
		}

		gem.Row = adjRow
		gem.Column = adjCol
		gem.Edge = movedEdge

		moved = append(moved, gem)

		// check for gateway reached
	nextGateway:
		for _, gateway := range b.Gateways {
			for _, loc := range gateway.Locations {
				if loc[0] == gem.Row && loc[1] == gem.Column && strings.Contains(gateway.Edges, gem.Edge) {
					gem.gateway = gateway
					break nextGateway
				}
			}
		}
	}

	if len(moved) > 0 {
		// NOTE only gems moved the first iteration could be moved again so do not need to concat returned gems on future iterations
		_, err := b.moveGems(-1, -1)
		if err != nil {
			return nil, err
		}
	}

	return moved, nil
}

// getAdjacent returns the adjacent row, col, and edge
func (b *board) getAdjacent(row, col int, edge string) (adjRow, adjCol int, adjEdge string) {
	edgeToRowColTop := map[string][2]int{A: {-1, -1}, B: {-1, 0}, C: {0, 1}, D: {1, 1}, E: {1, 0}, F: {0, -1}}
	edgeToRowColBot := map[string][2]int{A: {-1, 0}, B: {-1, 1}, C: {0, 1}, D: {1, 0}, E: {1, -1}, F: {0, -1}}
	var edgeMap map[string][2]int
	if row < rows/2 {
		edgeMap = edgeToRowColTop
	} else if row > rows/2 {
		edgeMap = edgeToRowColBot
	} else {
		if strings.Contains("ab", edge) {
			edgeMap = edgeToRowColTop
		} else if strings.Contains("de", edge) {
			edgeMap = edgeToRowColBot
		} else {
			edgeMap = edgeToRowColTop
		}
	}
	return row + edgeMap[edge][0], col + edgeMap[edge][1], edgeToEdge[edge]
}

func (b *board) gemsInPlay() int {
	count := 0
	for _, gem := range b.Gems {
		if !gem.collided && gem.gateway == nil {
			count++
		}
	}
	return count
}
