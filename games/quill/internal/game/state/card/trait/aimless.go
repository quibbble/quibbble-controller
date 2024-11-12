package trait

import "math/rand"

const (
	AimlessTrait = "Aimless"
)

type AimlessArgs struct{}

func BuildAimlessCodex(r *rand.Rand) string {
	codex := ""
	for i := 0; i < 8; i++ {
		if r.Intn(4) == 0 {
			codex += "1"
		} else {
			codex += "0"
		}
	}
	return codex
}
