package main

import (
	"context"
	"log"
	"spotify-go/api"
	"spotify-go/utils"
	"spotify-go/tasks"
)

func main() {
	ctx := context.Background()
	scheduler := utils.NewTaskScheduler()
	sClient, err := api.NewSpotifyClient()
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	taskInteractor := tasks.NewTaskInteractor(sClient,scheduler)
	taskInteractor.StartCurrentlyListeningWatcher(ctx)

	scheduler.Start()

	select {}
}
