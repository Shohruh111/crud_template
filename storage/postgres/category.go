package postgres

import (
	"app/models"
	"app/pkg/helper"
	"database/sql"
	"errors"
	"fmt"

	uuid "github.com/google/uuid"
)

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *categoryRepo {
	return &categoryRepo{
		db: db,
	}
}

func (r *categoryRepo) Create(req *models.CreateCategory) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO category(id, title, parent_id, updated_at)
		VALUES ($1, $2, $3, NOW())
	`

	_, err := r.db.Exec(query,
		id,
		req.Title,
		helper.NewNullString(req.ParentID),
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *categoryRepo) GetByID(req *models.CategoryPrimaryKey) (*models.Category, error) {

	var (
		resp  models.Category
		query string
	)

	query = `
		SELECT
			id,
			title,
			COALESCE(parent_id::VARCHAR, ''),
			created_at,
			updated_at
		FROM category
		WHERE id = $1
	`

	err := r.db.QueryRow(query, req.Id).Scan(
		&resp.Id,
		&resp.Title,
		&resp.ParentID,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r *categoryRepo) GetList(req *models.CategoryGetListRequest) (*models.CategoryGetListResponse, error) {

	var (
		resp   = &models.CategoryGetListResponse{}
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
			parent_id,
			created_at,
			updated_at
		FROM category
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
			category models.Category
			parentId sql.NullString
		)
		err := rows.Scan(
			&resp.Count,
			&category.Id,
			&category.Title,
			&parentId,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		category.ParentID = parentId.String
		resp.Categories = append(resp.Categories, &category)
	}

	return resp, nil
}
func (r *categoryRepo) Update(req *models.UpdateCategory) (*models.Category, error) {

	var (
		query string
		resp  = &models.Category{}
	)
	query = `
		update category() 
		set 
			title=$1	
			parent_id=$2
			updated_at=NOW()
		where id=$3
	`
	_, err := r.db.Query(query, req.Title, req.ParentID, req.Id)
	if err != nil {
		return nil, err
	}
	resp, err = r.GetByID(&models.CategoryPrimaryKey{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *categoryRepo) Delete(req *models.CategoryPrimaryKey) error {
	var (
		query string
	)
	query = `
		DELETE FROM CATEGORY WHERE ID=$1
	`
	_, err := r.db.Query(query, req.Id)
	if err != nil {
		return err
	}

	return errors.New("Successfully deleted!")
}
