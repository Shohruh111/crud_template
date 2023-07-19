package models

type Product struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Price     int `json:"price"`
	Category  string `json:"category"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateProduct struct {
	Title    string `json:"title"`
	Price    int `json:"price"`
	Category string `json:"category"`
}

type UpdateProduct struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Price     int `json:"price"`
	Category  string `json:"category"`
	UpdatedAt string `json:"updated_at"`
}
type ProductPrimaryKey struct {
	Id string `json:"id"`
}

type ProductGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type ProductGetListResponse struct {
	Count    int        `json:"count"`
	Products []*Product `json:"categories"`
}
