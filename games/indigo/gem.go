package indigo

const (
	Amber    = "Amber"
	Emerald  = "Emerald"
	Sapphire = "Sapphire"
)

var (
	colorToPoints = map[string]int{
		Amber:    1,
		Emerald:  2,
		Sapphire: 3,
	}
)

type gem struct {
	Color       string
	Edge        string
	Row, Column int
	collided    bool
	gateway     *gateway
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
