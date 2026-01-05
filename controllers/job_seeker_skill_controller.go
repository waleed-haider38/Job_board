package controllers

import (
	"net/http"

	"myjob/config"
	"myjob/models"
	"myjob/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

/*
|--------------------------------------------------------------------------
| INPUT STRUCT
|--------------------------------------------------------------------------
*/

type SkillInput struct {
	SkillID int `json:"skill_id"`
}

/*
|--------------------------------------------------------------------------
| ADD SKILL TO LOGGED-IN JOB SEEKER
|--------------------------------------------------------------------------
| POST /job-seeker/skills
*/

func AddSkillToJobSeeker(c echo.Context) error {

	// 🔐 Get user from JWT
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

	// 🔎 Find job seeker
	var jobSeeker models.JobSeeker
	if err := config.GormDB.
		Where("user_id = ?", userID).
		First(&jobSeeker).Error; err != nil {

		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Job seeker profile not found",
		})
	}

	// 📥 Bind input
	var input SkillInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	// 🚫 Prevent duplicate skill
	var exists int64
	config.GormDB.
		Model(&models.JobSeekerSkill{}).
		Where("job_seeker_id = ? AND skill_id = ?", jobSeeker.JobSeekerID, input.SkillID).
		Count(&exists)

	if exists > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Skill already added",
		})
	}

	// ✅ Create relation
	jss := models.JobSeekerSkill{
		JobSeekerID: jobSeeker.JobSeekerID,
		SkillID:     input.SkillID,
	}

	if err := config.GormDB.Create(&jss).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Skill added successfully",
	})
}

/*
|--------------------------------------------------------------------------
| GET LOGGED-IN JOB SEEKER SKILLS
|--------------------------------------------------------------------------
| GET /job-seeker/skills?skill=go&page=1
*/

func GetMySkills(c echo.Context) error {

	// 🔐 JWT
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

	// 🔎 Job seeker
	var jobSeeker models.JobSeeker
	if err := config.GormDB.
		Where("user_id = ?", userID).
		First(&jobSeeker).Error; err != nil {

		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Job seeker not found",
		})
	}

	var skills []models.JobSeekerSkill
	var total int64

	p := utils.GetPagination(c)
	search := c.QueryParam("skill")

	query := config.GormDB.
		Model(&models.JobSeekerSkill{}).
		Joins("JOIN skills ON skills.skill_id = job_seeker_skills.skill_id").
		Where("job_seeker_skills.job_seeker_id = ?", jobSeeker.JobSeekerID).
		Preload("Skill")

	if search != "" {
		query = query.Where("skills.skill_name ILIKE ?", "%"+search+"%")
	}

	// 🔢 Count
	if err := query.Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	// 📄 Fetch
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

/*
|--------------------------------------------------------------------------
| REMOVE SKILL FROM LOGGED-IN JOB SEEKER
|--------------------------------------------------------------------------
| DELETE /job-seeker/skills/:skill_id
*/

func RemoveSkillFromJobSeeker(c echo.Context) error {

	// 🔐 JWT
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

	skillID := c.Param("skill_id")

	// 🔎 Job seeker
	var jobSeeker models.JobSeeker
	if err := config.GormDB.
		Where("user_id = ?", userID).
		First(&jobSeeker).Error; err != nil {

		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Job seeker not found",
		})
	}

	// ❌ Delete relation
	if err := config.GormDB.
		Delete(&models.JobSeekerSkill{},
			"job_seeker_id = ? AND skill_id = ?",
			jobSeeker.JobSeekerID,
			skillID,
		).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Skill removed successfully",
	})
}
