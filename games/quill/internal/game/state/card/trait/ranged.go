package trait

import "math"

const (
	RangedTrait = "Ranged"
)

type RangedArgs struct {
	Amount int
}

func (a *RangedArgs) CheckRange(x1, y1, x2, y2 int) bool {
	x := int(math.Abs(float64(x2 - x1)))
	y := int(math.Abs(float64(y2 - y1)))

	if x > a.Amount || y > a.Amount {
		return false
	}
	return true
}
