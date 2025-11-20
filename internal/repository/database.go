package repository

import (
	"context"
	"log"

	"order-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ OrderRepository = (*DB)(nil)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(connectionString string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 10
	cfg.MinConns = 1

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &DB{pool: pool}, nil
}

func (p *DB) Close() {
	p.pool.Close()
}

// 🔥 ИСПРАВЛЕНО: Добавлен ctx как параметр
func (p *DB) SaveOrder(ctx context.Context, order *models.Order) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Order
	_, err = tx.Exec(ctx,
		`INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, 
		 customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}

	// Delivery
	_, err = tx.Exec(ctx,
		`INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return err
	}

	// Payment
	_, err = tx.Exec(ctx,
		`INSERT INTO payments (order_uid, transaction, request_id, currency, provider, 
		 amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return err
	}

	// Items
	for _, item := range order.Items {
		_, err = tx.Exec(ctx,
			`INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, 
			 total_price, nm_id, brand, status) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// 🔥 ИСПРАВЛЕНО: Добавлен ctx как параметр
func (p *DB) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	query := `
		SELECT 
			o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
			o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
			d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
			p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt,
			p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders o
		LEFT JOIN deliveries d ON o.order_uid = d.order_uid
		LEFT JOIN payments p ON o.order_uid = p.order_uid
		WHERE o.order_uid = $1
	`

	row := p.pool.QueryRow(ctx, query, orderUID)

	var order models.Order
	var delivery models.Delivery
	var payment models.Payment

	err := row.Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
		&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email,
		&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDt,
		&payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
	)
	if err != nil {
		return nil, err
	}

	order.Delivery = delivery
	order.Payment = payment

	// Items
	rows, err := p.pool.Query(ctx,
		`SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status 
		 FROM items WHERE order_uid = $1`, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	order.Items = items
	return &order, nil
}

// 🔥 ИСПРАВЛЕНО: Добавлен ctx как параметр
func (p *DB) GetAllOrders(ctx context.Context) (map[string]*models.Order, error) {
	rows, err := p.pool.Query(ctx, `SELECT order_uid FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make(map[string]*models.Order)
	for rows.Next() {
		var orderUID string
		if err := rows.Scan(&orderUID); err != nil {
			continue
		}

		order, err := p.GetOrder(ctx, orderUID) // 🔥 Передаем ctx
		if err != nil {
			log.Printf("Ошибка загрузки заказа %s: %v", orderUID, err)
			continue
		}

		orders[orderUID] = order
	}

	return orders, nil
}

// 🔥 ИСПРАВЛЕНО: Используем переданный ctx
func (p *DB) HealthCheck(ctx context.Context) error {
	return p.pool.Ping(ctx)
}
