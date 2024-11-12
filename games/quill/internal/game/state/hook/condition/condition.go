package condition

import (
	"context"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	st "github.com/quibbble/quibbble-controller/games/quill/internal/game/state"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

type Condition struct {
	uuid uuid.UUID

	Type string
	Not  bool
	Args interface{}

	pass func(c *Condition, ctx context.Context, engine *en.Engine, state *st.State) (bool, error)
}

func NewCondition(uuid uuid.UUID, typ string, not bool, args interface{}) (en.ICondition, error) {
	p, ok := ConditionMap[typ]
	if !ok {
		return nil, errors.ErrMissingMapKey
	}
	decoded := reflect.New(p.Type).Elem().Interface()
	if err := mapstructure.Decode(args, &decoded); err != nil {
		return nil, errors.Wrap(err)
	}
	return &Condition{
		uuid: uuid,
		Type: typ,
		Not:  not,
		Args: decoded,
		pass: p.Pass,
	}, nil
}

func (c *Condition) GetUUID() uuid.UUID {
	return c.uuid
}

func (c *Condition) GetType() string {
	return c.Type
}

func (c *Condition) GetNot() bool {
	return c.Not
}

func (c *Condition) GetArgs() interface{} {
	return c.Args
}

func (c *Condition) Pass(ctx context.Context, engine en.IEngine, state en.IState) (bool, error) {
	eng, ok := engine.(*en.Engine)
	if !ok {
		return false, errors.ErrInterfaceConversion
	}
	sta, ok := state.(*st.State)
	if !ok {
		return false, errors.ErrInterfaceConversion
	}
	pass, err := c.pass(c, ctx, eng, sta)
	if err != nil {
		return false, errors.Wrap(err)
	}
	if c.Not {
		return !pass, nil
	}
	return pass, nil
}

func SliceToConditions(conditions []*Condition) en.Conditions {
	c := make([]en.ICondition, 0)
	for _, condition := range conditions {
		c = append(c, condition)
	}
	return en.Conditions(c)
}
