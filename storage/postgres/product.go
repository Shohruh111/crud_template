package postgres

import (
	"app/models"
	"database/sql"
	"fmt"

	uuid "github.com/google/uuid"
)

type productRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *productRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) CreateProduct(req *models.CreateProduct) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO product(id, title,price, category, updated_at)
		VALUES ($1, $2, $3,$4, NOW())
	`

	_, err := r.db.Exec(query,
		id,
		req.Title,
		req.Price,
		req.Category,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *productRepo) GetByIDProduct(req *models.ProductPrimaryKey) (*models.Product, error) {

	var (
		resp  models.Product
		query string
	)

	query = `
		SELECT
			id,
			title,
			price,
			COALESCE(category::VARCHAR, ''),
			created_at,
			updated_at
		FROM product
		WHERE id = $1
	`

	err := r.db.QueryRow(query, req.Id).Scan(
		&resp.Id,
		&resp.Title,
		&resp.Price,
		&resp.Category,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r *productRepo) GetListProduct(req *models.ProductGetListRequest) (*models.ProductGetListResponse, error) {

	var (
		resp   = &models.ProductGetListResponse{}
		query  string
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			title,
			category,
			created_at,
			updated_at
		FROM product
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND title ILIKE '%' || '` + req.Search + `' || '%'`
	}

	query += where + offset + limit

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			product  models.Product
			category sql.NullString
		)
		err := rows.Scan(
			&resp.Count,
			&product.Id,
			&product.Title,
			&product.Category,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		product.Category = category.String
		resp.Products = append(resp.Products, &product)
	}

	return resp, nil
}
func (r *productRepo) UpdateProduct(req *models.UpdateProduct) (*models.Product, error) {

	var (
		query string
		resp  = &models.Product{}
	)
	query = `
		update product 
		set 
			title=$1,
			price=$2,
			category=$3,
			updated_at=NOW()
		where id=$4
	`
	_, err := r.db.Query(query, req.Title, req.Price, req.Category, req.Id)
	if err != nil {
		return nil, err
	}
	resp, err = r.GetByIDProduct(&models.ProductPrimaryKey{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *productRepo) DeleteProduct(req *models.ProductPrimaryKey) error {
	var (
		query string
	)
	query = `
		DELETE FROM PRODUCT WHERE ID=$1
	`
	_, err := r.db.Query(query, req.Id)
	if err != nil {
		return err
	}

	return nil
}
