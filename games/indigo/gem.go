package indigo

const (
	Amber    = "amber"
	Emerald  = "emerald"
	Sapphire = "sapphire"
)

var (
	colorToPoints = map[string]int{
		Amber:    1,
		Emerald:  2,
		Sapphire: 3,
	}
)

type gem struct {
	Color    string `json:"color"`
	Edge     string `json:"edge"`
	Row      int    `json:"row"`
	Column   int    `json:"col"`
	collided bool
	gateway  *gateway
}

func newGem(color, edge string, row, column int) *gem {
	return &gem{
		Color:    color,
		Edge:     edge,
		Row:      row,
		Column:   column,
		collided: false,
		gateway:  nil,
	}
}
