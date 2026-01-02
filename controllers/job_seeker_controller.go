package controllers

import (
	"net/http"
	"myjob/config"
	"myjob/models"
	"myjob/utils"

	"github.com/labstack/echo/v4"
)


// CREATE Job Seeker
func CreateJobSeeker(c echo.Context) error {
	jobSeeker := new(models.JobSeeker)

	// JSON → struct
	if err := c.Bind(jobSeeker); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	// Save to DB
	if err := config.GormDB.Create(jobSeeker).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, jobSeeker)
}


// GET all Job Seekers
func GetJobSeekers(c echo.Context) error {
	var jobSeekers []models.JobSeeker
	var total int64

	// Pagination
	p := utils.GetPagination(c)

	// Query params (filters)
	name := c.QueryParam("name")
	userID := c.QueryParam("user_id")

	// Base query
	query := config.GormDB.
		Model(&models.JobSeeker{}).
		Preload("User")

	// Dynamic filters
	if name != "" {
		query = query.Where("full_name ILIKE ?", "%"+name+"%")
	}

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Count AFTER filters
	if err := query.Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	// Fetch paginated job seekers
	if err := query.
		Limit(p.PerPage).
		Offset(p.Offset).
		Find(&jobSeekers).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": jobSeekers,
		"meta": echo.Map{
			"page":      p.Page,
			"per_page": p.PerPage,
			"total":    total,
		},
	})
}



// GET Job Seeker by ID
func GetJobSeekerByID(c echo.Context) error {
	id := c.Param("id")
	var jobSeeker models.JobSeeker

	if err := config.GormDB.Preload("User").First(&jobSeeker, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Job seeker not found",
		})
	}

	return c.JSON(http.StatusOK, jobSeeker)
}


// UPDATE Job Seeker
func UpdateJobSeeker(c echo.Context) error {
	id := c.Param("id")
	var jobSeeker models.JobSeeker

	if err := config.GormDB.First(&jobSeeker, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Job seeker not found",
		})
	}

	if err := c.Bind(&jobSeeker); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	config.GormDB.Save(&jobSeeker)
	return c.JSON(http.StatusOK, jobSeeker)
}


// DELETE Job Seeker
func DeleteJobSeeker(c echo.Context) error {
	id := c.Param("id")

	if err := config.GormDB.Delete(&models.JobSeeker{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Job seeker deleted successfully",
	})
}
