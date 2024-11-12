package event

import (
	"context"
	"reflect"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
)

var EventMap map[string]struct {
	Type   reflect.Type
	Affect func(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error
}

func init() {
	EventMap = map[string]struct {
		Type   reflect.Type
		Affect func(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error
	}{
		AddItemToUnitEvent: {
			Type:   reflect.TypeOf(&AddItemToUnitArgs{}),
			Affect: AddItemToUnitAffect,
		},
		AddTraitToCard: {
			Type:   reflect.TypeOf(&AddTraitToCardArgs{}),
			Affect: AddTraitToCardAffect,
		},
		AttackUnitEvent: {
			Type:   reflect.TypeOf(&AttackUnitArgs{}),
			Affect: AttackUnitAffect,
		},
		BurnCardEvent: {
			Type:   reflect.TypeOf(&BurnCardArgs{}),
			Affect: BurnCardAffect,
		},
		CooldownEvent: {
			Type:   reflect.TypeOf(&CooldownArgs{}),
			Affect: CooldownAffect,
		},
		DamageUnitEvent: {
			Type:   reflect.TypeOf(&DamageUnitArgs{}),
			Affect: DamageUnitAffect,
		},
		DamageUnitsEvent: {
			Type:   reflect.TypeOf(&DamageUnitsArgs{}),
			Affect: DamageUnitsAffect,
		},
		DiscardCardEvent: {
			Type:   reflect.TypeOf(&DiscardCardArgs{}),
			Affect: DiscardCardAffect,
		},
		DrainBaseManaEvent: {
			Type:   reflect.TypeOf(&DrainBaseManaArgs{}),
			Affect: DrainBaseManaAffect,
		},
		DrainManaEvent: {
			Type:   reflect.TypeOf(&DrainManaArgs{}),
			Affect: DrainManaAffect,
		},
		DrawCardEvent: {
			Type:   reflect.TypeOf(&DrawCardArgs{}),
			Affect: DrawCardAffect,
		},
		EndGameEvent: {
			Type:   reflect.TypeOf(&EndGameArgs{}),
			Affect: EndGameAffect,
		},
		EndTurnEvent: {
			Type:   reflect.TypeOf(&EndTurnArgs{}),
			Affect: EndTurnAffect,
		},
		GainBaseManaEvent: {
			Type:   reflect.TypeOf(&GainBaseManaArgs{}),
			Affect: GainBaseManaAffect,
		},
		GainManaEvent: {
			Type:   reflect.TypeOf(&GainManaArgs{}),
			Affect: GainManaAffect,
		},
		HealUnitEvent: {
			Type:   reflect.TypeOf(&HealUnitArgs{}),
			Affect: HealUnitAffect,
		},
		HealUnitsEvent: {
			Type:   reflect.TypeOf(&HealUnitsArgs{}),
			Affect: HealUnitsAffect,
		},
		KillUnitEvent: {
			Type:   reflect.TypeOf(&KillUnitArgs{}),
			Affect: KillUnitAffect,
		},
		ModifyUnitEvent: {
			Type:   reflect.TypeOf(&ModifyUnitArgs{}),
			Affect: ModifyUnitAffect,
		},
		ModifyUnitsEvent: {
			Type:   reflect.TypeOf(&ModifyUnitsArgs{}),
			Affect: ModifyUnitsAffect,
		},
		MoveUnitEvent: {
			Type:   reflect.TypeOf(&MoveUnitArgs{}),
			Affect: MoveUnitAffect,
		},
		PlaceUnitEvent: {
			Type:   reflect.TypeOf(&PlaceUnitArgs{}),
			Affect: PlaceUnitAffect,
		},
		PlayCardEvent: {
			Type:   reflect.TypeOf(&PlayCardArgs{}),
			Affect: PlayCardAffect,
		},
		RecallUnitEvent: {
			Type:   reflect.TypeOf(&RecallUnitArgs{}),
			Affect: RecallUnitAffect,
		},
		RecycleDeckEvent: {
			Type:   reflect.TypeOf(&RecycleDeckArgs{}),
			Affect: RecycleDeckAffect,
		},
		RefreshMovementEvent: {
			Type:   reflect.TypeOf(&RefreshMovementArgs{}),
			Affect: RefreshMovementAffect,
		},
		RemoveItemFromUnitEvent: {
			Type:   reflect.TypeOf(&RemoveItemFromUnitArgs{}),
			Affect: RemoveItemFromUnitAffect,
		},
		RemoveTraitFromCard: {
			Type:   reflect.TypeOf(&RemoveTraitFromCardArgs{}),
			Affect: RemoveTraitFromCardAffect,
		},
		RemoveTraitsFromCard: {
			Type:   reflect.TypeOf(&RemoveTraitsFromCardArgs{}),
			Affect: RemoveTraitsFromCardAffect,
		},
		SackCardEvent: {
			Type:   reflect.TypeOf(&SackCardArgs{}),
			Affect: SackCardAffect,
		},
		SummonUnitEvent: {
			Type:   reflect.TypeOf(&SummonUnitArgs{}),
			Affect: SummonUnitAffect,
		},
		SwapStatsEvent: {
			Type:   reflect.TypeOf(&SwapStatsArgs{}),
			Affect: SwapStatsAffect,
		},
		SwapUnitsEvent: {
			Type:   reflect.TypeOf(&SwapUnitsArgs{}),
			Affect: SwapUnitsAffect,
		},
	}
}
