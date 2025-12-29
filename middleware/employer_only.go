package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func EmployerOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// 1️. role context se nikalo
		role, ok := c.Get("role").(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Role not found in token",
			})
		}

		// 2️. check karo role employer hai ya nahi
		if role != "employer" {
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": "Only employers can access this resource",
			})
		}

		// 3️. role sahi hai → aage jao
		return next(c)
	}
}
