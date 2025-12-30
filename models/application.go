package models

import "time"

type Application struct {
	ApplicationID int `gorm:"column:application_id;primaryKey" json:"application_id"`

	JobID       int `gorm:"column:job_id;not null" json:"job_id"`
	JobSeekerID int `gorm:"column:job_seeker_id;not null" json:"job_seeker_id"`

	CoverLetter string `gorm:"column:cover_letter" json:"cover_letter"`
	Status      string `gorm:"column:status;default:applied" json:"status"`

	AppliedAt time.Time `gorm:"column:applied_at;autoCreateTime" json:"applied_at"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// 🔗 Relations
	Job       Job       `gorm:"foreignKey:JobID"`
	JobSeeker JobSeeker `gorm:"foreignKey:JobSeekerID"`
}

func (Application) TableName() string {
	return "applications"
}
