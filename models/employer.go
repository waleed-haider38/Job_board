package models

import "time"

type Employer struct {
	EmployerID    int       `gorm:"column:employer_id;primaryKey" json:"employer_id"`
	UserID        int       `gorm:"column:user_id" json:"user_id"`
	EmployerName  string    `gorm:"column:employer_name" json:"employer_name"`
	EmployerEmail string    `gorm:"column:employer_email" json:"employer_email"`

	// we can use gorm.Models but it has id column as primary key by default init. but it is not neccessary for now.
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	// 🔗 Relation
	User User `gorm:"foreignKey:UserID"`
}

func (Employer) TableName() string {
	return "employers"
}
