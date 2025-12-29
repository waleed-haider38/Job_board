package models

import "time"

type Application struct {
	ApplicationID int       `gorm:"column:application_id;primaryKey" json:"application_id"`
	JobID         int       `gorm:"column:job_id" json:"job_id"`
	JobSeekerID   int       `gorm:"column:job_seeker_id" json:"job_seeker_id"`

	Status      string    `gorm:"column:status" json:"status"`
	CoverLetter string    `gorm:"column:cover_letter" json:"cover_letter"`
	AppliedAt   time.Time `gorm:"column:applied_at" json:"applied_at"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	// 🔗 Relations
	Job       Job       `gorm:"foreignKey:JobID"`
	JobSeeker JobSeeker `gorm:"foreignKey:JobSeekerID"`
}

func (Application) TableName() string {
	return "applications"
}
