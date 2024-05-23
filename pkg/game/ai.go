package game

import (
	"math"
)

type GameAI interface {
	// Gives a number between -inf and inf
	// -inf means this action guarantees a lose
	// inf means this action guarantees a win
	Score(g Game, team string) float64
}

// alphabeta is a variation on alpha beta pruning.
// This takes into account games with mutiple actions
// per turn as well as games with more than two players.
func alphabeta(b GameBuilder, ai GameAI, node Game, depth int, alpha, beta float64, team string) (score float64, action *Action, err error) {
	snapshotJSON, err := node.GetSnapshotJSON()
	if err != nil {
		return 0, nil, err
	}
	snapshotQGN, err := node.GetSnapshotQGN()
	if err != nil {
		return 0, nil, err
	}
	if depth == 0 || len(snapshotJSON.Winners) > 0 {
		return ai.Score(node, team), nil, nil
	}
	if team == snapshotJSON.Turn {
		value := math.Inf(-1)
		for _, a := range snapshotJSON.Actions {
			copy, err := b.Create(snapshotQGN)
			if err != nil {
				return 0, nil, err
			}
			if err := copy.Do(a); err != nil {
				return 0, nil, err
			}

			snapshotJSONCopy, err := copy.GetSnapshotJSON()
			if err != nil {
				return 0, nil, err
			}
			newDepth := depth
			if snapshotJSON.Turn != snapshotJSONCopy.Turn {
				newDepth -= 1
			}
			v, _, err := alphabeta(b, ai, copy, newDepth, alpha, beta, team)
			if err != nil {
				return 0, nil, err
			}

			if v >= value {
				value = v
				action = a
			}
			if value > beta {
				break
			}
			alpha = math.Max(alpha, value)
		}
		return value, action, nil
	} else {
		value := math.Inf(1)
		for _, a := range snapshotJSON.Actions {
			copy, err := b.Create(snapshotQGN)
			if err != nil {
				return 0, nil, err
			}
			if err := copy.Do(a); err != nil {
				return 0, nil, err
			}

			snapshotJSONCopy, err := copy.GetSnapshotJSON()
			if err != nil {
				return 0, nil, err
			}
			newDepth := depth
			if snapshotJSON.Turn != snapshotJSONCopy.Turn {
				newDepth -= 1
			}
			v, _, err := alphabeta(b, ai, copy, newDepth, alpha, beta, team)
			if err != nil {
				return 0, nil, err
			}

			if v <= value {
				value = v
				action = a
			}
			if value < alpha {
				break
			}
			beta = math.Min(beta, value)
		}
		return value, action, nil
	}
}
