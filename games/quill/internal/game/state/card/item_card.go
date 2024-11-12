package card

import (
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

type ItemCard struct {
	*Card

	// UnitCard that is holding this item
	Holder *uuid.UUID

	// Traits applied to a unit when an item is held
	HeldTraits []st.ITrait
}

func (c *ItemCard) AddTrait(trait st.ITrait) error {
	return c.Card.addTrait(trait, c)
}

func (c *ItemCard) RemoveTrait(trait uuid.UUID) error {
	return c.Card.removeTrait(trait, c)
}

func (c *ItemCard) Reset(build st.BuildCard) {
	card, _ := build(c.GetID(), c.Player, c.Token)
	item := card.(*ItemCard)
	item.UUID = c.UUID
	c = item
}
