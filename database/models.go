package database

import "time"

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	StockQty    int       `json:"stock_qty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Password  string    `json:"password"`
	Role      string    `json:"role" gorm:"default:user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Cart      uint
}

type Order struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserId    uint      `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Cart      uint
}

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderId   uint    `json:"order_id"`
	ProductId uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Cart struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserId    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	CartItems []CartItem
}

type CartItem struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	CartId    uint `json:"cart_id"`
	ProductId uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}
