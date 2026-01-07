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

    // Filters
    title := c.QueryParam("title")
    location := c.QueryParam("location")
    jobType := c.QueryParam("job_type")
    employerID := c.QueryParam("employer_id")
    salaryMin := c.QueryParam("salary_min")
    salaryMax := c.QueryParam("salary_max")

    // Base query
    query := config.GormDB.Model(&models.Job{}).
        Preload("Employer").
        Preload("Skills")

    // Apply filters dynamically
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
    if salaryMin != "" {
        query = query.Where("salary_min >= ?", salaryMin)
    }
    if salaryMax != "" {
        query = query.Where("salary_max <= ?", salaryMax)
    }

    // Count total records after filters
    if err := query.Count(&total).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    // Fetch paginated jobs
    if err := query.Limit(p.PerPage).Offset(p.Offset).Find(&jobs).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
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
