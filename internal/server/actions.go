package server

import qg "github.com/quibbble/quibbble-controller/pkg/game"

// Valid Game Server Actions
const (
	Join = "join"
	Chat = "chat"
)

type Action struct {
	*qg.Action
	*Player
}
