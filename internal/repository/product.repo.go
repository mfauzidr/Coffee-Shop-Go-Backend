package repository

import (
	"database/sql"
	"fmt"

	"github.com/mfauzidr/coffeeshop-go-backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type RepoProduct struct {
	*sqlx.DB
}

func NewProduct(db *sqlx.DB) *RepoProduct {
	return &RepoProduct{db}
}

func (r *RepoProduct) CreateProduct(data *models.Product) error {
	query := `INSERT INTO public.product(
		name,
		price,
		category,
		description,
		stock
	) VALUES(
	 	:name,
		:price,
		:category,
		:description,
		:stock
	)`

	_, err := r.NamedExec(query, data)
	return err
}

func (r *RepoProduct) GetAllProduct(que *models.ProductQuery) (*models.Products, error) {
	query := `SELECT * FROM public.product`
	var values []interface{}
	condition := false

	if que.Name != "" {
		query += fmt.Sprintf(` WHERE name ILIKE $%d`, len(values)+1)
		values = append(values, "%"+que.Name+"%")
		condition = true
	}
	if que.MinPrice != 0 {
		if condition {
			query += " AND "
		} else {
			query += " WHERE "
		}
		query += fmt.Sprintf(` price > $%d`, len(values)+1)
		values = append(values, que.MinPrice)
		condition = true
	}
	if que.MaxPrice != 0 {
		if condition {
			query += " AND "
		} else {
			query += " WHERE "
		}
		query += fmt.Sprintf(` price < $%d`, len(values)+1)
		values = append(values, que.MaxPrice)
		condition = true
	}
	if que.Category != "" {
		if condition {
			query += " AND "
		} else {
			query += " WHERE "
		}
		query += fmt.Sprintf(` category = $%d`, len(values)+1)
		values = append(values, que.Category)
		condition = true
	}

	switch que.Sort {
	case "alphabet":
		query += " ORDER BY name ASC"
	case "price":
		query += " ORDER BY price ASC"
	case "asc":
		query += " ORDER BY created_at DESC"
	case "desc":
		query += " ORDER BY created_at ASC"
	}

	if que.Page > 0 {
		limit := 5
		offset := (que.Page - 1) * limit
		query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(values)+1, len(values)+2)
		values = append(values, limit, offset)
	}

	rows, err := r.DB.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data models.Products
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.UUID,
			&product.Name,
			&product.Price,
			&product.Category,
			&product.Description,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *RepoProduct) GetDetailProduct(id int) (*models.Product, error) {
	query := `SELECT * FROM public.product WHERE product_id = :product_id`
	data := models.Product{}

	rows, err := r.DB.NamedQuery(query, map[string]interface{}{
		"product_id": id,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&data)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}

	return nil, nil
}
func (r *RepoProduct) DeleteProduct(id int) error {
	query := `DELETE FROM public.product WHERE product_id = :product_id`

	_, err := r.DB.NamedExec(query, map[string]interface{}{
		"product_id": id,
	})
	return err
}

func (r *RepoProduct) UpdateProduct(data *models.Product, id int) (*models.Product, error) {
	query := `UPDATE public.product SET`
	var values []interface{}
	condition := false

	if data.Name != "" {
		query += fmt.Sprintf(` name = $%d`, len(values)+1)
		values = append(values, data.Name)
		condition = true
	}
	if data.Price != 0 {
		if condition {
			query += ","
		}
		query += fmt.Sprintf(` price = $%d`, len(values)+1)
		values = append(values, data.Price)
		condition = true
	}
	if data.Category != "" {
		if condition {
			query += ","
		}
		query += fmt.Sprintf(` category = $%d`, len(values)+1)
		values = append(values, data.Category)
		condition = true
	}
	if data.Description != "" {
		if condition {
			query += ","
		}
		query += fmt.Sprintf(` description = $%d`, len(values)+1)
		values = append(values, data.Description)
		condition = true
	}

	if !condition {
		return nil, fmt.Errorf("no fields to update")
	}

	query += fmt.Sprintf(`, updated_at = now() WHERE product_id = $%d RETURNING *`, len(values)+1)
	values = append(values, id)

	row := r.DB.QueryRow(query, values...)
	var product models.Product
	err := row.Scan(
		&product.ID,
		&product.UUID,
		&product.Name,
		&product.Price,
		&product.Category,
		&product.Description,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(`product with id %d not found`, id)
		}
		return nil, fmt.Errorf(`query execution error: %w`, err)
	}

	return &product, nil
}
