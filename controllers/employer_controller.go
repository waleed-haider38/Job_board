package controllers

import (
	"net/http"
	"myjob/config"
	"myjob/models"

	"github.com/labstack/echo/v4"
)

// CREATE Employer
func CreateEmployer(c echo.Context) error {
	employer := new(models.Employer)

	//bind method request json ko struct mein   bind karta ha

	if err := c.Bind(employer); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	//GormDB sql query built krta ha 
		
	if err := config.GormDB.Create(employer).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, employer)
}

// GET all Employers
func GetEmployers(c echo.Context) error {
	var employers []models.Employer

	// linked user info b fetch krta ha.
	if err := config.GormDB.Preload("User").Find(&employers).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, employers)
}

// GET single Employer by ID
func GetEmployerByID(c echo.Context) error {
	id := c.Param("id")
	var employer models.Employer

	//Preload linked user info b fetch krta ha.
	if err := config.GormDB.Preload("User").First(&employer, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Employer not found"})
	}

	return c.JSON(http.StatusOK, employer)
}

// UPDATE Employer
func UpdateEmployer(c echo.Context) error {
	// param url sy id nikalta ha.(Built-in method)
	id := c.Param("id")
	var employer models.Employer

	//GormDB query create krta ha ju automatically. method k according.
	if err := config.GormDB.First(&employer, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Employer not found"})
	}

	/*bind method request k andr jo json data hta ha usko struct k sth bind krta ha 
	error ki surat  mein usy  print kr dyta ha.
	*/
	if err := c.Bind(&employer); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	

	config.GormDB.Save(&employer)

	return c.JSON(http.StatusOK, employer)
}

// DELETE Employer
func DeleteEmployer(c echo.Context) error {
	id := c.Param("id")

	if err := config.GormDB.Delete(&models.Employer{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Employer deleted successfully"})
}
