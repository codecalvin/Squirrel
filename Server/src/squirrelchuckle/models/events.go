package models

import (
	"container/list"
)

type EventType int

const (
	EVENT_JOIN EventType = iota
	EVENT_LEAVE
	EVENT_MESSAGE
)

type Event struct {
	EventType
	User string
	Content string
	Time int
}

const archiveSize = 20

// Event archives.
var archive = list.New()

// NewArchive saves new event to archive list.
func NewArchive(event Event) {
	if archive.Len() >= archiveSize {
		archive.Remove(archive.Front())
	}
	archive.PushBack(event)
}