package main

import (
	"github.com/deminzhang/qimen-go/world"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"net/http"
)

func main() {
	game := world.NewWorld()
	go func() {
		if err := http.ListenAndServe(":6061", nil); err != nil {
			log.Fatal(err)
		}
	}()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
