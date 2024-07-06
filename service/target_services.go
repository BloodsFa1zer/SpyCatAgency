package service

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"spyCat/database"
	"spyCat/database/models"
)

type TargetServiceInterface interface {
	UpdateTarget(target *models.Target) (*models.Target, error, int)
	UpdateTargetNotes(targetID int, notes string) (*models.Target, error, int)
	CompleteTarget(missionID, targetID int) (error, int)
	DeleteTarget(missionID, targetID int) (error, int)
	AddTarget(missionID int, target *models.Target) (*models.Target, error, int)
}

type TargetService struct {
	DbTarget database.TargetDatabaseInterface
	validate *validator.Validate
}

func NewTargetService(DbTarget database.TargetDatabaseInterface, validate *validator.Validate) *TargetService {
	return &TargetService{DbTarget: DbTarget, validate: validate}
}

func (ts *TargetService) UpdateTarget(target *models.Target) (*models.Target, error, int) {
	missionCompleted, err := ts.DbTarget.IsMissionCompleted(target.MissionID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if missionCompleted {
		return nil, errors.New("cannot update target of a completed mission"), http.StatusBadRequest
	}

	if target.Status == "completed" {
		return nil, errors.New("cannot update a completed target"), http.StatusBadRequest
	}

	err = ts.DbTarget.UpdateTarget(target)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return target, nil, http.StatusOK
}

func (ts *TargetService) UpdateTargetNotes(targetID int, notes string) (*models.Target, error, int) {
	target, err := ts.DbTarget.GetTarget(targetID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	isLinked, err := ts.DbTarget.IsTargetLinkedToMission(target.MissionID, target.ID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if !isLinked {
		return nil, errors.New("the specified target is not linked to this mission"), http.StatusBadRequest
	}

	missionCompleted, err := ts.DbTarget.IsMissionCompleted(target.MissionID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if missionCompleted {
		return nil, errors.New("cannot update notes of a target in a completed mission"), http.StatusBadRequest
	}

	if target.Status == "completed" {
		return nil, errors.New("cannot update notes of a completed target"), http.StatusConflict
	}

	err = ts.DbTarget.UpdateTargetNotes(targetID, notes)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	updatedTarget, err := ts.DbTarget.GetTarget(targetID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return updatedTarget, nil, http.StatusOK
}

func (ts *TargetService) CompleteTarget(missionID, targetID int) (error, int) {
	isLinked, err := ts.DbTarget.IsTargetLinkedToMission(missionID, targetID)
	if err != nil {
		return err, http.StatusBadRequest
	}
	if !isLinked {
		return errors.New("the specified target is not linked to this mission"), http.StatusBadRequest
	}

	missionCompleted, err := ts.DbTarget.IsMissionCompleted(missionID)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if missionCompleted {
		return errors.New("cannot complete a target of an already completed mission"), http.StatusConflict
	}

	return ts.DbTarget.CompleteTarget(missionID, targetID), http.StatusOK
}

func (ts *TargetService) DeleteTarget(missionID, targetID int) (error, int) {
	targetCompleted, err := ts.DbTarget.IsTargetCompleted(targetID)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if targetCompleted {
		return errors.New("cannot delete a completed target"), http.StatusBadRequest
	}

	isLinked, err := ts.DbTarget.IsTargetLinkedToMission(missionID, targetID)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if !isLinked {
		return errors.New("the specified target is not linked to this mission"), http.StatusBadRequest
	}

	mission, err := ts.DbTarget.GetMission(missionID)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	if len(mission.Targets) <= 1 {
		return errors.New("cannot delete the last target of a mission"), http.StatusConflict
	}

	return ts.DbTarget.DeleteTarget(missionID, targetID), http.StatusOK
}

func (ts *TargetService) AddTarget(missionID int, target *models.Target) (*models.Target, error, int) {
	missionCompleted, err := ts.DbTarget.IsMissionCompleted(missionID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if missionCompleted {
		return nil, errors.New("cannot add target to a completed mission"), http.StatusBadRequest
	}

	mission, err := ts.DbTarget.GetMission(missionID)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	if len(mission.Targets) >= 3 {
		return nil, errors.New("a mission cannot have more than 3 targets"), http.StatusConflict
	}

	target.MissionID = missionID
	target.Status = "in_progress"

	err = ts.DbTarget.AddTarget(target)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return target, nil, http.StatusOK
}
