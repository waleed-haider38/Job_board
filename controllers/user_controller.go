package controllers

import (
	//"fmt"
	"myjob/config"
	"myjob/models"
	"myjob/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Me(c echo.Context) error {

	// 1️. context se user_id nikalo
	userID := c.Get("user_id").(int)

	db := config.ConnectDB()

	user := models.User{}

	// 2️ DB se user fetch karo
	err := db.QueryRow(
		`SELECT user_id, email, role 
		 FROM users 
		 WHERE user_id = $1`,
		userID,
	).Scan(&user.UserID, &user.Email, &user.Role)

	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
		})
	}

	// 3️. response
	return c.JSON(http.StatusOK, echo.Map{
		"user_id": user.UserID,
		"email":   user.Email,
		"role":    user.Role,
	})
}

// CREATE USER
func CreateUser(c echo.Context) error {
	db := config.GormDB
	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if err := db.Create(user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, user)
}

//Get User

func GetUsers(c echo.Context) error {
	var users []models.User
	var total int64

	db := config.GormDB

	// Pagination
	p := utils.GetPagination(c)

	// Query params (search filters)
	email := c.QueryParam("email")
	role := c.QueryParam("role")
	isActive := c.QueryParam("is_active") // "true" / "false"

	// Base query
	query := db.Model(&models.User{})

	// Dynamic filters
	if email != "" {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}

	if role != "" {
		query = query.Where("role = ?", role)
	}

	if isActive != "" {
		if isActive == "true" {
			query = query.Where("is_active = ?", true)
		} else if isActive == "false" {
			query = query.Where("is_active = ?", false)
		}
	}

	// Count AFTER filters
	if err := query.Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Fetch paginated data
	if err := query.
		Limit(p.PerPage).
		Offset(p.Offset).
		Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Return response
	return c.JSON(http.StatusOK, echo.Map{
		"data": users,
		"meta": echo.Map{
			"page":     p.Page,
			"per_page": p.PerPage,
			"total":    total,
		},
	})
}


//Get User by ID

func GetUserByID(c echo.Context) error {
	db := config.GormDB
	id := c.Param("id")

	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "User not found",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// Update User

func UpdateUser(c echo.Context) error {
	db := config.GormDB
	id := c.Param("id")

	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "User not found",
		})
	}

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	db.Save(&user)

	return c.JSON(http.StatusOK, user)
}

// Delete User

func DeleteUser(c echo.Context) error {
	db := config.GormDB
	id := c.Param("id")

	if err := db.Delete(&models.User{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "User deleted successfully",
	})
}
