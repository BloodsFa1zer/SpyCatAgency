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

type TargetHandler struct {
	TargetService service.TargetServiceInterface
}

func NewTargetHandler(service service.TargetServiceInterface) *TargetHandler {
	return &TargetHandler{TargetService: service}
}

type TargetHandlerInterface interface {
	UpdateTarget(c echo.Context) error
	UpdateTargetNotes(c echo.Context) error
	CompleteTarget(c echo.Context) error
	DeleteTarget(c echo.Context) error
	AddTarget(c echo.Context) error
}

// UpdateTarget updates a target
func (th *TargetHandler) UpdateTarget(c echo.Context) error {
	target := new(models.Target)
	if err := c.Bind(target); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	updatedTarget, err, respStatus := th.TargetService.UpdateTarget(target)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedTarget}})
}

// UpdateTargetNotes updates the notes of a target
func (th *TargetHandler) UpdateTargetNotes(c echo.Context) error {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	var requestBody struct {
		Notes string `json:"Notes"`
	}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	_, err, respStatus := th.TargetService.UpdateTargetNotes(targetID, requestBody.Notes)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	resp := fmt.Sprintf("Notes on target %d updated", targetID)
	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": resp}})
}

// CompleteTarget marks a target as completed
func (th *TargetHandler) CompleteTarget(c echo.Context) error {
	missionID, err := strconv.Atoi(c.Param("missionId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	targetID, err := strconv.Atoi(c.Param("targetId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	err, respStatus := th.TargetService.CompleteTarget(missionID, targetID)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	resp := fmt.Sprintf("Target %d on Mission %d marked as completed", targetID, missionID)
	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": resp}})
}

// DeleteTarget deletes a target from a Target
func (th *TargetHandler) DeleteTarget(c echo.Context) error {
	missionID, err := strconv.Atoi(c.Param("missionId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	targetID, err := strconv.Atoi(c.Param("targetId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	err, respStatus := th.TargetService.DeleteTarget(missionID, targetID)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	resp := fmt.Sprintf("Target %d on Mission %d deleted", targetID, missionID)
	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": resp}})
}

// AddTarget adds a new target to a Target
func (th *TargetHandler) AddTarget(c echo.Context) error {
	TargetID, err := strconv.Atoi(c.Param("missionId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "Invalid Target ID"}})
	}

	var target models.Target
	if err := c.Bind(&target); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	createdTarget, err, respStatus := th.TargetService.AddTarget(TargetID, &target)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, response.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": createdTarget}})
}
