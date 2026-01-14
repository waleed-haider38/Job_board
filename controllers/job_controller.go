package controllers

import (
	"myjob/config"
	"myjob/models"
	"myjob/utils"
	"net/http"
	"time"

	//"github.com/golang-jwt/jwt/v5"
	// /"fmt"

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

	// Filters
	title := c.QueryParam("title")
	location := c.QueryParam("location")
	jobType := c.QueryParam("job_type")
	companyID := c.QueryParam("company_id")
	salaryMin := c.QueryParam("salary_min")
	salaryMax := c.QueryParam("salary_max")

	// Base query
	query := config.GormDB.Model(&models.Job{}).
		Preload("Company").
		Preload("Skills")

	// Apply filters
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if location != "" {
		query = query.Where("job_location ILIKE ?", "%"+location+"%")
	}
	if jobType != "" {
		query = query.Where("job_type = ?", jobType)
	}
	if companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}
	if salaryMin != "" {
		query = query.Where("salary_min >= ?", salaryMin)
	}
	if salaryMax != "" {
		query = query.Where("salary_max <= ?", salaryMax)
	}

	// Count
	if err := query.Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	// Fetch data
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
			"page":     p.Page,
			"per_page": p.PerPage,
			"total":    total,
		},
	})
}

// GET single Job by ID
func GetJobByID(c echo.Context) error {
	id := c.Param("id")
	var job models.Job

	if err := config.GormDB.
		Preload("Company").
		Preload("Skills").
		First(&job, id).Error; err != nil {

		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Job not found",
		})
	}

	return c.JSON(http.StatusOK, job)
}

// UPDATE Job (Employer Only, Safe)
func UpdateJob(c echo.Context) error {

	// 1️ JWT se user_id lo
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "missing or invalid token",
		})
	}

	// user_id ko int mein convert karo
	var userIDInt int
	switch v := userID.(type) {
	case int:
		userIDInt = v
	case float64:
		userIDInt = int(v)
	default:
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "invalid user_id type",
		})
	}

	// 2️ Employer table se employer nikaalo (by user_id)
	var employer models.Employer
	if err := config.GormDB.
		Where("user_id = ?", userIDInt).
		First(&employer).Error; err != nil {

		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Employer profile not found",
		})
	}

	// 3️ Job ID URL se lo
	jobID := c.Param("id")

	// 4️ Job load karo Company + Skills ke sath
	var job models.Job
	if err := config.GormDB.
		Preload("Company").
		Preload("Skills").
		First(&job, jobID).Error; err != nil {

		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Job not found",
		})
	}

	// 5️ Company existence check
	if job.Company == nil {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Job has no associated company",
		})
	}

	// 6️ OWNERSHIP CHECK (REAL FIX HERE)
	if job.Company.EmployerID != employer.EmployerID {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "You cannot update this job",
		})
	}

	// 7️ Safe update input
	type JobUpdateInput struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		JobType     string `json:"job_type"`
		JobLocation string `json:"job_location"`
		SalaryMin   int    `json:"salary_min"`
		SalaryMax   int    `json:"salary_max"`
		Status      string `json:"status"`
	}

	var input JobUpdateInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	// 8️ Fields update karo
	job.Title = input.Title
	job.Description = input.Description
	job.JobType = input.JobType
	job.JobLocation = input.JobLocation
	job.SalaryMin = input.SalaryMin
	job.SalaryMax = input.SalaryMax
	job.Status = input.Status
	job.UpdatedAt = time.Now()

	// 9️ Save
	if err := config.GormDB.Save(&job).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	// 10 Reload for clean response
	if err := config.GormDB.
		Preload("Company").
		Preload("Skills").
		First(&job, job.ID).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to reload job",
		})
	}

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

// AddSkillToJob adds a skill to an existing job (many-to-many relationship)
func AddSkillToJob(c echo.Context) error {

	// 1. Get the job ID from the URL parameter
	jobID := c.Param("job_id")

	// 2. Define the expected JSON input structure
	type SkillInput struct {
		SkillID int `json:"skill_id"` // ID of the skill to add
	}

	// 3. Bind the JSON body to SkillInput struct
	input := new(SkillInput)
	if err := c.Bind(input); err != nil {
		// If JSON is invalid, return 400 Bad Request
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// 4. Fetch the job from the database
	var job models.Job
	if err := config.GormDB.First(&job, jobID).Error; err != nil {
		// If job does not exist, return 404 Not Found
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Job not found"})
	}

	// 5. Fetch the skill from the database
	var skill models.Skill
	if err := config.GormDB.First(&skill, input.SkillID).Error; err != nil {
		// If skill does not exist, return 404 Not Found
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Skill not found"})
	}

	// 6. Append the skill to the job's Skills association (many-to-many)
	if err := config.GormDB.Model(&job).Association("Skills").Append(&skill); err != nil {
		// If there is a database error, return 500 Internal Server Error
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// 7. Return success message
	return c.JSON(http.StatusOK, echo.Map{"message": "Skill added to job"})
}
