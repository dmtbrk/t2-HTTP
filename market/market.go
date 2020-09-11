package market

import (
	"errors"
	"fmt"
)

// ErrPermission is an error returned when a user does not have rights
// to do Market methods.
type ErrPermission struct {
	Reason error
}

func (err *ErrPermission) Error() string {
	return fmt.Sprintf("permission denied: %s", err.Reason)
}

func (err *ErrPermission) Is(target error) bool {
	_, ok := target.(*ErrPermission)
	return ok
}

func (err *ErrPermission) Unwrap() error {
	return err.Reason
}

// Interface may be used by protocol layers for RPC or mocking.
type Interface interface {
	Products() ([]*Product, error)
	Product(id int) (*Product, error)
	AddProduct(p *Product, userID string) (*Product, error)
	ReplaceProduct(p *Product, userID string) (*Product, error)
	DeleteProduct(id int, userID string) error
}

// Market composes business logic from different services.
type Market struct {
	AuthService    AuthService
	UserService    UserService
	ProductService ProductService
}

// Products returns all products on the market.
func (m *Market) Products() ([]*Product, error) {
	ps, err := m.ProductService.Products()
	if err != nil {
		err = fmt.Errorf("products: %w", err)
		return nil, err
	}
	return ps, nil
}

// Product finds the product by its ID.
func (m *Market) Product(id int) (*Product, error) {
	p, err := m.ProductService.Product(id)
	if err != nil {
		err = fmt.Errorf("product: %w", err)
		return nil, err
	}
	return p, nil
}

func (m *Market) AddProduct(p *Product, userID string) (*Product, error) {
	// Check the user for permission. Only the existence of the user counts yet.
	_, err := m.UserService.User(userID)
	if errors.Is(err, &ErrUserNotFound{}) {
		err = fmt.Errorf("add product: %w", &ErrPermission{Reason: err})
		return nil, err
	}
	if err != nil {
		err = fmt.Errorf("add product: %w", err)
		return nil, err
	}

	// After the user check, add the product.
	p, err = m.ProductService.AddProduct(p)
	if err != nil {
		err = fmt.Errorf("add product: %w", err)
		return nil, err
	}
	return p, nil
}

// ReplaceProduct updates information about the product with the new one by product ID.
func (m *Market) ReplaceProduct(p *Product, userID string) (*Product, error) {
	// Check the user for permission. Only the existence of the user counts yet.
	_, err := m.UserService.User(userID)
	if errors.Is(err, &ErrUserNotFound{}) {
		err = fmt.Errorf("edit product: %w", &ErrPermission{Reason: err})
		return nil, err
	}
	if err != nil {
		err = fmt.Errorf("edit product: %w", err)
		return nil, err
	}

	// After the user check, add the product.
	p, err = m.ProductService.ReplaceProduct(p)
	if err != nil {
		err = fmt.Errorf("edit product: %w", err)
		return nil, err
	}
	return p, nil
}

// DeleteProduct deletes the product from the market by its ID
// checking the permission to do it by user ID.
func (m *Market) DeleteProduct(id int, userID string) error {
	// Obtain the product.
	product, err := m.ProductService.Product(id)
	if err != nil {
		err = fmt.Errorf("delete product: %w", err)
		return err
	}

	// Check the user for permission. The user existence does not count here.
	if product.Seller != userID {
		err = fmt.Errorf("delete product: %w", &ErrPermission{Reason: err})
		return err
	}

	// Perform deletion.
	err = m.ProductService.DeleteProduct(id)
	if err != nil {
		err = fmt.Errorf("delete product: %w", err)
		return err
	}
	return nil
}
