package main

import (
	"fmt"
	"log"
	"os"

	"github.com/deminzhang/qimen-go/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joho/godotenv"
)

func loadEnv() error {
	return godotenv.Load("config.env")
}

func getConfig() (url, model, apiKey string) {
	url = os.Getenv("AI_API_URL")
	model = os.Getenv("AI_MODEL")
	apiKey = os.Getenv("AI_API_KEY")
	return
}
func main() {
	if err := loadEnv(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	url, model, apiKey := getConfig()
	fmt.Println(url, model, apiKey)

	game := world.NewWorld()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
