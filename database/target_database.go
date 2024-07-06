package database

import (
	"database/sql"
	"errors"
	"spyCat/database/models"
	"time"
)

type TargetDatabaseInterface interface {
	IsTargetCompleted(targetID int) (bool, error)
	GetTarget(targetID int) (*models.Target, error)
	UpdateTargetNotes(targetID int, notes string) error
	UpdateTarget(target *models.Target) error
	CompleteTarget(missionID, targetID int) error
	DeleteTarget(missionID, targetID int) error
	AddTarget(target *models.Target) error
	IsMissionCompleted(missionID int) (bool, error)
	GetMission(id int) (*models.Mission, error)
	IsTargetLinkedToMission(missionID, targetID int) (bool, error)
}

type TargetDatabase struct {
	*Database
}

func NewTargetDatabase(Conn *Database) *TargetDatabase {
	return &TargetDatabase{Conn}
}

func (td *TargetDatabase) UpdateTarget(target *models.Target) error {
	_, err := td.Connection.Exec(`
		UPDATE targets
		SET name = $1, country = $2, notes = $3, status = $4, updated_at = $5
		WHERE id = $6 AND mission_id = $7
	`, target.Name, target.Country, target.Notes, target.Status, time.Now(), target.ID, target.MissionID)
	return err
}

func (td *TargetDatabase) CompleteTarget(missionID, targetID int) error {
	_, err := td.Connection.Exec(`
		UPDATE targets
		SET status = 'completed', updated_at = $1
		WHERE id = $2 AND mission_id = $3
	`, time.Now(), targetID, missionID)
	return err
}

func (td *TargetDatabase) DeleteTarget(missionID, targetID int) error {
	result, err := td.Connection.Exec("DELETE FROM targets WHERE id = $1 AND mission_id = $2", targetID, missionID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return err
}

func (td *TargetDatabase) AddTarget(target *models.Target) error {
	return td.Connection.QueryRow(`
		INSERT INTO targets (mission_id, name, country, notes, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, target.MissionID, target.Name, target.Country, target.Notes, target.Status, time.Now(), time.Now()).Scan(&target.ID)
}

func (td *TargetDatabase) GetTarget(targetID int) (*models.Target, error) {
	var target models.Target
	var createdAt, updatedAt time.Time

	err := td.Connection.QueryRow(`
        SELECT id, mission_id, name, country, notes, status, created_at, updated_at
        FROM targets
        WHERE id = $1
    `, targetID).Scan(&target.ID, &target.MissionID, &target.Name, &target.Country, &target.Notes, &target.Status, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	target.CreatedAt = createdAt.Format("15:04:05 02:01:06")
	target.UpdatedAt = updatedAt.Format("15:04:05 02:01:06")

	return &target, nil
}

func (td *TargetDatabase) IsTargetCompleted(targetID int) (bool, error) {
	var status string
	err := td.Connection.QueryRow("SELECT status FROM targets WHERE id = $1", targetID).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return status == "completed", nil
}

func (td *TargetDatabase) UpdateTargetNotes(targetID int, notes string) error {
	_, err := td.Connection.Exec(`
        UPDATE targets
        SET notes = $1, updated_at = $2
        WHERE id = $3
    `, notes, time.Now(), targetID)
	return err
}

func (td *TargetDatabase) IsMissionCompleted(missionID int) (bool, error) {
	var status string
	err := td.Connection.QueryRow("SELECT status FROM missions WHERE id = $1", missionID).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return status == "completed", nil
}

func (td *TargetDatabase) GetMission(id int) (*models.Mission, error) {
	var mission models.Mission
	err := td.Connection.QueryRow(`
		SELECT id, cat_id, status, created_at, updated_at
		FROM missions
		WHERE id = $1
	`, id).Scan(&mission.ID, &mission.CatID, &mission.Status, &mission.CreatedAt, &mission.UpdatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := td.Connection.Query(`
		SELECT id, name, country, notes, status, created_at, updated_at
		FROM targets
		WHERE mission_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var target models.Target
		err := rows.Scan(&target.ID, &target.Name, &target.Country, &target.Notes, &target.Status, &target.CreatedAt, &target.UpdatedAt)
		if err != nil {
			return nil, err
		}
		target.MissionID = mission.ID
		mission.Targets = append(mission.Targets, target)
	}

	return &mission, nil
}

func (td *TargetDatabase) IsTargetLinkedToMission(missionID, targetID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM targets WHERE id = $1 AND mission_id = $2)`
	err := td.Connection.QueryRow(query, targetID, missionID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
