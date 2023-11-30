package models

type User struct {
	Email     string `gorm:"column:email" json:"email"`
	FirstName string `gorm:"column:firstname" json:"firstName"`
	LastName  string `gorm:"column:lastname" json:"lastName"`
	Password  string `gorm:"column:password" json:"password"`
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
