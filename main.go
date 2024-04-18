package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"qimen/world"
)

func main() {
	game := world.NewWorld()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
