package models

type TargetStatus string

const (
	TargetStatusInProgress TargetStatus = "in_progress"
	TargetStatusCompleted  TargetStatus = "completed"
)

type Target struct {
	ID        int          `db:"id" json:"ID"`
	MissionID int          `db:"mission_id" json:"MissionID"`
	Name      string       `db:"name" json:"Name" validate:"required"`
	Country   string       `db:"country" json:"Country" validate:"required"`
	Notes     string       `db:"notes" json:"Notes"`
	Status    TargetStatus `db:"salary" json:"Status" validate:"required"`
	CreatedAt string       `db:"created_at" json:"CreatedAt"`
	UpdatedAt string       `db:"updated_at" json:"UpdatedAt,omitempty"`
}
