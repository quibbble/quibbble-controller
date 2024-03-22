package tsuro

import (
	"errors"
	"math"
	"math/rand"
)

type token struct {
	Row   int    `json:"row"`
	Col   int    `json:"col"`
	Notch string `json:"notch"`
}

func newToken(row, col int, notch string) *token {
	return &token{
		Row:   row,
		Col:   col,
		Notch: notch,
	}
}

func randomToken(random *rand.Rand) *token {
	options := "ABCDEFGH"
	notch := string(options[random.Intn(8)])                            // which notch to lie on
	side := random.Intn(int(math.Min(float64(rows), float64(columns)))) // which side to lie on
	var row = 0
	var col = 0
	switch notch {
	case "A", "B":
		row = 0
		col = side
	case "C", "D":
		row = side
		col = columns - 1
	case "E", "F":
		row = rows - 1
		col = side
	case "G", "H":
		row = side
		col = 0
	}
	return newToken(row, col, notch)
}

func uniqueRandomToken(tokens map[string]*token, random *rand.Rand) *token {
	token := randomToken(random)
	for _, tok := range tokens {
		if token.Row == tok.Row && token.Col == tok.Col {
			return uniqueRandomToken(tokens, random)
		}
	}
	return token
}

func (t *token) equals(t2 *token) bool {
	if t.Row == t2.Row && t.Col == t2.Col && t.Notch == t2.Notch {
		return true
	}
	return false
}

func (t *token) collided(t2 *token) bool {
	adjacent := map[string]string{"A": "F", "B": "E", "C": "H", "D": "G", "E": "B", "F": "A", "G": "D", "H": "C"}
	switch t.Notch {
	case "A", "B":
		return t2.equals(&token{Row: t.Row - 1, Col: t.Col, Notch: adjacent[t.Notch]})
	case "C", "D":
		return t2.equals(&token{Row: t.Row, Col: t.Col + 1, Notch: adjacent[t.Notch]})
	case "E", "F":
		return t2.equals(&token{Row: t.Row + 1, Col: t.Col, Notch: adjacent[t.Notch]})
	case "G", "H":
		return t2.equals(&token{Row: t.Row, Col: t.Col - 1, Notch: adjacent[t.Notch]})
	}
	return false
}

func (t *token) getAdjacent() (*token, error) {
	adjacent := map[string]string{"A": "F", "B": "E", "C": "H", "D": "G", "E": "B", "F": "A", "G": "D", "H": "C"}
	invalidErr := errors.New("invalid token")
	switch t.Notch {
	case "A", "B":
		if t.Row <= 0 {
			return nil, invalidErr
		}
		return &token{Row: t.Row - 1, Col: t.Col, Notch: adjacent[t.Notch]}, nil
	case "C", "D":
		if t.Col >= columns {
			return nil, invalidErr
		}
		return &token{Row: t.Row, Col: t.Col + 1, Notch: adjacent[t.Notch]}, nil
	case "E", "F":
		if t.Row >= rows {
			return nil, invalidErr
		}
		return &token{Row: t.Row + 1, Col: t.Col, Notch: adjacent[t.Notch]}, nil
	case "G", "H":
		if t.Col <= 0 {
			return nil, invalidErr
		}
		return &token{Row: t.Row, Col: t.Col - 1, Notch: adjacent[t.Notch]}, nil
	}
	return nil, invalidErr
}
