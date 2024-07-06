package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"spyCat/database"
	"spyCat/handler"
	"spyCat/service"
)

var validate = validator.New()
var catHandler = handler.NewCatHandler(service.NewCatService(database.NewCatDatabase(database.NewDatabase()), validate))
var missionHandler = handler.NewMissionHandler(service.NewMissionService(database.NewMissionDatabase(database.NewDatabase()), validate))
var targetHandler = handler.NewTargetHandler(service.NewTargetService(database.NewTargetDatabase(database.NewDatabase()), validate))

func UserRoute(e *echo.Echo) {
	e.POST("/cats", catHandler.CreateCat)
	e.GET("/cats/:id", catHandler.GetCat)
	e.GET("/cats", catHandler.GetAllCats)
	e.PUT("/cats/:id", catHandler.UpdateCatSalary)
	e.DELETE("/cats/:id", catHandler.DeleteCat)

	e.POST("/missions", missionHandler.CreateMission)
	e.DELETE("/missions/:id", missionHandler.DeleteMission)
	e.PUT("/missions/:id/complete", missionHandler.CompleteMission)
	e.PUT("/missions/:id/assign", missionHandler.AssignCatToMission)
	e.GET("/missions", missionHandler.ListMissions)
	e.GET("/missions/:id", missionHandler.GetMission)

	e.PUT("/targets", targetHandler.UpdateTarget)
	e.PUT("/targets/:id/notes", targetHandler.UpdateTargetNotes)
	e.PUT("/missions/:missionId/targets/:targetId/complete", targetHandler.CompleteTarget)
	e.DELETE("/missions/:missionId/targets/:targetId", targetHandler.DeleteTarget)
	e.POST("/missions/:missionId/targets", targetHandler.AddTarget)

}
