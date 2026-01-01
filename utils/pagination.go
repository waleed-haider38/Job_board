package utils

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

type Pagination struct {
	Page    int
	PerPage int
	Offset  int
}

func GetPagination(c echo.Context) Pagination {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(c.QueryParam("per_page"))
	if err != nil || perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	return Pagination{
		Page:    page,
		PerPage: perPage,
		Offset:  offset,
	}
}
