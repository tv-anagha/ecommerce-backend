package model

import "time"

// Order is a persisted checkout for one user.
type Order struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	UserID      uint        `json:"userId" gorm:"column:user_id;not null"`
	TotalAmount float64     `json:"totalAmount" gorm:"column:total_amount;not null"`
	CreatedAt   time.Time   `json:"createdAt"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

// OrderItem snapshots product name and price at checkout time.
type OrderItem struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	OrderID     uint    `json:"orderId" gorm:"column:order_id;not null"`
	ProductID   uint    `json:"productId" gorm:"column:product_id;not null"`
	ProductName string  `json:"productName" gorm:"column:product_name;not null"`
	Price       float64 `json:"price" gorm:"not null"`
	Quantity    int     `json:"quantity" gorm:"not null"`
}
