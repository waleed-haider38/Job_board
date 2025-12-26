
package controllers

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/labstack/echo/v4"
	"github.com/golang-jwt/jwt/v5"
	"myjob/models"
	"myjob/config"
)

// RegisterInput struct
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // client se milega
}



// LoginInput struct
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JWT Claims (team standard)
type JwtCustomClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}


// Bind binds the path params , query params and the request body into provided type.
//The new built-in function allocates memory. The first argument is a type, not a value, and the value returned is a pointer to a newly allocated zero value of that type.


func Login(c echo.Context) error {
	input := new(LoginInput)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid input",
		})
	}

	db := config.ConnectDB()

	// Get user from DB
	user := models.User{}

	//QueryRow is used to return single row of data from the database.
	err := db.QueryRow(
		`SELECT user_id, email, password_hash, role 
		 FROM users 
		 WHERE email = $1 AND is_active = true`,
		input.Email,

	//Scan copies the columns from the matched row into the values.
	).Scan(&user.UserID, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "Invalid email or password",
		})
	}

	// Compare password
	//CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext equivalent. Returns nil on success, or an error on failure.
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "Invalid email or password",
		})
	}

	// Create JWT claims
	claims := &JwtCustomClaims{
		UserID: user.UserID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	// Generate token
	//NewWithClaims creates a new Token with the specified signing method and claims.
	//SignedString creates and returns a complete, signed JWT. The token is signed using the SigningMethod specified in the token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Could not generate token",
		})
	}

	// Success response
	return c.JSON(http.StatusOK, echo.Map{
		"token": tokenString,
		"user": echo.Map{
			"id":    user.UserID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
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
