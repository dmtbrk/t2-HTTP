package market

import (
	"errors"
	"fmt"
)

//go:generate mockgen -destination=./mock/product_service.go  -package=mock . ProductService

var ErrProductNotFound = errors.New("product not found")

// ProductService represents a product data backend.
type ProductService interface {
	Products() ([]*Product, error)
	Product(int) (*Product, error)
	AddProduct(*Product) (*Product, error)
	ReplaceProduct(*Product) (*Product, error)
	DeleteProduct(int) error
}

type Product struct {
	ID     int
	Name   string
	Price  int
	Seller string
}

func (p *Product) String() string {
	return fmt.Sprintf("Product{ ID: %d, Name: %s, Price: %d, Seller: %v }", p.ID, p.Name, p.Price, p.Seller)
}
