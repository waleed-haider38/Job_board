package controllers

import (
	"time"
	"net/http"
	"myjob/config"
	"myjob/models"

	"github.com/labstack/echo/v4"
)

// CREATE Job
func CreateJob(c echo.Context) error {
	job := new(models.Job)

	// Bind JSON body to struct
	if err := c.Bind(job); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Optional: set default timestamps if not using autoCreateTime
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()

	if err := config.GormDB.Create(job).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, job)
}

// GET all Jobs
func GetJobs(c echo.Context) error {
	var jobs []models.Job

	if err := config.GormDB.Preload("Employer").Find(&jobs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, jobs)
}

// GET single Job by ID
func GetJobByID(c echo.Context) error {
	id := c.Param("id")
	var job models.Job

	if err := config.GormDB.Preload("Employer").First(&job, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Job not found"})
	}

	return c.JSON(http.StatusOK, job)
}

// UPDATE Job
func UpdateJob(c echo.Context) error {
	id := c.Param("id")
	var job models.Job

	if err := config.GormDB.First(&job, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Job not found"})
	}

	if err := c.Bind(&job); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	job.UpdatedAt = time.Now()
	config.GormDB.Save(&job)

	return c.JSON(http.StatusOK, job)
}

// DELETE Job
func DeleteJob(c echo.Context) error {
	id := c.Param("id")

	if err := config.GormDB.Delete(&models.Job{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Job deleted successfully"})
}
