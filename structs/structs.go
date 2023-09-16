package structs

import "fmt"

type Status int

const (
	StatusStarted Status = iota
	StartedAuth
)

type CurrentlyPlayingSong struct {
	Name     string
	ID       string
	IsActive bool
}

func (cs *CurrentlyPlayingSong) GetSongStatusString() string {
	if cs != nil && cs.IsActive {
		return fmt.Sprintf(" '%s' playing", cs.Name)
	} else if cs != nil && !cs.IsActive {
		return fmt.Sprintf(" '%s' paused", cs.Name)
	} 
	return "User inactive"
}
