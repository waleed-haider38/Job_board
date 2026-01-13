package models

type Employer struct {
	EmployerID    int    `gorm:"column:employer_id;primaryKey" json:"employer_id"`
	UserID        int    `gorm:"column:user_id;not null" json:"user_id"`
	EmployerName  string `gorm:"column:employer_name" json:"employer_name"`
	EmployerEmail string `gorm:"column:employer_email" json:"employer_email"`

	//  Relation
	User User `gorm:"foreignKey:UserID;references:UserID"`
}

func (Employer) TableName() string {
	return "employers"
}
