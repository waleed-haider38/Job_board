package controllers

import (
	"net/http"
	"myjob/config"
	"myjob/models"
	"myjob/utils"

	"github.com/labstack/echo/v4"
)

// CREATE Company
func CreateCompany(c echo.Context) error {
	company := new(models.Company)

	// Bind JSON body to struct
	if err := c.Bind(company); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := config.GormDB.Create(company).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, company)
}

// GET all Companies
func GetCompanies(c echo.Context) error {
	var companies []models.Company
	var total int64

	p := utils.GetPagination(c)

	config.GormDB.Model(&models.Company{}).Count(&total)

	if err := config.GormDB.
		Limit(p.PerPage).
		Offset(p.Offset).
		Find(&companies).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": companies,
		"meta": echo.Map{
			"page":      p.Page,
			"per_page": p.PerPage,
			"total":    total,
		},
	})
}


// GET single Company by ID
func GetCompanyByID(c echo.Context) error {
	id := c.Param("id")
	var company models.Company

	if err := config.GormDB.First(&company, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Company not found"})
	}

	return c.JSON(http.StatusOK, company)
}

// UPDATE Company
func UpdateCompany(c echo.Context) error {
	id := c.Param("id")
	var company models.Company

	if err := config.GormDB.First(&company, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Company not found"})
	}

	if err := c.Bind(&company); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	config.GormDB.Save(&company)
	return c.JSON(http.StatusOK, company)
}

// DELETE Company
func DeleteCompany(c echo.Context) error {
	id := c.Param("id")

	if err := config.GormDB.Delete(&models.Company{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Company deleted successfully"})
}
