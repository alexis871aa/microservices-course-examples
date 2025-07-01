package model

import "time"

type UFORecordedEvent struct {
	UUID        string
	ObservedAt  *time.Time
	Location    string
	Description string
}
