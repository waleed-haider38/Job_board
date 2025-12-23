package models

import "time"

type User struct {
    UserID      int       `db:"user_id" json:"user_id"`
    Email       string    `db:"email" json:"email"`
    PasswordHash string   `db:"password_hash" json:"-"` // json:"-" ka matlab hai ye response mein nahi jayega
    Role        string    `db:"role" json:"role"`       // "job_seeker" ya "employer"
    IsActive    bool      `db:"is_active" json:"is_active"`
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
    UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
