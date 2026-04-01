package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/green-api/maxbot-api-client-go/pkg/client"
	"github.com/green-api/maxbot-api-client-go/pkg/models"
	"github.com/green-api/maxbot-chatbot-go/pkg/bot"
	n "github.com/green-api/maxbot-chatbot-go/pkg/notification"
	"github.com/green-api/maxbot-chatbot-go/pkg/state"

	s "github.com/green-api/maxbot-demo-chatbot-go/scenes"
)

type ExecutableScene interface {
	state.Scene
	Execute(n *n.Notification)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	appBot, err := bot.NewBot(client.Config{
		BaseURL:   os.Getenv("BASE_URL"),
		Token:     os.Getenv("TOKEN"),
		GlobalRPS: 25,
		Timeout:   35 * time.Second,
	})

	if err != nil {
		log.Fatal().Msgf("Bot initialization error: %v", err)
	}

	startScene := &s.StartScene{}
	appBot.StateManager = state.NewMapStateManager(map[string]any{})
	appBot.StateManager.SetStartScene(startScene)

	startTime := time.Now().Unix()

	sceneHandler := func(n *n.Notification) {
		if n.Update != nil && int64(n.Update.Timestamp) < startTime {
			return
		}
		n.CreateStateId()

		if appBot.StateManager.Get(n.StateId) == nil {
			appBot.StateManager.Create(n.StateId)
		}

		currentScene := n.GetCurrentScene()
		if currentScene == nil {
			currentScene = startScene
			n.ActivateNextScene(currentScene)
		}

		if execScene, ok := currentScene.(ExecutableScene); ok {
			execScene.Execute(n)
		} else {
			log.Error().Msg("Current scene does not implement ExecutableScene")
		}
	}

	appBot.Router.Register(models.TypeMessageCreated, sceneHandler)
	appBot.Router.Register(models.TypeMessageCallback, sceneHandler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go appBot.StartPolling(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info().Msg("The bot has been stopped")
}
