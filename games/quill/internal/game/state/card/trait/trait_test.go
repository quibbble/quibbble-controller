package trait

import (
	"math/rand"
	"testing"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	cd "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/card"
	ch "github.com/quibbble/quibbble-controller/games/quill/internal/game/state/hook/choose"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_NewTrait(t *testing.T) {
	args := map[string]interface{}{
		"Amount": 1,
	}
	trait, err := NewTrait(uuid.Nil, nil, AssassinTrait, args)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, trait.GetArgs().(*AssassinArgs).Amount)
}

func Test_ModifyingTraitArgs(t *testing.T) {
	gen := uuid.NewGen(rand.New(rand.NewSource(0)))
	trait, err := NewTrait(uuid.Nil, nil, FriendsTrait, &FriendsArgs{
		ChooseUnits: struct {
			Type string
			Args interface{}
		}{
			Type: ch.UUIDChoice,
			Args: ch.UUIDArgs{
				UUID: gen.New(en.UnitUUID),
			},
		},
		Trait: struct {
			Type string
			Args interface{}
		}{
			Type: BuffTrait,
			Args: BuffArgs{
				Stat:   cd.AttackStat,
				Amount: 1,
			},
		},
		Current: make([]uuid.UUID, 0),
	})
	if err != nil {
		t.Fatal(err)
	}

	args := trait.GetArgs().(*FriendsArgs)
	assert.True(t, len(trait.GetArgs().(*FriendsArgs).Current) == 0)
	args.Current = []uuid.UUID{gen.New(en.UnitUUID), gen.New(en.UnitUUID), gen.New(en.UnitUUID)}
	assert.True(t, len(trait.GetArgs().(*FriendsArgs).Current) > 0)
}
