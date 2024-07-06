package models

type MissionStatus string

const (
	MissionStatusInProgress MissionStatus = "in_progress"
	MissionStatusCompleted  MissionStatus = "completed"
)

type Mission struct {
	ID        int           `db:"id" json:"ID"`
	CatID     int           `db:"cat_id" json:"CatID" validate:"required"`
	Status    MissionStatus `db:"status" json:"Status" validate:"required"`
	Targets   []Target      `json:"Targets" validate:"required"`
	CreatedAt string        `db:"created_at" json:"CreatedAt"`
	UpdatedAt string        `db:"updated_at" json:"UpdatedAt,omitempty"`
}
