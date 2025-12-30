package controllers

import (
	"net/http"
	"time"

	"myjob/config"
	"myjob/models"

	"github.com/labstack/echo/v4"
)

func CreateApplication(c echo.Context) error {
	app := new(models.Application)

	if err := c.Bind(app); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	// auto set applied time
	app.AppliedAt = time.Now()

	if err := config.GormDB.Create(app).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, app)
}

func GetApplications(c echo.Context) error {
	var applications []models.Application

	if err := config.GormDB.
		Preload("Job").
		Preload("JobSeeker").
		Find(&applications).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, applications)
}

func GetApplicationByID(c echo.Context) error {
	id := c.Param("id")
	var app models.Application

	if err := config.GormDB.
		Preload("Job").
		Preload("JobSeeker").
		First(&app, id).Error; err != nil {

		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Application not found",
		})
	}

	return c.JSON(http.StatusOK, app)
}


func UpdateApplication(c echo.Context) error {
	id := c.Param("id")
	var app models.Application

	if err := config.GormDB.First(&app, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Application not found",
		})
	}

	if err := c.Bind(&app); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if err := config.GormDB.Save(&app).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, app)
}


func DeleteApplication(c echo.Context) error {
	id := c.Param("id")

	if err := config.GormDB.Delete(&models.Application{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Application deleted successfully",
	})
}
