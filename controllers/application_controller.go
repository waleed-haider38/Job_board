package controllers

import (
	"net/http"
	"time"

	"myjob/config"
	"myjob/models"
	"myjob/utils"

	//"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Apply Jobs
func ApplyToJob(c echo.Context) error {

	//  JWT se user_id nikalo
	userID := c.Get("user_id").(int)

	// user → job seeker mapping
	var jobSeeker models.JobSeeker
	if err := config.GormDB.
		Where("user_id = ?", userID).
		First(&jobSeeker).Error; err != nil {

		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Only job seekers can apply",
		})
	}

	type ApplyInput struct {
		JobID       int    `json:"job_id"`
		CoverLetter string `json:"cover_letter"`
	}

	input := new(ApplyInput)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	//  Duplicate apply check
	var existing models.Application
	if err := config.GormDB.
		Where("job_id = ? AND job_seeker_id = ?", input.JobID, jobSeeker.JobSeekerID).
		First(&existing).Error; err == nil {

		return c.JSON(http.StatusConflict, echo.Map{
			"error": "You have already applied to this job",
		})
	}

	application := models.Application{
		JobID:       input.JobID,
		JobSeekerID: jobSeeker.JobSeekerID,
		CoverLetter: input.CoverLetter,
		Status:      "applied",
	}

	if err := config.GormDB.Create(&application).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, application)
}



//Application Detail code when we know who is the applicant
// Get applications of a specific Job Seeker
func GetMyApplications(c echo.Context) error {

	userID := c.Get("user_id").(int)

	var jobSeeker models.JobSeeker
	if err := config.GormDB.
		Where("user_id = ?", userID).
		First(&jobSeeker).Error; err != nil {

		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Job seeker not found",
		})
	}

	var applications []models.Application

	if err := config.GormDB.
		Where("job_seeker_id = ?", jobSeeker.JobSeekerID).
		Preload("Job").
		Find(&applications).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, applications)
}


// Get all applications for a specific job (Employer view)
func GetApplicationsForJob(c echo.Context) error {
	jobID := c.Param("job_id")

	var applications []models.Application

	if err := config.GormDB.
		Where("job_id = ?", jobID).
		Preload("JobSeeker").
		Find(&applications).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, applications)
}

// Update the status of an application (Employer action)
func UpdateApplicationStatus(c echo.Context) error {
	id := c.Param("id")

	type StatusInput struct {
		Status string `json:"status"`
	}

	input := new(StatusInput)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	allowed := map[string]bool{
		"reviewed":    true,
		"shortlisted": true,
		"rejected":    true,
		"hired":       true,
	}
	var app models.Application
	config.GormDB.First(&app, id)

	// Current status
	switch app.Status {
	case "pending":
		if input.Status != "reviewed" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid status transition from pending"})
		}
	case "reviewed":
		if input.Status != "shortlisted" && input.Status != "rejected" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid status transition from reviewed"})
		}
	case "shortlisted":
		if input.Status != "hired" && input.Status != "rejected" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid status transition from shortlisted"})
		}
	default:
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Cannot change status from " + app.Status})
	}


	if !allowed[input.Status] {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid status",
		})
	}

	if err := config.GormDB.
		Model(&models.Application{}).
		Where("application_id = ?", id).
		Update("status", input.Status).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Application status updated",
	})
}



// Create Application
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

// Get Application
func GetApplications(c echo.Context) error {
	var applications []models.Application
	var total int64

	// Pagination
	p := utils.GetPagination(c)

	// Query params (filters)
	jobID := c.QueryParam("job_id")
	jobSeekerID := c.QueryParam("job_seeker_id")
	status := c.QueryParam("status")

	// Base query
	query := config.GormDB.
		Model(&models.Application{}).
		Preload("Job").
		Preload("JobSeeker")

	// Dynamic filters
	if jobID != "" {
		query = query.Where("job_id = ?", jobID)
	}

	if jobSeekerID != "" {
		query = query.Where("job_seeker_id = ?", jobSeekerID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count AFTER filters
	if err := query.Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	// Fetch paginated data
	if err := query.
		Limit(p.PerPage).
		Offset(p.Offset).
		Find(&applications).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	// Response
	return c.JSON(http.StatusOK, echo.Map{
		"data": applications,
		"meta": echo.Map{
			"page":     p.Page,
			"per_page": p.PerPage,
			"total":    total,
		},
	})
}

// Get Application by ID

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

// Update the Application

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

// Delete the Application

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
