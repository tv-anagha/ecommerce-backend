package model

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name" gorm:"column:product_name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	ImageURL string  `json:"image_url" gorm:"column:image_url"`
	Quantity int     `json:"quantity" gorm:"column:quantity"`
}

func (Product) TableName() string {
	return "products"
}
