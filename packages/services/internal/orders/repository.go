package orders

import (
	"context"

	"gorm.io/gorm"
)

// Repository provides persistence for orders.
type Repository struct {
	db *gorm.DB
}

// NewRepository constructs the order repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create persists a new order record.
func (r *Repository) Create(ctx context.Context, order *Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// UpdateStatus updates order status by primary key.
func (r *Repository) UpdateStatus(ctx context.Context, id uint, status OrderStatus) error {
	return r.db.WithContext(ctx).Model(&Order{}).Where("id = ?", id).Update("status", status).Error
}

// ListActive returns orders matching filter criteria.
func (r *Repository) ListActive(ctx context.Context, side string, collection string) ([]Order, error) {
	var result []Order
	query := r.db.WithContext(ctx).Where("status = ?", OrderStatusActive)
	if side != "" {
		query = query.Where("side = ?", side)
	}
	if collection != "" {
		query = query.Where("nft_address = ?", collection)
	}
	if err := query.Order("created_at DESC").Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// ListByStatus 根据状态查询订单（支持 filled, cancelled 等）
func (r *Repository) ListByStatus(ctx context.Context, status string, side string, collection string) ([]Order, error) {
	var result []Order
	query := r.db.WithContext(ctx).Where("status = ?", status)
	if side != "" {
		query = query.Where("side = ?", side)
	}
	if collection != "" {
		query = query.Where("nft_address = ?", collection)
	}
	if err := query.Order("updated_at DESC").Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// FindByID retrieves an order by primary key.
func (r *Repository) FindByID(ctx context.Context, id uint) (*Order, error) {
	var ord Order
	if err := r.db.WithContext(ctx).First(&ord, id).Error; err != nil {
		return nil, err
	}
	return &ord, nil
}
