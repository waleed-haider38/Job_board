package controllers

import (
	"net/http"
	"myjob/config"
	"myjob/models"
	"myjob/utils"

	"github.com/labstack/echo/v4"
)

// CREATE / ASSIGN a skill to a Job Seeker
func AddSkillToJobSeeker(c echo.Context) error {
	jss := new(models.JobSeekerSkill)

	if err := c.Bind(jss); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := config.GormDB.Create(jss).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, jss)
}

// GET all skills of a Job Seeker
func GetSkillsOfJobSeeker(c echo.Context) error {
	jobSeekerID := c.Param("id")

	var skills []models.JobSeekerSkill
	var total int64

	// get pagination values
	p := utils.GetPagination(c)

	// count total records
	if err := config.GormDB.
		Model(&models.JobSeekerSkill{}).
		Where("job_seeker_id = ?", jobSeekerID).
		Count(&total).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// fetch paginated data
	if err := config.GormDB.
		Preload("Skill").
		Where("job_seeker_id = ?", jobSeekerID).
		Limit(p.PerPage).
		Offset(p.Offset).
		Find(&skills).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": skills,
		"meta": echo.Map{
			"page":     p.Page,
			"per_page": p.PerPage,
			"total":    total,
		},
	})
}


// DELETE a skill from a Job Seeker
func RemoveSkillFromJobSeeker(c echo.Context) error {
	jobSeekerID := c.Param("jobSeekerID")
	skillID := c.Param("skillID")

	if err := config.GormDB.Delete(&models.JobSeekerSkill{}, "job_seeker_id = ? AND skill_id = ?", jobSeekerID, skillID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Skill removed successfully"})
}
