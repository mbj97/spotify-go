package api

import (
	"fmt"
	"log"
	"spotify-go/structs"

	"golang.org/x/net/context"
)

func (s *SpotifyClient) GetCurrentlyPlayingSong(ctx context.Context) *structs.CurrentlyPlayingSong {
	currentlyPlaying, err := s.C.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		log.Fatalf("Unable to get currently playing track: %v", err)
		return nil
	}
	if currentlyPlaying.Item != nil {
		return &structs.CurrentlyPlayingSong{
			Name:     fmt.Sprintf("%s - %s", currentlyPlaying.Item.Name, currentlyPlaying.Item.Artists[0].Name),
			ID:       currentlyPlaying.Item.ID.String(),
			IsActive: currentlyPlaying.Playing,
		}
	}
	return nil
}
