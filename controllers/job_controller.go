package controllers

import (
	"time"
	"net/http"
	"myjob/config"
	"myjob/models"
	"myjob/utils"

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
	var total int64

	// Pagination
	p := utils.GetPagination(c)

	// Query params (filters)
	title := c.QueryParam("title")
	location := c.QueryParam("location")
	jobType := c.QueryParam("job_type")
	employerID := c.QueryParam("employer_id")
	salaryMin := c.QueryParam("salary_min")
	salaryMax := c.QueryParam("salary_max")

	// Base query
	query := config.GormDB.
		Model(&models.Job{}).
		Preload("Employer")

	// Dynamic filters
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	if location != "" {
		query = query.Where("job_location ILIKE ?", "%"+location+"%")
	}

	if jobType != "" {
		query = query.Where("job_type = ?", jobType)
	}

	if employerID != "" {
		query = query.Where("employer_id = ?", employerID)
	}

	// Salary range filters
	if salaryMin != "" {
		query = query.Where("salary_min >= ?", salaryMin)
	}

	if salaryMax != "" {
		query = query.Where("salary_max <= ?", salaryMax)
	}

	// Count AFTER filters
	if err := query.Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	// Fetch paginated jobs
	if err := query.
		Limit(p.PerPage).
		Offset(p.Offset).
		Find(&jobs).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": jobs,
		"meta": echo.Map{
			"page":      p.Page,
			"per_page": p.PerPage,
			"total":    total,
		},
	})
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
