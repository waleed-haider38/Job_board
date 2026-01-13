package models

import "time"

type Application struct {
	ApplicationID int `gorm:"column:application_id;primaryKey" json:"application_id"`

	JobID       int `gorm:"column:job_id;not null" json:"job_id"`
	JobSeekerID int `gorm:"column:job_seeker_id;not null" json:"job_seeker_id"`

	CoverLetter string `gorm:"column:cover_letter" json:"cover_letter"`
	Status      string `gorm:"column:status;default:pending" json:"status"`

	AppliedAt time.Time `gorm:"column:applied_at;autoCreateTime" json:"applied_at"`

	//  Relations
	Job       Job       `gorm:"foreignKey:JobID;references:JobID"`
	JobSeeker JobSeeker `gorm:"foreignKey:JobSeekerID;references:JobSeekerID"`
}

func (Application) TableName() string {
	return "applications"
}
