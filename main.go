package main

import (
	"github.com/deminzhang/qimen-go/world"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	game := world.NewWorld()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
