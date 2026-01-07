package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func JobSeekerOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func (c echo.Context) error {

		// we need to fectch role from the context
		role , ok := c.Get("role").(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized,echo.Map{
				"message": "Role not Found",
			})
		}

		//Check the role that role is of job seeker.
		if role != "job_seeker"{
			return c.JSON(http.StatusForbidden,echo.Map{
				"message": "Only Job Seeker can apply",
			})
		}

		//if role is correct pass next
		return next(c)
	}
}

