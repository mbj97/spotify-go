package handler

import (
	"context"
	"fmt"
	"net/http"
	"spotify-go/api"
	"spotify-go/tasks"
	"spotify-go/utils"
	"sync"
)

type Status int

const (
	STATUS_NEW Status = iota
	STATUS_READY
)

type UserTaskGroup struct {
	id             string
	status         Status
	taskScheduler  *utils.TaskScheduler
	taskInteractor *tasks.TaskInteractor
	spotifyClient  *api.SpotifyClient
}

var (
	userTaskGroups = make(map[string]*UserTaskGroup)
	mu             sync.Mutex
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	username := params.Get("username")

	if username == "" {
		http.Error(w, "Missing username parameter", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, ok := userTaskGroups[username]; !ok {
		spotifyClient, err := api.CreateNewClient(username)
		fmt.Println("NEW CLIENT")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		userTaskGroups[username] = &UserTaskGroup{
			status:         STATUS_NEW,
			id:             username,
			taskScheduler:  utils.NewTaskScheduler(),
			taskInteractor: nil,
			spotifyClient:  spotifyClient,
		}
	}
	utg := userTaskGroups[username]
	if utg.spotifyClient.Status == api.STATUS_AUTHORIZED {
		fmt.Fprintf(w, "Authenticated successfully")
		if utg.status != STATUS_READY {
			utg.taskInteractor = tasks.NewTaskInteractor(utg.spotifyClient, utg.taskScheduler)
			utg.status = STATUS_READY
		}
		return
	}
	fmt.Fprintf(w, "URL: %s", utg.spotifyClient.LoginURL)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	username := params.Get("state")
	if username == "" {
		http.Error(w, "Missing state parameter", http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := userTaskGroups[username]; ok {
		utg := userTaskGroups[username]
		utg.spotifyClient.CompleteAuth(w, r)
		if utg.spotifyClient.Status == api.STATUS_AUTHORIZED {
			if utg.status != STATUS_READY {
				utg.taskInteractor = tasks.NewTaskInteractor(utg.spotifyClient, utg.taskScheduler)
				utg.status = STATUS_READY
			}
		}
	}
}

func CurrentlyPlayingHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	username := params.Get("username")
	if username == "" {
		http.Error(w, "Missing username parameter", http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := userTaskGroups[username]; ok {
		utg := userTaskGroups[username]
		if utg.spotifyClient == nil || utg.spotifyClient.Status != api.STATUS_AUTHORIZED {
			fmt.Fprintf(w, "Not authenticated")
			return
		}
		songStatus := utg.spotifyClient.GetCurrentlyPlayingSong(context.Background()).GetSongStatusString()
		fmt.Fprintf(w, songStatus)
		return
	}
	fmt.Fprintf(w, "Not authenticated")
}

func StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	username := params.Get("username")
	if username == "" {
		http.Error(w, "Missing username parameter", http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := userTaskGroups[username]; ok {
		utg := userTaskGroups[username]
		if utg.spotifyClient == nil || utg.spotifyClient.Status != api.STATUS_AUTHORIZED {
			fmt.Fprintf(w, "Not authenticated")
			return
		}
		utg.taskInteractor.StartCurrentlyListeningWatcher(context.Background())
		utg.taskScheduler.Start()
		fmt.Fprint(w, "Started listener")
		return
	}
	fmt.Fprintf(w, "Not authenticated")
}
