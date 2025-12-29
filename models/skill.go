package models

type Skill struct {
	SkillID   int    `gorm:"column:skill_id;primaryKey" json:"skill_id"`
	SkillName string `gorm:"column:skill_name" json:"skill_name"`
}

func (Skill) TableName() string {
	return "skills"
}
