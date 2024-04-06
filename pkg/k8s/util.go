package k8s

import (
	"fmt"
	"strings"
)

const (
	Component = "component"

	ControllerComponennt = "controller"
	GameComponent        = "game"
)

const ChartName = "quibbble-controller"

const Namespace = "quibbble"

func Name(key, id string) string {
	return fmt.Sprintf("%s.%s", key, id)
}

func KeyID(name string) (string, string) {
	s := strings.Split(name, ".")
	return s[0], s[1]
}
