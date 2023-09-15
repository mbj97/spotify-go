package main

import (
	"context"
	"fmt"
	"log"
	"spotify-go/api"
)

func main() {
	ctx := context.Background()
	sClient, err := api.NewSpotifyClient()
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	currentlyPlaying := sClient.GetCurrentlyPlayingSong(ctx)
	if currentlyPlaying == nil {
		fmt.Println("User is currently not playing a song")
	} else {
		fmt.Println(currentlyPlaying.Name)
	}

}
