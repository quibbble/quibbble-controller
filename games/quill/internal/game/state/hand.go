package state

import (
	c "github.com/quibbble/quibbble-controller/pkg/collection"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

const (
	InitHandSize = 5
	MaxHandSize  = 10
)

type Hand struct {
	c.Collection[ICard]
}

func NewHand(seed int64, card ...ICard) *Hand {
	collection := c.NewCollection[ICard](seed)
	collection.Add(card...)
	return &Hand{
		Collection: *collection,
	}
}

func (h *Hand) GetCard(card uuid.UUID) (ICard, error) {
	for _, it := range h.GetItems() {
		if it.GetUUID() == card {
			return it, nil
		}
	}
	return nil, ErrNotFound(card)
}

func (h *Hand) RemoveCard(card uuid.UUID) error {
	for i, it := range h.GetItems() {
		if it.GetUUID() == card {
			if err := h.Remove(i); err != nil {
				return errors.Wrap(err)
			}
			return nil
		}
	}
	return ErrNotFound(card)
}
