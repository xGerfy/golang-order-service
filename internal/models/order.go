package models

import "time"

type Order struct {
	OrderUID          string    `json:"order_uid" db:"order_uid" validate:"required,min=1"`
	TrackNumber       string    `json:"track_number" db:"track_number" validate:"required,min=1"`
	Entry             string    `json:"entry" db:"entry" validate:"required,min=1"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required,min=1,dive"`
	Locale            string    `json:"locale" db:"locale" validate:"required,min=1"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	CustomerID        string    `json:"customer_id" db:"customer_id" validate:"required,min=1"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service" validate:"required,min=1"`
	Shardkey          string    `json:"shardkey" db:"shardkey" validate:"required,min=1"`
	SmID              int       `json:"sm_id" db:"sm_id" validate:"required,min=1"`
	DateCreated       time.Time `json:"date_created" db:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" db:"oof_shard" validate:"required,min=1"`
}

type Delivery struct {
	OrderUID string `json:"order_uid" db:"order_uid"`
	Name     string `json:"name" db:"name" validate:"required,min=1"`
	Phone    string `json:"phone" db:"phone" validate:"required,min=1"`
	Zip      string `json:"zip" db:"zip" validate:"required,min=1"`
	City     string `json:"city" db:"city" validate:"required,min=1"`
	Address  string `json:"address" db:"address" validate:"required,min=1"`
	Region   string `json:"region" db:"region" validate:"required,min=1"`
	Email    string `json:"email" db:"email" validate:"required,email"`
}

type Payment struct {
	OrderUID     string `json:"order_uid" db:"order_uid"`
	Transaction  string `json:"transaction" db:"transaction" validate:"required,min=1"`
	RequestID    string `json:"request_id" db:"request_id"`
	Currency     string `json:"currency" db:"currency" validate:"required,min=1"`
	Provider     string `json:"provider" db:"provider" validate:"required,min=1"`
	Amount       int    `json:"amount" db:"amount" validate:"required,min=0"`
	PaymentDt    int64  `json:"payment_dt" db:"payment_dt" validate:"required,min=0"`
	Bank         string `json:"bank" db:"bank" validate:"required,min=1"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost" validate:"min=0"`
	GoodsTotal   int    `json:"goods_total" db:"goods_total" validate:"required,min=0"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee" validate:"min=0"`
}

type Item struct {
	OrderUID    string `json:"order_uid" db:"order_uid"`
	ChrtID      int    `json:"chrt_id" db:"chrt_id" validate:"required,min=1"`
	TrackNumber string `json:"track_number" db:"track_number" validate:"required,min=0"`
	Price       int    `json:"price" db:"price" validate:"required,min=1"`
	Rid         string `json:"rid" db:"rid" validate:"required,min=1"`
	Name        string `json:"name" db:"name" validate:"required,min=1"`
	Sale        int    `json:"sale" db:"sale" validate:"min=0"`
	Size        string `json:"size" db:"size" validate:"required,min=1"`
	TotalPrice  int    `json:"total_price" db:"total_price" validate:"required,min=0"`
	NmID        int    `json:"nm_id" db:"nm_id" validate:"required,min=1"`
	Brand       string `json:"brand" db:"brand" validate:"required,min=1"`
	Status      int    `json:"status" db:"status" validate:"required,min=0"`
}
