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

	// Pagination
	p := utils.GetPagination(c)

	// Query param (search)
	skill := c.QueryParam("skill") // e.g. ?skill=go

	// Base query (job_seeker_id is mandatory)
	query := config.GormDB.
		Model(&models.JobSeekerSkill{}).
		Joins("JOIN skills ON skills.skill_id = job_seeker_skills.skill_id").
		Where("job_seeker_skills.job_seeker_id = ?", jobSeekerID).
		Preload("Skill")

	// Optional skill-name search
	if skill != "" {
		query = query.Where("skills.skill_name ILIKE ?", "%"+skill+"%")
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
		Find(&skills).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": skills,
		"meta": echo.Map{
			"page":      p.Page,
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
