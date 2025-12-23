
package controllers

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/labstack/echo/v4"
	"myjob/models"
	"myjob/config"
)

// RegisterInput struct
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // client se milega
}

// Register handler
func Register(c echo.Context) error {
	input := new(RegisterInput)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	db := config.ConnectDB() // tumhara DB connection

	// Check if email already exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", input.Email).Scan(&exists)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "DB error"})
	}
	if exists {
		return c.JSON(http.StatusConflict, map[string]string{"message": "Email already registered"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error hashing password"})
	}

	// Map to User model
	user := &models.User{
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Role:         input.Role,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Insert user into DB
	_, err = db.Exec(
		`INSERT INTO users (email, password_hash, role, is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		user.Email, user.PasswordHash, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error saving user"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "registered successfully"})
}
