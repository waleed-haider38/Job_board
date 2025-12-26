package controllers

import (
	"net/http"
	"myjob/config"
	"myjob/models"

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
