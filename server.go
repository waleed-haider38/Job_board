package main

import (
	"fmt"
	"myjob/config"
	"myjob/controllers"
	"myjob/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Connect to database
	db := config.ConnectDB()

	// Connection to database using GORM.
	config.ConnectGorm()
	fmt.Println(db)

	// Example: you can use db in handlers
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World! I am Waleed. Database is connected successfully!")
	})
	// To register a User.
	e.POST("/api/register", controllers.Register)
	//To login a User.
	e.POST("/api/login", controllers.Login)

	// In this route only register user can go.

	api := e.Group("/api")
	api.Use(middleware.JWTMiddleware)
	api.GET("/me", controllers.Me)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
