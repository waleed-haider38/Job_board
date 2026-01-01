package controllers

import (
	"net/http"
	"myjob/config"
	"myjob/models"
	"myjob/utils"

	"github.com/labstack/echo/v4"
)

// CREATE Skill
func CreateSkill(c echo.Context) error {
	skill := new(models.Skill)

	if err := c.Bind(skill); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := config.GormDB.Create(skill).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, skill)
}

// GET all Skills
func GetSkills(c echo.Context) error {
	var skills []models.Skill
	var total int64

	// get pagination values
	p := utils.GetPagination(c)

	// count total skills
	if err := config.GormDB.
		Model(&models.Skill{}).
		Count(&total).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// fetch paginated skills
	if err := config.GormDB.
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


// GET Skill by ID
func GetSkillByID(c echo.Context) error {
	id := c.Param("id")
	var skill models.Skill

	if err := config.GormDB.First(&skill, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Skill not found"})
	}

	return c.JSON(http.StatusOK, skill)
}

// UPDATE Skill
func UpdateSkill(c echo.Context) error {
	id := c.Param("id")
	var skill models.Skill

	if err := config.GormDB.First(&skill, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Skill not found"})
	}

	if err := c.Bind(&skill); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	config.GormDB.Save(&skill)
	return c.JSON(http.StatusOK, skill)
}

// DELETE Skill
func DeleteSkill(c echo.Context) error {
	id := c.Param("id")

	if err := config.GormDB.Delete(&models.Skill{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Skill deleted successfully"})
}
