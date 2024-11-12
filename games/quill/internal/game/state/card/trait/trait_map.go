package trait

import (
	"reflect"

	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
)

var TraitMap map[string]struct {
	Type   reflect.Type
	Add    func(t *Trait, card st.ICard) error
	Remove func(t *Trait, card st.ICard) error
}

func init() {
	TraitMap = map[string]struct {
		Type   reflect.Type
		Add    func(t *Trait, card st.ICard) error
		Remove func(t *Trait, card st.ICard) error
	}{
		AssassinTrait: {
			Type: reflect.TypeOf(&AssassinArgs{}),
		},
		BattleCryTrait: {
			Type: reflect.TypeOf(&BattleCryArgs{}),
		},
		BerserkTrait: {
			Type: reflect.TypeOf(&BerserkArgs{}),
		},
		BuffTrait: {
			Type:   reflect.TypeOf(&BuffArgs{}),
			Add:    AddBuff,
			Remove: RemoveBuff,
		},
		AimlessTrait: {
			Type: reflect.TypeOf(&AimlessArgs{}),
		},
		DeathCryTrait: {
			Type: reflect.TypeOf(&DeathCryArgs{}),
		},
		DebuffTrait: {
			Type:   reflect.TypeOf(&DebuffArgs{}),
			Add:    AddDebuff,
			Remove: RemoveDebuff,
		},
		DodgeTrait: {
			Type: reflect.TypeOf(&DodgeArgs{}),
		},
		EnemiesTrait: {
			Type: reflect.TypeOf(&EnemiesArgs{}),
		},
		EnrageTrait: {
			Type: reflect.TypeOf(&EnrageArgs{}),
		},
		EternalTrait: {
			Type: reflect.TypeOf(&EternalArgs{}),
		},
		ExecuteTrait: {
			Type: reflect.TypeOf(&ExecuteArgs{}),
		},
		FriendsTrait: {
			Type: reflect.TypeOf(&FriendsArgs{}),
		},
		GiftTrait: {
			Type: reflect.TypeOf(&GiftArgs{}),
		},
		HasteTrait: {
			Type: reflect.TypeOf(&HasteArgs{}),
		},
		LobberTrait: {
			Type: reflect.TypeOf(&LobberArgs{}),
		},
		PillageTrait: {
			Type: reflect.TypeOf(&PillageArgs{}),
		},
		PoisonTrait: {
			Type: reflect.TypeOf(&PoisonArgs{}),
		},
		PurityTrait: {
			Type: reflect.TypeOf(&PurityArgs{}),
		},
		RangedTrait: {
			Type: reflect.TypeOf(&RangedArgs{}),
		},
		RecodeTrait: {
			Type:   reflect.TypeOf(&RecodeArgs{}),
			Add:    AddRecode,
			Remove: RemoveRecode,
		},
		ShieldTrait: {
			Type: reflect.TypeOf(&ShieldArgs{}),
		},
		SpikyTrait: {
			Type: reflect.TypeOf(&SpikyArgs{}),
		},
		SurgeTrait: {
			Type: reflect.TypeOf(&SurgeArgs{}),
		},
		ThiefTrait: {
			Type: reflect.TypeOf(&ThiefArgs{}),
		},
		TiredTrait: {
			Type: reflect.TypeOf(&TiredArgs{}),
		},
		WardTrait: {
			Type: reflect.TypeOf(&WardArgs{}),
		},
	}
}
