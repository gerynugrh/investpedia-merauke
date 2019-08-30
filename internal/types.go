package internal

import (
	"database/sql"
)

// Args used for this application
type Args struct {

	// Port used by this service
	Port int
}

// Handler object used to handle the HTTP API
type Handler struct {

	// DB object that'll be used
	DB *sql.DB
}

// User struct for database query
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Product struct {
	ID int `json:"id"`
	ProductName string `json:"product_name"`
	Price string `json:"price"`
}
type Investment struct {
	ID int `json:"id"`
	RoomId string `json:"room_id"`
	ProductId int `json:"product_id"`
	TotalPayment string `json:"total_payment"`
}

