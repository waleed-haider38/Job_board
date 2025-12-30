package models

type JobSeekerSkill struct {
	JobSeekerID int `gorm:"column:job_seeker_id;primaryKey" json:"job_seeker_id"`
	SkillID     int `gorm:"column:skill_id;primaryKey" json:"skill_id"`

	// 🔗 Relations
	JobSeeker JobSeeker `gorm:"foreignKey:JobSeekerID;references:JobSeekerID"` // ensure references correct PK
	Skill     Skill     `gorm:"foreignKey:SkillID;references:SkillID"`
}

func (JobSeekerSkill) TableName() string {
	return "job_seeker_skills"
}
