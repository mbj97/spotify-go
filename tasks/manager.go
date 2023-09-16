package tasks

import (
	"context"
	"fmt"
	"spotify-go/api"
	"spotify-go/utils"
)

type TaskInteractor struct {
	spotifyClient *api.SpotifyClient
	taskScheduler *utils.TaskScheduler
}

func NewTaskInteractor(spotifyClient *api.SpotifyClient, taskScheduler *utils.TaskScheduler) *TaskInteractor {
	ti := TaskInteractor{spotifyClient: spotifyClient, taskScheduler: taskScheduler}
	return &ti
}

func (ti *TaskInteractor) StartCurrentlyListeningWatcher(ctx context.Context) {
	wrappedFunc := func() {
		song := ti.spotifyClient.GetCurrentlyPlayingSong(ctx)
		if song != nil && song.IsActive {
			utils.Log(fmt.Sprintf(" '%s' playing", song.Name))
		} else if song != nil && !song.IsActive {
			utils.Log(fmt.Sprintf(" '%s' paused", song.Name))
		} else {
			utils.Log("User inactive")
		}
	}
	ti.taskScheduler.CreateTask("listener", wrappedFunc, 5)
}
