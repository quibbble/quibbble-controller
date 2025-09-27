package controller

import (
	"regexp"

	"github.com/pkg/errors"
	"github.com/quibbble/quibbble-controller/pkg/proto/controller"
)

var imageRegex = regexp.MustCompile(`^([a-z0-9]+(?:[._-][a-z0-9]+)*\/)?([a-z0-9]+(?:[._-][a-z0-9]+)*\/)*([a-z0-9]+(?:[._-][a-z0-9]+)*)(?::([a-zA-Z0-9._-]+))?$`)

func validate(gk *controller.GameKey) error {
	if !imageRegex.MatchString(gk.Repository + ":" + gk.Tag) {
		return errors.New("provided image is invalid")
	}
	return nil
}
