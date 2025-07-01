package model

import "time"

type SightingInfo struct {
	ObservedAt      *time.Time
	Location        string
	Description     string
	Color           *string
	Sound           *bool
	DurationSeconds *int32
}

type SightingUpdateInfo struct {
	ObservedAt      *time.Time
	Location        *string
	Description     *string
	Color           *string
	Sound           *bool
	DurationSeconds *int32
}

type Sighting struct {
	Uuid      string
	Info      SightingInfo
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
