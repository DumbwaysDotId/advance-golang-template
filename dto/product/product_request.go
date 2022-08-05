package productdto

type ProductRequest struct {
	Name       string `json:"name" form:"name" gorm:"type: varchar(255)"`
	Desc       string `json:"desc" gorm:"type:text" form:"desc"`
	Price      int    `json:"price" form:"price" gorm:"type: int"`
	Image      string `json:"image" form:"image" gorm:"type: varchar(255)"`
	Qty        int    `json:"qty" form:"qty" gorm:"type: int"`
	UserID     int    `json:"user_id" gorm:"type: int"`
	CategoryID int    `json:"category_id" form:"category_id" gorm:"type: int"`
}
