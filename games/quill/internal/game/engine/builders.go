package engine

import (
	"github.com/quibbble/quibbble-controller/games/quill/parse"
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

type Builders struct {
	BuildCondition
	BuildEvent
	BuildHook
	BuildChoose
}

func (b *Builders) NewHook(gen *uuid.Gen, card uuid.UUID, hook parse.Hook) (IHook, error) {
	conditions := make([]ICondition, 0)
	for _, c := range hook.Conditions {
		condition, err := b.BuildCondition(gen.New(ConditionUUID), c.Type, c.Not, c.Args)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		conditions = append(conditions, condition)
	}
	events := make([]IEvent, 0)
	for _, e := range hook.Events {
		event, err := b.BuildEvent(gen.New(EventUUID), e.Type, e.Args)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		events = append(events, event)
	}
	reuse := make([]ICondition, 0)
	for _, c := range hook.ReuseConditions {
		re, err := b.BuildCondition(gen.New(ConditionUUID), c.Type, c.Not, c.Args)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		reuse = append(reuse, re)
	}
	return b.BuildHook(gen.New(HookUUID), card, hook.When, hook.Types, conditions, events, reuse)
}
