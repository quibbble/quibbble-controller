package state

type Mana struct {
	Amount     int
	BaseAmount int
}

func NewMana() *Mana {
	return &Mana{
		Amount:     1,
		BaseAmount: 1,
	}
}
