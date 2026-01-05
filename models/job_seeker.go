package models

import "time"

type JobSeeker struct {
	JobSeekerID int `gorm:"column:job_seeker_id;primaryKey" json:"job_seeker_id"`
	//User and UserID are not the same because it is the column in our job_seeker table and User defines the relation as we can see foreign key:UserID.
	UserID      int `gorm:"column:user_id" json:"user_id"`

	FullName  string `gorm:"column:full_name" json:"full_name"`
	ResumeURL string `gorm:"column:resume_url" json:"resume_url"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	// 🔗 Relation
	User User `gorm:"foreignKey:UserID"`
	Skills []Skill `gorm:"many2many:job_seeker_skills;joinForeignKey:JobSeekerID;JoinReferences:SkillID" json:"skills"`

}

func (JobSeeker) TableName() string {
	return "job_seekers"
}
