package carcassonne

// Token types
const (
	Farmer = "farmer"
	Knight = "knight"
	Thief  = "thief"
	Monk   = "monk"
)

var TokenTypes = []string{Farmer, Knight, Thief, Monk}

type token struct {
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Team string `json:"team"`
	Type string `json:"type"` // Farmer, Knight, Thief, Monk
	Side string `json:"side"` // normal side if Knight or Thief, farm side if Farmer, empty if Monk
}

func newToken(x, y int, team, typ, side string) *token {
	return &token{
		X:    x,
		Y:    y,
		Team: team,
		Type: typ,
		Side: side,
	}
}
