package controllers

import (
	"fmt"
	"net/http"
	"time"

	"myjob/config"
	"myjob/models"

	"github.com/labstack/echo/v4"
)

// ===============================
// Job Seeker applies to a Job
// ===============================
func ApplyToJob(c echo.Context) error {

	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "invalid token",
		})
	}

	// Job seeker fetch
	var jobSeeker models.JobSeeker
	if err := config.GormDB.
		Where("user_id = ?", userID).
		First(&jobSeeker).Error; err != nil {

		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Job seeker profile not found",
		})
	}

	type ApplyInput struct {
		JobID       int    `json:"job_id"`
		CoverLetter string `json:"cover_letter"`
	}

	var input ApplyInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Job exists check
	var job models.Job
	if err := config.GormDB.First(&job, input.JobID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Job not found"})
	}

	// Duplicate check
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
		Status:      "pending",
		AppliedAt:   time.Now(),
	}

	if err := config.GormDB.Create(&application).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, application)
}

// ===============================
// Job Seeker views own applications
// ===============================
func GetMyApplications(c echo.Context) error {

	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
	}

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

// ===============================
// Employer views applications for a Job
// ===============================
func GetApplicationsForJob(c echo.Context) error {

	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
	}

	jobID := c.Param("job_id")

	// Employer fetch
	var employer models.Employer
	if err := config.GormDB.
		Where("user_id = ?", userID).
		First(&employer).Error; err != nil {

		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Employer not found",
		})
	}

	// Job + Company ownership check
	var job models.Job
	if err := config.GormDB.
		Preload("Company").
		Where("job_id = ?", jobID).
		First(&job).Error; err != nil {

		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Job not found",
		})
	}

	fmt.Printf("JOB: %+v\n", job)
	fmt.Printf("COMPANY: %+v\n", job.Company)
	fmt.Println("EmployerID from token:", employer.EmployerID)

	if job.Company == nil || job.Company.EmployerID != employer.EmployerID {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "You are not allowed to view applications for this job",
		})
	}

	var applications []models.Application
	if err := config.GormDB.
		Where("job_id = ?", jobID).
		Preload("JobSeeker").
		Find(&applications).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, applications)
}

// ===============================
// Employer updates application status
// ===============================
func UpdateApplicationStatus(c echo.Context) error {

	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
	}

	appID := c.Param("id")

	type StatusInput struct {
		Status string `json:"status"`
	}

	var input StatusInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	allowed := map[string]bool{
		"reviewed":    true,
		"shortlisted": true,
		"rejected":    true,
		"hired":       true,
	}

	if !allowed[input.Status] {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid status"})
	}

	// Employer fetch
	var employer models.Employer
	if err := config.GormDB.
		Where("user_id = ?", userID).
		First(&employer).Error; err != nil {

		return c.JSON(http.StatusForbidden, echo.Map{"error": "Employer not found"})
	}

	// Ownership check via job -> company -> employer
	var application models.Application
	if err := config.GormDB.
		Joins("JOIN jobs ON jobs.job_id = applications.job_id").
		Joins("JOIN companies ON companies.company_id = jobs.company_id").
		Where(
			"applications.application_id = ? AND companies.employer_id = ?",
			appID,
			employer.EmployerID,
		).
		First(&application).Error; err != nil {

		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "You are not allowed to update this application",
		})
	}

	if err := config.GormDB.
		Model(&application).
		Update("status", input.Status).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Application status updated successfully",
	})
}
