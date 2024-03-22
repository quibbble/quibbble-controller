package tsuro

import (
	"errors"
	"math/rand"
)

type deck struct {
	deck []*tile
	rand *rand.Rand
}

func newDeck(random *rand.Rand) *deck {
	d := make([]*tile, 0)
	for _, edges := range tiles {
		t, _ := newTile(edges)
		d = append(d, t)
	}
	result := &deck{
		deck: d,
		rand: random,
	}
	result.Shuffle()
	return result
}

func (d *deck) Remove(tile *tile) error {
	for idx, t := range d.deck {
		if tile.equals(t) {
			d.deck = append(d.deck[:idx], d.deck[idx+1:]...)
			return nil
		}
	}
	return errors.New("tile not found")
}

func (d *deck) Add(tiles ...*tile) {
	d.deck = append(d.deck, tiles...)
	d.Shuffle()
}

func (d *deck) Draw() (*tile, error) {
	size := len(d.deck)
	if size <= 0 {
		return nil, errors.New("deck is empty so cannot draw")
	}
	tile := d.deck[size-1]
	d.deck = d.deck[:size-1]
	return tile, nil
}

func (d *deck) Shuffle() {
	for i := 0; i < len(d.deck); i++ {
		r := d.rand.Intn(len(d.deck))
		if i != r {
			d.deck[r], d.deck[i] = d.deck[i], d.deck[r]
		}
	}
}
