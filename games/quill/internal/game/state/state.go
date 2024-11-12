package state

import (
	"math/rand"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

type State struct {
	Seed int64
	Rand *rand.Rand
	Gen  *uuid.Gen

	Turn   int
	Teams  []uuid.UUID
	Winner *uuid.UUID

	Board *Board

	Deck    map[uuid.UUID]*Deck
	Discard map[uuid.UUID]*Deck
	Trash   map[uuid.UUID]*Deck
	Hand    map[uuid.UUID]*Hand
	Mana    map[uuid.UUID]*Mana
	Recycle map[uuid.UUID]int
	Sacked  map[uuid.UUID]bool

	*en.Builders
	BuildCard
}

func NewState(seed int64, buildCard BuildCard, engineBuilders *en.Builders, player1, player2 uuid.UUID, deck1, deck2 []string) (*State, error) {
	r := rand.New(rand.NewSource(seed))
	gen := uuid.NewGen(r)

	board, err := NewBoard(buildCard, gen, player1, player2)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	if len(deck1) != InitDeckSize || len(deck2) != InitDeckSize {
		return nil, errors.Errorf("decks must be of size %d", InitDeckSize)
	}

	d1, err := NewDeck(seed, buildCard, player1, deck1...)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	d1.Shuffle()
	d2, err := NewDeck(seed, buildCard, player2, deck2...)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	d2.Shuffle()
	d2.Shuffle()

	hand1 := make([]ICard, 0)
	hand2 := make([]ICard, 0)
	for i := 0; i < InitHandSize; i++ {
		card1, err := d1.Draw()
		if err != nil {
			return nil, errors.Wrap(err)
		}
		card2, err := d2.Draw()
		if err != nil {
			return nil, errors.Wrap(err)
		}
		hand1 = append(hand1, *card1)
		hand2 = append(hand2, *card2)
	}

	return &State{
		Seed: seed,
		Rand: r,
		Gen:  gen,

		Turn:   0,
		Teams:  []uuid.UUID{player1, player2},
		Winner: nil,

		Board:   board,
		Deck:    map[uuid.UUID]*Deck{player1: d1, player2: d2},
		Discard: map[uuid.UUID]*Deck{player1: NewEmptyDeck(seed), player2: NewEmptyDeck(seed)},
		Trash:   map[uuid.UUID]*Deck{player1: NewEmptyDeck(seed), player2: NewEmptyDeck(seed)},
		Hand:    map[uuid.UUID]*Hand{player1: NewHand(seed, hand1...), player2: NewHand(seed, hand2...)},
		Mana:    map[uuid.UUID]*Mana{player1: NewMana(), player2: NewMana()},
		Recycle: map[uuid.UUID]int{player1: 0, player2: 0},
		Sacked:  map[uuid.UUID]bool{player1: false, player2: false},

		Builders:  engineBuilders,
		BuildCard: buildCard,
	}, nil
}

func (s *State) GetTurn() uuid.UUID {
	return s.Teams[s.Turn%len(s.Teams)]
}

func (s *State) GetCard(uuid uuid.UUID) ICard {
	for _, hand := range s.Hand {
		for _, card := range hand.GetItems() {
			if card.GetUUID() == uuid {
				return card
			}
		}
	}
	for _, tile := range s.Board.UUIDs {
		if tile.Unit != nil && tile.Unit.GetUUID() == uuid {
			return tile.Unit
		}
	}
	for _, deck := range s.Deck {
		for _, card := range deck.GetItems() {
			if card.GetUUID() == uuid {
				return card
			}
		}
	}
	for _, discard := range s.Discard {
		for _, card := range discard.GetItems() {
			if card.GetUUID() == uuid {
				return card
			}
		}
	}
	for _, trash := range s.Trash {
		for _, card := range trash.GetItems() {
			if card.GetUUID() == uuid {
				return card
			}
		}
	}
	return nil
}

func (s *State) GetOpponent(player uuid.UUID) uuid.UUID {
	if s.Teams[0] == player {
		return s.Teams[1]
	}
	return s.Teams[0]
}

func (s *State) GameOver() bool {
	return s.Winner != nil
}
