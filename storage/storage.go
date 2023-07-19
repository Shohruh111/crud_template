package storage

import "app/models"

type StorageI interface {
	Close()
	Category() CategoryRepoI
	Product() ProductRepoI
}

type CategoryRepoI interface {
	Create(*models.CreateCategory) (string, error)
	GetByID(*models.CategoryPrimaryKey) (*models.Category, error)
	GetList(*models.CategoryGetListRequest) (*models.CategoryGetListResponse, error)
	Update(*models.UpdateCategory) (*models.Category, error)
	Delete(*models.CategoryPrimaryKey) error
}

type ProductRepoI interface {
	CreateProduct(*models.CreateProduct) (string, error)
	GetByIDProduct(*models.ProductPrimaryKey) (*models.Product, error)
	GetListProduct(*models.ProductGetListRequest) (*models.ProductGetListResponse, error)
	UpdateProduct(*models.UpdateProduct) (*models.Product, error)
	DeleteProduct(*models.ProductPrimaryKey) error
}
