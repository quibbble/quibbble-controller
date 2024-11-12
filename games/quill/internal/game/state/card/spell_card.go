package card

import (
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

type SpellCard struct {
	*Card
}

func (c *SpellCard) AddTrait(trait st.ITrait) error {
	return c.Card.addTrait(trait, c)
}

func (c *SpellCard) RemoveTrait(trait uuid.UUID) error {
	return c.Card.removeTrait(trait, c)
}

func (c *SpellCard) Reset(build st.BuildCard) {
	card, _ := build(c.GetID(), c.Player, c.Token)
	spell := card.(*SpellCard)
	spell.UUID = c.UUID
	c = spell
}
