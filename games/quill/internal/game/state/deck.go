package state

import (
	c "github.com/quibbble/quibbble-controller/pkg/collection"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const (
	InitDeckSize = 30
)

type Deck struct {
	c.Collection[ICard]
}

func NewEmptyDeck(seed int64) *Deck {
	return &Deck{
		Collection: *c.NewCollection[ICard](seed),
	}
}

func NewDeck(seed int64, build BuildCard, player uuid.UUID, ids ...string) (*Deck, error) {
	deck := NewEmptyDeck(seed)
	for _, id := range ids {
		card, err := build(id, player, false)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		deck.Add(card)
	}
	return deck, nil
}
