package authdto

type RegisterRequest struct {
	Name     string `gorm:"type: varchar(255)" json:"name" validate:"required"`
	Email    string `gorm:"type: varchar(255)" json:"email" validate:"required"`
	Password string `gorm:"type: varchar(255)" json:"password" validate:"required"`
}
