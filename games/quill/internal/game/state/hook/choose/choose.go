package choose

import (
	"context"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

type Choose struct {
	uuid uuid.UUID

	Type string
	Args interface{}

	retrieve func(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error)
}

func NewChoose(uuid uuid.UUID, typ string, args interface{}) (en.IChoose, error) {
	c, ok := ChooseMap[typ]
	if !ok {
		return nil, errors.ErrMissingMapKey
	}
	decoded := reflect.New(c.Type).Elem().Interface()
	if err := mapstructure.Decode(args, &decoded); err != nil {
		return nil, errors.Wrap(err)
	}
	return &Choose{
		uuid:     uuid,
		Type:     typ,
		Args:     decoded,
		retrieve: c.Retrieve,
	}, nil
}

func (c *Choose) GetType() string {
	return c.Type
}

func (c *Choose) GetArgs() interface{} {
	return c.Args
}

func (c *Choose) Retrieve(ctx context.Context, engine en.IEngine, state en.IState) ([]uuid.UUID, error) {
	eng, ok := engine.(*en.Engine)
	if !ok {
		return nil, errors.ErrInterfaceConversion
	}
	sta, ok := state.(*st.State)
	if !ok {
		return nil, errors.ErrInterfaceConversion
	}
	return c.retrieve(c, ctx, eng, sta)
}
