package state

import (
	"github.com/quibbble/quibbble-controller/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/uuid"
)

var (
	ErrNotFound = func(uuid uuid.UUID) error { return errors.Errorf("%s not found", uuid) }
)
