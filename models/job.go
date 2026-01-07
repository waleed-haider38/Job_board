package models

import "time"

type Job struct {
    JobID      int       `gorm:"column:job_id;primaryKey" json:"job_id"`
    EmployerID int       `gorm:"column:employer_id" json:"employer_id"`

    Title       string    `gorm:"column:title" json:"title"`
    Description string    `gorm:"column:description" json:"description"`
    JobType     string    `gorm:"column:job_type" json:"job_type"`
    JobLocation string    `gorm:"column:job_location" json:"job_location"`
    SalaryMin   int       `gorm:"column:salary_min" json:"salary_min"`
    SalaryMax   int       `gorm:"column:salary_max" json:"salary_max"`

    CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
    UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

    // Relations
    Employer User    `gorm:"foreignKey:EmployerID"`
    Skills   []Skill `gorm:"many2many:job_skills;joinForeignKey:JobID;JoinReferences:SkillID" json:"skills"`
}

func (Job) TableName() string {
    return "jobs"
}
