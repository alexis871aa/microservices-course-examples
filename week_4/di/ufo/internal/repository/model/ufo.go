package model

import "time"

type SightingInfo struct {
	ObservedAt      *time.Time `bson:"observed_at,omitempty"`
	Location        string     `bson:"location"`
	Description     string     `bson:"description"`
	Color           *string    `bson:"color,omitempty"`
	Sound           *bool      `bson:"sound,omitempty"`
	DurationSeconds *int32     `bson:"duration_seconds,omitempty"`
}

type SightingUpdateInfo struct {
	ObservedAt      *time.Time `bson:"observed_at,omitempty"`
	Location        *string    `bson:"location,omitempty"`
	Description     *string    `bson:"description,omitempty"`
	Color           *string    `bson:"color,omitempty"`
	Sound           *bool      `bson:"sound,omitempty"`
	DurationSeconds *int32     `bson:"duration_seconds,omitempty"`
}

type Sighting struct {
	Uuid      string       `bson:"_id"`
	Info      SightingInfo `bson:"info"`
	CreatedAt time.Time    `bson:"created_at"`
	UpdatedAt *time.Time   `bson:"updated_at,omitempty"`
	DeletedAt *time.Time   `bson:"deleted_at,omitempty"`
}
