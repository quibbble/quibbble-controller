package hook

import (
	"context"
	"slices"

	en "github.com/quibbble/quibbble-controller/games/quill/internal/game/engine"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

type Hook struct {
	uuid       uuid.UUID
	cardUUID   uuid.UUID // the card that registered the hook
	when       en.When
	types      []string
	conditions en.Conditions
	events     []en.IEvent
	reuse      en.Conditions
}

func NewHook(uuid, cardUUID uuid.UUID, when string, types []string, conditions []en.ICondition, events []en.IEvent, reuse []en.ICondition) (en.IHook, error) {
	return &Hook{
		uuid:       uuid,
		cardUUID:   cardUUID,
		when:       en.When(when),
		types:      types,
		conditions: conditions,
		events:     events,
		reuse:      reuse,
	}, nil
}

func (h *Hook) GetUUID() uuid.UUID {
	return h.uuid
}

func (h *Hook) GetCardUUID() uuid.UUID {
	return h.cardUUID
}

func (h *Hook) GetTypes() []string {
	return h.types
}

func (h *Hook) Trigger(when en.When, typ string) bool {
	return h.when == when && slices.Contains(h.types, typ)
}

func (h *Hook) Pass(ctx context.Context, engine en.IEngine, state en.IState) (bool, error) {
	return h.conditions.Pass(ctx, engine, state)
}

func (h *Hook) Events() []en.IEvent {
	return h.events
}

func (h *Hook) Reuse(ctx context.Context, engine en.IEngine, state en.IState) (bool, error) {
	return h.reuse.Pass(ctx, engine, state)
}
