package models

import "time"

type Job struct {
	ID        int `gorm:"column:job_id;primaryKey" json:"job_id"`
	CompanyID int `gorm:"column:company_id" json:"company_id"`

	Title       string `gorm:"column:title" json:"title"`
	Description string `gorm:"column:description" json:"description"`
	JobType     string `gorm:"column:job_type" json:"job_type"`
	JobLocation string `gorm:"column:job_location" json:"job_location"`
	SalaryMin   int    `gorm:"column:salary_min" json:"salary_min"`
	SalaryMax   int    `gorm:"column:salary_max" json:"salary_max"`
	Status      string `gorm:"column:status" json:"status"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	Company Company `gorm:"foreignKey:CompanyID;references:CompanyID" json:"company"`

	//  NOW GORM WILL USE: job_id + skill_id
	Skills []Skill `gorm:"many2many:job_skills" json:"skills"`
}

func (Job) TableName() string {
	return "jobs"
}
