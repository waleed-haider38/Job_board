package controllers

import (
	"net/http"
	"myjob/config"
	"myjob/models"
	"myjob/utils"

	"github.com/labstack/echo/v4"
)

// CREATE Company (Private, Employer Only)
func CreateCompany(c echo.Context) error {
	company := new(models.Company)

	if err := c.Bind(company); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// JWT se employerID nikal kar assign karo
	userID := c.Get("user_id")
	var userIDInt int
	switch v := userID.(type) {
	case int:
		userIDInt = v
	case float64:
		userIDInt = int(v)
	default:
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid user_id type"})
	}

	// Employer table se actual employerID nikal lo
	var employer models.Employer
	if err := config.GormDB.Where("user_id = ?", userIDInt).First(&employer).Error; err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Employer profile not found"})
	}

	company.EmployerID = employer.EmployerID

	if err := config.GormDB.Create(company).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, company)
}

// GET all companies (Public, with filters + pagination)
func GetCompanies(c echo.Context) error {
	var companies []models.Company
	var total int64
	p := utils.GetPagination(c)

	name := c.QueryParam("name")
	product := c.QueryParam("product")
	employerID := c.QueryParam("employer_id")

	query := config.GormDB.Model(&models.Company{})

	if name != "" {
		query = query.Where("company_name ILIKE ?", "%"+name+"%")
	}
	if product != "" {
		query = query.Where("company_product ILIKE ?", "%"+product+"%")
	}
	if employerID != "" {
		query = query.Where("employer_id = ?", employerID)
	}

	if err := query.Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if err := query.Limit(p.PerPage).Offset(p.Offset).Find(&companies).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": companies,
		"meta": echo.Map{
			"page":     p.Page,
			"per_page": p.PerPage,
			"total":    total,
		},
	})
}

// GET single Company by ID (Private, owner only)
func GetCompanyByID(c echo.Context) error {
	id := c.Param("id")
	var company models.Company

	if err := config.GormDB.First(&company, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Company not found"})
	}

	// Ownership check
	userID := c.Get("user_id")
	var userIDInt int
	switch v := userID.(type) {
	case int:
		userIDInt = v
	case float64:
		userIDInt = int(v)
	default:
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid user_id type"})
	}

	var employer models.Employer
	if err := config.GormDB.Where("user_id = ?", userIDInt).First(&employer).Error; err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Employer profile not found"})
	}

	if company.EmployerID != employer.EmployerID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "You cannot access this company"})
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

	// Ownership check (same as GetCompanyByID)
	userID := c.Get("user_id")
	var userIDInt int
	switch v := userID.(type) {
	case int:
		userIDInt = v
	case float64:
		userIDInt = int(v)
	default:
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid user_id type"})
	}

	var employer models.Employer
	if err := config.GormDB.Where("user_id = ?", userIDInt).First(&employer).Error; err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Employer profile not found"})
	}

	if company.EmployerID != employer.EmployerID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "You cannot update this company"})
	}

	// Bind new data
	if err := c.Bind(&company); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := config.GormDB.Save(&company).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, company)
}

// DELETE Company
func DeleteCompany(c echo.Context) error {
	id := c.Param("id")
	var company models.Company

	if err := config.GormDB.First(&company, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Company not found"})
	}

	// Ownership check
	userID := c.Get("user_id")
	var userIDInt int
	switch v := userID.(type) {
	case int:
		userIDInt = v
	case float64:
		userIDInt = int(v)
	default:
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid user_id type"})
	}

	var employer models.Employer
	if err := config.GormDB.Where("user_id = ?", userIDInt).First(&employer).Error; err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Employer profile not found"})
	}

	if company.EmployerID != employer.EmployerID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "You cannot delete this company"})
	}

	if err := config.GormDB.Delete(&company).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Company deleted successfully"})
}
