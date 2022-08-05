package authdto

type AuthRequest struct {
	ID       int    `json:"id"`
	Name     string `gorm:"type: varchar(255)" json:"name"`
	Email    string `gorm:"type: varchar(255)" json:"email"`
	Password string `gorm:"type: varchar(255)" json:"password"`
}
