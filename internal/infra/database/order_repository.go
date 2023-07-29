package database

import (
	"database/sql"

	"github.com/twsm000/goexp-cleanarch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) ListAll() ([]*entity.Order, error) {
	rows, err := r.Db.Query(`SELECT id, price, tax, final_price FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*entity.Order{}
	for rows.Next() {
		var order entity.Order
		err = rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
		if err != nil {
			return nil, err
		}
		result = append(result, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
