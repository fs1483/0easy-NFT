package orders

import "time"

// OrderStatus represents the lifecycle state of an order.
type OrderStatus string

const (
	OrderStatusActive    OrderStatus = "active"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusFilled    OrderStatus = "filled"
)

// Order models a signed order stored off-chain.
type Order struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	Maker        string      `gorm:"type:varchar(66);index:idx_orders_maker_nonce,unique" json:"maker"`
	NFTAddress   string      `gorm:"type:varchar(66);index;column:nft_address" json:"nftAddress"`
	TokenID      string      `gorm:"type:numeric;column:token_id" json:"tokenId"`
	PaymentToken string      `gorm:"type:varchar(66);column:payment_token" json:"paymentToken"`
	Price        string      `gorm:"type:numeric" json:"price"`
	Expiry       time.Time   `gorm:"index" json:"expiry"`
	Nonce        string      `gorm:"type:numeric;index:idx_orders_maker_nonce,unique" json:"nonce"`
	Side         string      `gorm:"type:varchar(4);index" json:"side"`
	Status       OrderStatus `gorm:"type:varchar(16);index" json:"status"`
	Signature    string      `gorm:"type:varchar(132)" json:"signature"`
	Hash         string      `gorm:"type:varchar(66);index" json:"hash"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}

// TableName overrides default table name.
func (Order) TableName() string {
	return "orders"
}
