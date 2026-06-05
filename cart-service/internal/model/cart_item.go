package model

// CartItem is one row in a user's shopping cart (one row per user + product).
type CartItem struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	UserID    uint `json:"userId" gorm:"column:user_id;not null;uniqueIndex:idx_user_product"`
	ProductID uint `json:"productId" gorm:"column:product_id;not null;uniqueIndex:idx_user_product"`
	Quantity  int  `json:"quantity" gorm:"not null"`
}

func NewCartItem(userID, productID uint, quantity int) CartItem {
	return CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
}
