package service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"net/http"
	"spyCat/database"
	"spyCat/database/models"
)

type MissionServiceInterface interface {
	CreateMission(mission *models.Mission) (*models.Mission, error, int)
	DeleteMission(id int) error
	UpdateMission(mission *models.Mission) (*models.Mission, error, int)
	CompleteMission(id int) error
	AssignCatToMission(missionID, catID int) (error, int)
	ListMissions() (*[]models.Mission, error)
	GetMission(id int) (*models.Mission, error)
}

type MissionService struct {
	DbMission database.MissionDatabaseInterface
	validate  *validator.Validate
}

func NewMissionService(DbMission database.MissionDatabaseInterface, validate *validator.Validate) *MissionService {
	return &MissionService{DbMission: DbMission, validate: validate}
}

func (ms *MissionService) CreateMission(mission *models.Mission) (*models.Mission, error, int) {
	isAvailable, err := ms.DbMission.IsCatAvailable(mission.CatID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if !isAvailable {
		return nil, errors.New("the selected cat is already assigned to an active mission"), http.StatusConflict
	}

	if len(mission.Targets) < 1 || len(mission.Targets) > 3 {
		return nil, errors.New("a mission must have between 1 and 3 targets"), http.StatusBadRequest
	}

	mission.Status = "in_progress"
	for i := range mission.Targets {
		mission.Targets[i].Status = "in_progress"
	}

	err = ms.DbMission.CreateMission(mission)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return mission, nil, http.StatusCreated
}

func (ms *MissionService) DeleteMission(id int) error {
	assigned, err := ms.DbMission.IsMissionAssignedToCat(id)
	if err != nil {
		return err
	}
	if assigned {
		return errors.New("cannot delete a mission assigned to a cat")
	}

	return ms.DbMission.DeleteMission(id)
}

func (ms *MissionService) UpdateMission(mission *models.Mission) (*models.Mission, error, int) {
	if mission.CatID != 0 {
		catExists, err := ms.DbMission.DoesCatExist(mission.CatID)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}
		if !catExists {
			return nil, errors.New("the specified cat does not exist"), http.StatusBadRequest
		}

		isAvailable, err := ms.DbMission.IsCatAvailable(mission.CatID)
		if err != nil {
			return nil, err, http.StatusBadRequest
		}
		if !isAvailable {
			return nil, errors.New("the selected cat is already assigned to an active mission"), http.StatusConflict
		}
	}

	err := ms.DbMission.UpdateMission(mission)
	if err != nil {
		// Check for foreign key violation
		var pqErr *pq.Error
		ok := errors.As(err, &pqErr)
		if ok {
			switch pqErr.Code {
			case "23503": // foreign_key_violation
				return nil, errors.New("invalid cat ID: the specified cat does not exist"), http.StatusBadRequest
			}
		}
		return nil, err, http.StatusInternalServerError
	}

	return mission, nil, http.StatusOK
}

func (ms *MissionService) CompleteMission(id int) error {
	mission, err := ms.DbMission.GetMission(id)
	if err != nil {
		return err
	}

	for _, target := range mission.Targets {
		if target.Status != "completed" {
			return errors.New("all targets must be completed before completing the mission")
		}
	}

	return ms.DbMission.CompleteMission(id)
}

func (ms *MissionService) ListMissions() (*[]models.Mission, error) {
	mission, err := ms.DbMission.ListMissions()
	if err != nil {
		return nil, err
	}

	return mission, err
}

func (ms *MissionService) GetMission(id int) (*models.Mission, error) {
	return ms.DbMission.GetMission(id)
}

func (ms *MissionService) AssignCatToMission(missionID, catID int) (error, int) {
	// Check if the cat is available
	isAvailable, err := ms.DbMission.IsCatAvailable(catID)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if !isAvailable {
		return errors.New("the selected cat is already assigned to an active mission"), http.StatusConflict
	}

	isMissionAssigned, err := ms.DbMission.IsMissionAssigned(missionID)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if isMissionAssigned {
		return errors.New("this mission is already assigned to a cat"), http.StatusBadRequest
	}

	return ms.DbMission.AssignCatToMission(missionID, catID), http.StatusOK
}
