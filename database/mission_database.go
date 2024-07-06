package database

import (
	"database/sql"
	"errors"
	"spyCat/database/models"
	"time"
)

type MissionDatabaseInterface interface {
	CreateMission(mission *models.Mission) error
	DeleteMission(id int) error
	UpdateMission(mission *models.Mission) error
	CompleteMission(id int) error
	AssignCatToMission(missionID, catID int) error
	ListMissions() (*[]models.Mission, error)
	GetMission(id int) (*models.Mission, error)
	IsMissionAssignedToCat(missionID int) (bool, error)
	IsMissionCompleted(missionID int) (bool, error)
	IsCatAvailable(catID int) (bool, error)
	IsMissionAssigned(missionID int) (bool, error)
	DoesCatExist(catID int) (bool, error)
}

type MissionDatabase struct {
	*Database
}

func NewMissionDatabase(Conn *Database) *MissionDatabase {
	return &MissionDatabase{Conn}
}

func (md *MissionDatabase) CreateMission(mission *models.Mission) error {
	tx, err := md.Connection.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRow(`
		INSERT INTO missions (cat_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, mission.CatID, mission.Status, time.Now(), time.Now()).Scan(&mission.ID)
	if err != nil {
		return err
	}

	for i := range mission.Targets {
		err = tx.QueryRow(`
			INSERT INTO targets (mission_id, name, country, notes, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id
		`, mission.ID, mission.Targets[i].Name, mission.Targets[i].Country, mission.Targets[i].Notes,
			mission.Targets[i].Status, time.Now(), time.Now()).Scan(&mission.Targets[i].ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (md *MissionDatabase) DeleteMission(id int) error {
	result, err := md.Connection.Exec("DELETE FROM missions WHERE id = $1", id)
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

func (md *MissionDatabase) UpdateMission(mission *models.Mission) error {
	result, err := md.Connection.Exec(`UPDATE missions
	SET cat_id = $1, status = $2, updated_at = $3
	WHERE id = $4`, mission.CatID, mission.Status, time.Now(), mission.ID)
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

func (md *MissionDatabase) CompleteMission(id int) error {
	_, err := md.Connection.Exec(`
		UPDATE missions
		SET status = 'completed', updated_at = $1
		WHERE id = $2
	`, time.Now(), id)
	return err
}

func (md *MissionDatabase) ListMissions() (*[]models.Mission, error) {
	rows, err := md.Connection.Query(`
		SELECT id, cat_id, status, created_at, updated_at
		FROM missions
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var createdAt, updatedAt time.Time
	var missions []models.Mission
	for rows.Next() {
		var mission models.Mission
		err := rows.Scan(&mission.ID, &mission.CatID, &mission.Status, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		mission.CreatedAt = createdAt.Format("15:04:05 02:01:06")
		mission.UpdatedAt = updatedAt.Format("15:04:05 02:01:06")

		missions = append(missions, mission)
	}

	return &missions, nil
}

func (md *MissionDatabase) GetMission(id int) (*models.Mission, error) {
	var mission models.Mission
	var createdAt, updatedAt time.Time
	err := md.Connection.QueryRow(`
		SELECT id, cat_id, status, created_at, updated_at
		FROM missions
		WHERE id = $1
	`, id).Scan(&mission.ID, &mission.CatID, &mission.Status, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := md.Connection.Query(`
		SELECT id, name, country, notes, status, created_at, updated_at
		FROM targets
		WHERE mission_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mission.CreatedAt = createdAt.Format("15:04:05 02:01:06")
	mission.UpdatedAt = updatedAt.Format("15:04:05 02:01:06")

	for rows.Next() {
		var target models.Target
		err := rows.Scan(&target.ID, &target.Name, &target.Country, &target.Notes, &target.Status, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		target.MissionID = mission.ID
		target.CreatedAt = createdAt.Format("15:04:05 02:01:06")
		target.UpdatedAt = updatedAt.Format("15:04:05 02:01:06")
		mission.Targets = append(mission.Targets, target)
	}

	return &mission, nil
}

func (md *MissionDatabase) IsMissionAssignedToCat(missionID int) (bool, error) {
	var catID sql.NullInt64
	err := md.Connection.QueryRow("SELECT cat_id FROM missions WHERE id = $1", missionID).Scan(&catID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return catID.Valid, nil
}

func (md *MissionDatabase) IsMissionCompleted(missionID int) (bool, error) {
	var status string
	err := md.Connection.QueryRow("SELECT status FROM missions WHERE id = $1", missionID).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return status == "completed", nil
}

func (md *MissionDatabase) AssignCatToMission(missionID, catID int) error {
	_, err := md.Connection.Exec(`
		UPDATE missions
		SET cat_id = $1, updated_at = $2
		WHERE id = $3
	`, catID, time.Now(), missionID)
	return err
}

func (md *MissionDatabase) IsCatAvailable(catID int) (bool, error) {
	var count int
	err := md.Connection.QueryRow("SELECT COUNT(*) FROM missions WHERE cat_id = $1 AND status = 'in_progress'", catID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (md *MissionDatabase) IsMissionAssigned(missionID int) (bool, error) {
	var catID sql.NullInt64
	err := md.Connection.QueryRow("SELECT cat_id FROM missions WHERE id = $1", missionID).Scan(&catID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // Mission doesn't exist
		}
		return false, err
	}
	return catID.Valid, nil // If catID is valid (not NULL), the mission is assigned
}

func (md *MissionDatabase) DoesCatExist(catID int) (bool, error) {
	var exists bool
	err := md.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM spy_cats WHERE id = $1)", catID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
