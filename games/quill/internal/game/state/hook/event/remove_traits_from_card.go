package event

import (
	"context"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
)

const RemoveTraitsFromCard = "RemoveTraitsFromCard"

type RemoveTraitsFromCardArgs struct {
	ChooseTraits parse.Choose
	ChooseCard   parse.Choose
}

func RemoveTraitsFromCardAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*RemoveTraitsFromCardArgs)
	traitChoices, err := ch.GetChoices(ctx, a.ChooseTraits, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	for _, choice := range traitChoices {
		event, err := NewEvent(state.Gen.New(en.EventUUID), RemoveTraitFromCard, RemoveTraitFromCardArgs{
			ChooseTrait: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: choice,
				},
			},
			ChooseCard: a.ChooseCard,
		})
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(ctx, event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
