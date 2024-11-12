package state

import (
	"github.com/quibbble/quibbble-controller/games/quill/parse"
)

var (
	CardIDByCostMap = map[int][]string{}
	CardIDByTypeMap = map[string][]string{}
)

var idToType = map[rune]string{
	'S': "Spell",
	'U': "Unit",
	'I': "Item",
}

func init() {
	cards, err := parse.AllCards()
	if err != nil {
		panic(err)
	}
	for _, id := range cards {
		card, err := parse.ParseCard(id)
		if err == parse.ErrNotEnabled {
			continue
		}
		if err != nil {
			panic(err)
		}
		l, ok := CardIDByCostMap[card.GetCost()]
		if !ok {
			CardIDByCostMap[card.GetCost()] = []string{card.GetID()}
		} else {
			CardIDByCostMap[card.GetCost()] = append(l, card.GetID())
		}

		typ := idToType[rune(card.GetID()[0])]
		l, ok = CardIDByTypeMap[typ]
		if !ok {
			CardIDByTypeMap[typ] = []string{card.GetID()}
		} else {
			CardIDByTypeMap[typ] = append(l, card.GetID())
		}
	}
}
