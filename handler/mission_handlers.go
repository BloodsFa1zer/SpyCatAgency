package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"spyCat/database/models"
	"spyCat/response"
	"spyCat/service"
	"strconv"
)

type MissionHandler struct {
	MissionService service.MissionServiceInterface
}

func NewMissionHandler(service service.MissionServiceInterface) *MissionHandler {
	return &MissionHandler{MissionService: service}
}

type MissionHandlerInterface interface {
	CreateMission(c echo.Context) error
	DeleteMission(c echo.Context) error
	UpdateMission(c echo.Context) error
	CompleteMission(c echo.Context) error
	AssignCatToMission(c echo.Context) error
	ListMissions(c echo.Context) error
	GetMission(c echo.Context) error
}

// CreateMission creates a new mission
func (mh *MissionHandler) CreateMission(c echo.Context) error {
	mission := new(models.Mission)
	if err := c.Bind(mission); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	createdMission, err, respStatus := mh.MissionService.CreateMission(mission)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, response.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": createdMission}})
}

// DeleteMission deletes a mission by ID
func (mh *MissionHandler) DeleteMission(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	err = mh.MissionService.DeleteMission(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})

	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "Mission successfully deleted"}})
}

// UpdateMission updates a mission
func (mh *MissionHandler) UpdateMission(c echo.Context) error {
	mission := new(models.Mission)
	if err := c.Bind(mission); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	updatedMission, err, respStatus := mh.MissionService.UpdateMission(mission)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedMission}})
}

// CompleteMission marks a mission as completed
func (mh *MissionHandler) CompleteMission(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	err = mh.MissionService.CompleteMission(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	resp := fmt.Sprintf("Mission %d marked as completed", id)

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": resp}})
}

// AssignCatToMission assigns a cat to a mission
func (mh *MissionHandler) AssignCatToMission(c echo.Context) error {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "Invalid mission ID"}})
	}

	var requestBody struct {
		CatID int `json:"CatID"`
	}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "Invalid mission ID"}})
	}

	err, respStatus := mh.MissionService.AssignCatToMission(missionID, requestBody.CatID)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": "Invalid mission ID"}})
	}

	resp := fmt.Sprintf("Mission %d assigned to Cat %d", missionID, requestBody.CatID)
	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": resp}})
}

// ListMissions retrieves all missions
func (mh *MissionHandler) ListMissions(c echo.Context) error {
	missions, err := mh.MissionService.ListMissions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": "Invalid mission ID"}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": missions}})
}

// GetMission retrieves a specific mission by ID
func (mh *MissionHandler) GetMission(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "Invalid mission ID"}})
	}

	mission, err := mh.MissionService.GetMission(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": "Invalid mission ID"}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": mission}})
}
