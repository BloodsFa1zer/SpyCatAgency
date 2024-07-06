package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"spyCat/database/models"
	"spyCat/response"
	"spyCat/service"
	"strconv"
)

type CatHandler struct {
	catService service.CatServiceInterface
}

func NewCatHandler(service service.CatServiceInterface) *CatHandler {
	return &CatHandler{catService: service}
}

type CatHandlerInterface interface {
	UpdateCatSalary(c echo.Context) error
	CreateCat(c echo.Context) error
	GetCat(c echo.Context) error
	GetAllCats(c echo.Context) error
	DeleteCat(c echo.Context) error
}

func (ch *CatHandler) CreateCat(c echo.Context) error {
	var cat models.Cat

	if err := c.Bind(&cat); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	catID, err, respStatus := ch.catService.CreateCat(cat)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, response.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"use this ID to interact with cat`s profile": catID}})
}

func (ch *CatHandler) GetCat(c echo.Context) error {
	ID := c.Param("id")
	catID, err := strconv.Atoi(ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	cat, err, respStatus := ch.catService.GetCat(catID)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": cat}})
}

func (ch *CatHandler) UpdateCatSalary(c echo.Context) error {
	ID := c.Param("id")
	catID, err := strconv.Atoi(ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	var salaryUpdate struct {
		Salary float64 `json:"salary"`
	}

	if err := c.Bind(&salaryUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": "Invalid request body"}})
	}

	err, respStatus := ch.catService.EditCatSalary(catID, salaryUpdate.Salary)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success",
		Data: &echo.Map{"data": "cat salary successfully updated"}})
}

func (ch *CatHandler) GetAllCats(c echo.Context) error {
	cats, err := ch.catService.GetAllCats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": cats}})
}

func (ch *CatHandler) DeleteCat(c echo.Context) error {
	ID := c.Param("id")
	catID, err := strconv.Atoi(ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	err, respStatus := ch.catService.DeleteCat(catID)
	if err != nil {
		return c.JSON(respStatus, response.UserResponse{Status: respStatus, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, response.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "cat successfully deleted"}})
}
