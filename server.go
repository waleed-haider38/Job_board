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
	config.ConnectGorm(e)
	fmt.Println(db)

	// Example: you can use db in handlers
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World! I am Waleed. Database is connected successfully!")
	})	
	e.GET("/user", func(c echo.Context) error {
		return c.String(http.StatusOK, "Assalam-o-Alaikum! I am Waleed.")
	})	

	// To register a User.
	e.POST("/api/register", controllers.Register)
	//To login a User.
	e.POST("/api/login", controllers.Login)

	//Only register user can access this route.

	api := e.Group("/api")
	api.Use(middleware.JWTMiddleware)
	api.GET("/me", controllers.Me)

	//user CRUD routes
	e.POST("/users", controllers.CreateUser)
	e.GET("/users", controllers.GetUsers)
	e.GET("/users/:id", controllers.GetUserByID)
	e.PUT("/users/:id", controllers.UpdateUser)
	e.DELETE("/users/:id", controllers.DeleteUser)

	//Employer CRUD Routes
	e.POST("/employers", controllers.CreateEmployer)
	e.GET("/employers", controllers.GetEmployers)
	e.GET("/employers/:id", controllers.GetEmployerByID)
	e.PUT("/employers/:id", controllers.UpdateEmployer)
	e.DELETE("/employers/:id", controllers.DeleteEmployer)



	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
