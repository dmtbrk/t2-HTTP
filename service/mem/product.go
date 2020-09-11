package mem

import (
	"sync"

	"github.com/ortymid/t2-http/market"
)

type ProductService struct {
	mu       sync.RWMutex
	lastID   int
	products []*market.Product
}

func NewProductService() *ProductService {
	products := []*market.Product{
		{ID: 1, Name: "Banana", Price: 1500, Seller: "1"},
		{ID: 2, Name: "Carrot", Price: 1400, Seller: "2"},
	}
	return &ProductService{products: products, lastID: 2}
}

func (srv *ProductService) Products() ([]*market.Product, error) {
	srv.mu.RLock()
	defer srv.mu.RUnlock()

	return srv.products, nil
}

func (srv *ProductService) Product(id int) (*market.Product, error) {
	srv.mu.RLock()
	defer srv.mu.RUnlock()

	for _, p := range srv.products {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, market.ErrProductNotFound
}

func (srv *ProductService) AddProduct(p *market.Product) (*market.Product, error) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	srv.lastID++
	p.ID = srv.lastID
	srv.products = append(srv.products, p)
	return p, nil
}

func (srv *ProductService) ReplaceProduct(np *market.Product) (*market.Product, error) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	for i, op := range srv.products {
		if op.ID == np.ID {
			srv.products[i] = np
			return np, nil
		}
	}
	return nil, market.ErrProductNotFound
}

func (srv *ProductService) DeleteProduct(id int) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	for i, p := range srv.products {
		if p.ID == id {
			if i == len(srv.products)-1 {
				srv.products[i] = nil
				srv.products = srv.products[:i]
			}
			copy(srv.products[i:], srv.products[i+1:])
			return nil
		}
	}
	return market.ErrProductNotFound
}
