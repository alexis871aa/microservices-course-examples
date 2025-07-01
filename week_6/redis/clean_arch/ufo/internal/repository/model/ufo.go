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

// SightingRedisView - модель для хранения в Redis hash map
type SightingRedisView struct {
	UUID         string  `redis:"uuid"`
	ObservedAtNs *int64  `redis:"observed_at,omitempty"`
	Location     string  `redis:"location"`
	Description  string  `redis:"description"`
	Color        *string `redis:"color,omitempty"`
	Sound        *bool   `redis:"sound,omitempty"`
	Duration     *int32  `redis:"duration_seconds,omitempty"`
	CreatedAtNs  int64   `redis:"created_at"`
	UpdatedAtNs  *int64  `redis:"updated_at,omitempty"`
	DeletedAtNs  *int64  `redis:"deleted_at,omitempty"`
}
