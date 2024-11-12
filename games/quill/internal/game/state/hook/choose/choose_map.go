package choose

import (
	"context"
	"reflect"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

var ChooseMap map[string]struct {
	Type     reflect.Type
	Retrieve func(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error)
}

func init() {
	ChooseMap = map[string]struct {
		Type     reflect.Type
		Retrieve func(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error)
	}{
		AdjacentChoice: {
			Type:     reflect.TypeOf(&AdjacentArgs{}),
			Retrieve: RetrieveAdjacent,
		},
		CardIDByCostChoice: {
			Type:     reflect.TypeOf(&CardIDByCostArgs{}),
			Retrieve: RetrieveCardIDByCost,
		},
		CardIDByTypeChoice: {
			Type:     reflect.TypeOf(&CardIDByTypeArgs{}),
			Retrieve: RetrieveCardIDByType,
		},
		CodexChoice: {
			Type:     reflect.TypeOf(&CodexArgs{}),
			Retrieve: RetrieveCodex,
		},
		CardTraitsChoice: {
			Type:     reflect.TypeOf(&CardTraitsArgs{}),
			Retrieve: RetrieveCardTraits,
		},
		CompositeChoice: {
			Type:     reflect.TypeOf(&CompositeArgs{}),
			Retrieve: RetrieveComposite,
		},
		ConnectedChoice: {
			Type:     reflect.TypeOf(&ConnectedArgs{}),
			Retrieve: RetrieveConnected,
		},
		CurrentPlayerChoice: {
			Type:     reflect.TypeOf(&CurrentPlayerArgs{}),
			Retrieve: RetrieveCurrentPlayer,
		},
		HookEventTileChoice: {
			Type:     reflect.TypeOf(&HookEventTileArgs{}),
			Retrieve: RetrieveHookTileUnit,
		},
		HookEventUnitChoice: {
			Type:     reflect.TypeOf(&HookEventUnitArgs{}),
			Retrieve: RetrieveHookEventUnit,
		},
		ItemHolderChoice: {
			Type:     reflect.TypeOf(&ItemHolderArgs{}),
			Retrieve: RetrieveItemHolder,
		},
		OpposingPlayerChoice: {
			Type:     reflect.TypeOf(&OpposingPlayerArgs{}),
			Retrieve: RetrieveOpposingPlayer,
		},
		OwnedTilesChoice: {
			Type:     reflect.TypeOf(&OwnedTilesArgs{}),
			Retrieve: RetrieveOwnedTiles,
		},
		OwnedUnitsChoice: {
			Type:     reflect.TypeOf(&OwnedUnitsArgs{}),
			Retrieve: RetrieveOwnedUnits,
		},
		OwnerChoice: {
			Type:     reflect.TypeOf(&OwnerArgs{}),
			Retrieve: RetrieveOwner,
		},
		RandomChoice: {
			Type:     reflect.TypeOf(&RandomArgs{}),
			Retrieve: RetrieveRandom,
		},
		RangedChoice: {
			Type:     reflect.TypeOf(&RangedArgs{}),
			Retrieve: RetrieveRanged,
		},
		SelfChoice: {
			Type:     reflect.TypeOf(&SelfArgs{}),
			Retrieve: RetrieveSelf,
		},
		TargetChoice: {
			Type:     reflect.TypeOf(&TargetArgs{}),
			Retrieve: RetrieveTarget,
		},
		TilesChoice: {
			Type:     reflect.TypeOf(&TilesArgs{}),
			Retrieve: RetrieveTiles,
		},
		TraitEventTileChoice: {
			Type:     reflect.TypeOf(&TraitEventTileArgs{}),
			Retrieve: RetrieveTraitEventTile,
		},
		TraitEventUnitChoice: {
			Type:     reflect.TypeOf(&TraitEventUnitArgs{}),
			Retrieve: RetrieveTraitEventUnit,
		},
		UnitsChoice: {
			Type:     reflect.TypeOf(&UnitsArgs{}),
			Retrieve: RetrieveUnits,
		},
		UUIDChoice: {
			Type:     reflect.TypeOf(&UUIDArgs{}),
			Retrieve: RetrieveUUID,
		},
	}
}
