package card

import (
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const (
	Unit          = "Unit"
	BaseUnit      = "Base"
	StructureUnit = "Structure"
	CreatureUnit  = "Creature"
)

type UnitCard struct {
	*Card

	Type                   string
	DamageType             string
	Attack                 int
	Health                 int
	Cooldown, BaseCooldown int
	Movement, BaseMovement int
	Codex                  string

	// Items that apply held traits to this card
	Items []*ItemCard
}

func (c *UnitCard) AddTrait(trait st.ITrait) error {
	return c.Card.addTrait(trait, c)
}

func (c *UnitCard) RemoveTrait(trait uuid.UUID) error {
	return c.Card.removeTrait(trait, c)
}

func (c *UnitCard) AddItem(item *ItemCard) error {
	item.Holder = &c.UUID
	c.Items = append(c.Items, item)
	return nil
}

func (c *UnitCard) GetItem(item uuid.UUID) (*ItemCard, error) {
	idx := -1
	var itm *ItemCard
	for i, it := range c.Items {
		if it.UUID == item {
			idx = i
			itm = it
		}
	}
	if idx < 0 {
		return nil, errors.Errorf("'%s' not found on unit", item)
	}
	return itm, nil
}

func (c *UnitCard) RemoveItem(item uuid.UUID) error {
	_, err := c.GetAndRemoveItem(item)
	return err
}

func (c *UnitCard) GetAndRemoveItem(item uuid.UUID) (*ItemCard, error) {
	idx := -1
	var itm *ItemCard
	for i, it := range c.Items {
		if it.UUID == item {
			idx = i
			itm = it
		}
	}
	if idx < 0 {
		return nil, errors.Errorf("'%s' not found on unit", item)
	}
	c.Items = append(c.Items[:idx], c.Items[idx+1:]...)
	itm.Holder = nil
	return itm, nil
}

// CheckCodex checks whether the unit may move/attack from x1, y1 to x2, y2 with it's current codex
func (c *UnitCard) CheckCodex(x1, y1, x2, y2 int) bool {
	x := x2 - x1
	y := y2 - y1

	if (x < -1 || x > 1 || y < -1 || y > 1) || (x == 0 && y == 0) {
		return false
	}

	check := (x == 0 && y == 1 && c.Codex[0] == '1') || // up
		(x == 0 && y == -1 && c.Codex[1] == '1') || // down
		(x == -1 && y == 0 && c.Codex[2] == '1') || // left
		(x == 1 && y == 0 && c.Codex[3] == '1') || // right
		(x == -1 && y == 1 && c.Codex[4] == '1') || // upper-left
		(x == 1 && y == -1 && c.Codex[5] == '1') || // lower-right
		(x == -1 && y == -1 && c.Codex[6] == '1') || // lower-left
		(x == 1 && y == 1 && c.Codex[7] == '1') // upper-right

	return check
}

func (c *UnitCard) Reset(build st.BuildCard) {
	card, _ := build(c.GetID(), c.Player, c.Token)
	unit := card.(*UnitCard)
	unit.UUID = c.UUID
	c = unit
}
