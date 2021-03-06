package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/adrianolmedo/go-restapi/internal/domain"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r ProductRepository) Create(product *domain.Product) error {
	stmt, err := r.db.Prepare("INSERT INTO products (name, observations, price, created_at) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return err
	}
	defer stmt.Close()

	product.CreatedAt = time.Now()

	err = stmt.QueryRow(product.Name, product.Observations, product.Price, product.CreatedAt).Scan(&product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r ProductRepository) ByID(id int64) (*domain.Product, error) {
	stmt, err := r.db.Prepare("SELECT * FROM products WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	product, err := scanRowProduct(stmt.QueryRow(id))
	if errors.Is(err, sql.ErrNoRows) {
		return &domain.Product{}, domain.ErrProductNotFound
	}

	if err != nil {
		return &domain.Product{}, err
	}

	return product, nil
}

func (r ProductRepository) Update(product domain.Product) error {
	stmt, err := r.db.Prepare("UPDATE products SET name = $1, observations = $2, price = $3, updated_at = $4 WHERE id = $5")
	if err != nil {
		return err
	}
	defer stmt.Close()

	product.UpdatedAt = time.Now()

	result, err := stmt.Exec(product.Name, product.Observations, product.Price, timeToNull(product.UpdatedAt), product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func (r ProductRepository) All() ([]*domain.Product, error) {
	stmt, err := r.db.Prepare("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*domain.Product, 0)

	for rows.Next() {
		p, err := scanRowProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r ProductRepository) Delete(id int64) error {
	stmt, err := r.db.Prepare("DELETE FROM products WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r ProductRepository) DeleteAll() error {
	stmt, err := r.db.Prepare("TRUNCATE TABLE products RESTART IDENTITY CASCADE")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("can't truncate table: %v", err)
	}

	return nil
}

// scanRowUser return nulled fields of User parsed.
func scanRowProduct(s scanner) (*domain.Product, error) {
	var updatedAtNull sql.NullTime
	p := &domain.Product{}

	err := s.Scan(
		&p.ID,
		&p.Name,
		&p.Observations,
		&p.Price,
		&p.CreatedAt,
		&updatedAtNull,
	)
	if err != nil {
		return &domain.Product{}, err
	}

	p.UpdatedAt = updatedAtNull.Time

	return p, nil
}
