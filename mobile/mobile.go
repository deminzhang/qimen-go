package mobile

import (
	"github.com/hajimehoshi/ebiten/v2/mobile"

	"qimen/world"
)

func init() {
	game := world.NewWorld()
	mobile.SetGame(game)
}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}
