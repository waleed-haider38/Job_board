package models

type Company struct {
	CompanyID      int    `gorm:"column:company_id;primaryKey" json:"company_id"`
	EmployerID     int    `gorm:"column:employer_id" json:"employer_id"`
	CompanyName    string `gorm:"column:company_name" json:"company_name"`
	CompanyProduct string `gorm:"column:company_product" json:"company_product"`
}

func (Company) TableName() string {
	return "companies"
}
