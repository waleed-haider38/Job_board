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

	// Health check routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World! Database connected successfully!")
	})
	e.GET("/user", func(c echo.Context) error {
		return c.String(http.StatusOK, "Assalam-o-Alaikum!")
	})

	// Auth routes
	e.POST("/api/register", controllers.Register)
	e.POST("/api/login", controllers.Login)

	// Protected API group
	api := e.Group("/api", middleware.JWTMiddleware)
	api.GET("/me", controllers.Me)

	// User CRUD
	e.POST("/users", controllers.CreateUser)
	e.GET("/users", controllers.GetUsers)
	e.GET("/users/:id", controllers.GetUserByID)
	e.PUT("/users/:id", controllers.UpdateUser)
	e.DELETE("/users/:id", controllers.DeleteUser)

	// Employer CRUD
	e.POST("/employers", controllers.CreateEmployer)
	e.GET("/employers", controllers.GetEmployers)
	e.GET("/employers/:id", controllers.GetEmployerByID)
	e.PUT("/employers/:id", controllers.UpdateEmployer)
	e.DELETE("/employers/:id", controllers.DeleteEmployer)

	job := e.Group("/emp")
	job.Use(middleware.JWTMiddleware)
	job.Use(middleware.EmployerOnly)

	// Job CRUD
	job.POST("/jobs", controllers.CreateJob)
	e.GET("/jobs", controllers.GetJobs)
	e.GET("/jobs/:id", controllers.GetJobByID)
	job.PUT("/jobs/:id", controllers.UpdateJob)
	job.DELETE("/jobs/:id", controllers.DeleteJob)

	// Add skill to a job (Employer only)
	e.POST("/jobs/:job_id/skills", controllers.AddSkillToJob, middleware.JWTMiddleware, middleware.EmployerOnly)

	// Job Seeker CRUD
	e.POST("/job-seekers", controllers.CreateJobSeeker)
	e.GET("/job-seekers", controllers.GetJobSeekers)
	e.GET("/job-seekers/:id", controllers.GetJobSeekerByID)
	e.PUT("/job-seekers/:id", controllers.UpdateJobSeeker)
	e.DELETE("/job-seekers/:id", controllers.DeleteJobSeeker)

	// Application CRUD (basic)
	e.POST("/applications", controllers.CreateApplication)
	e.GET("/applications", controllers.GetApplications)
	e.GET("/applications/:id", controllers.GetApplicationByID)
	e.PUT("/applications/:id", controllers.UpdateApplication)
	e.DELETE("/applications/:id", controllers.DeleteApplication)

	// Job Seeker applies to a job
	e.POST(
		"/jobs/apply",
		controllers.ApplyToJob,
		middleware.JWTMiddleware,
		middleware.JobSeekerOnly,
	)

	// Job Seeker views their own applications
	e.GET(
		"/my-applications",
		controllers.GetMyApplications,
		middleware.JWTMiddleware,
		middleware.JobSeekerOnly,
	)

	// Employer views applications for a job
	e.GET(
		"/jobs/:job_id/applications",
		controllers.GetApplicationsForJob,
		middleware.JWTMiddleware,
		middleware.EmployerOnly,
	)

	// Employer updates application status
	e.PATCH(
		"/applications/:id/status",
		controllers.UpdateApplicationStatus,
		middleware.JWTMiddleware,
		middleware.EmployerOnly,
	)

	// Skills CRUD
	e.POST("/skills", controllers.CreateSkill)
	e.GET("/skills", controllers.GetSkills)
	e.GET("/skills/:id", controllers.GetSkillByID)
	e.PUT("/skills/:id", controllers.UpdateSkill)
	e.DELETE("/skills/:id", controllers.DeleteSkill)

	// Job Seeker skills (protected + role based)
	jobSeeker := e.Group(
		"/job-seeker",
		middleware.JWTMiddleware,
		middleware.JobSeekerOnly,
	)
	jobSeeker.POST("/skills", controllers.AddSkillToJobSeeker)
	jobSeeker.GET("/skills", controllers.GetMySkills)
	jobSeeker.DELETE("/skills/:skill_id", controllers.RemoveSkillFromJobSeeker)

	// Company CRUD
	e.POST("/companies", controllers.CreateCompany)
	e.GET("/companies", controllers.GetCompanies)
	e.GET("/companies/:id", controllers.GetCompanyByID)
	e.PUT("/companies/:id", controllers.UpdateCompany)
	e.DELETE("/companies/:id", controllers.DeleteCompany)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
