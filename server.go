package main

import (
	"fmt"
	"myjob/config"
	"net/http"
	"myjob/controllers" 
	"myjob/middleware" 

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Connect to database
	db := config.ConnectDB()
	fmt.Println(db)

	// Example: you can use db in handlers
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World! I am Waleed. Database is connected successfully!")
	})
	e.POST("/api/register", controllers.Register)
	e.POST("/api/login", controllers.Login)

	api := e.Group("/api")
	api.Use(middleware.JWTMiddleware)
	api.GET("/me", controllers.Me)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
