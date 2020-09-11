package market_test

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ortymid/t2-http/market"
	"github.com/ortymid/t2-http/market/mock"
)

type MockFuncUser struct {
	expect     bool
	argID      string
	returnUser *market.User
	returnErr  error
}
type MockUserService struct {
	User MockFuncUser
}

func (opt *MockUserService) Setup(m *mock.MockUserService) {
	if opt.User.expect {
		m.EXPECT().User(gomock.Eq(opt.User.argID)).Return(opt.User.returnUser, opt.User.returnErr)
	} else {
		m.EXPECT().User(nil).MaxTimes(0)
	}
}

type MockFuncProducts struct {
	expect         bool
	returnProducts []*market.Product
	returnErr      error
}
type MockFuncProduct struct {
	expect        bool
	argID         int
	returnProduct *market.Product
	returnErr     error
}
type MockFuncAddProduct struct {
	expect        bool
	argProduct    *market.Product
	returnProduct *market.Product
	returnErr     error
}
type MockFuncReplaceProduct struct {
	expect        bool
	argProduct    *market.Product
	returnProduct *market.Product
	returnErr     error
}
type MockFuncDeleteProduct struct {
	expect    bool
	argID     int
	returnErr error
}
type MockProductService struct {
	Products       MockFuncProducts
	Product        MockFuncProduct
	AddProduct     MockFuncAddProduct
	ReplaceProduct MockFuncReplaceProduct
	DeleteProduct  MockFuncDeleteProduct
}

func (opt *MockProductService) Setup(m *mock.MockProductService) {
	if opt.Products.expect {
		m.EXPECT().Products().Return(opt.Products.returnProducts, opt.Products.returnErr)
	} else {
		m.EXPECT().Products().MaxTimes(0)
	}
	if opt.Product.expect {
		m.EXPECT().Product(opt.Product.argID).Return(opt.Product.returnProduct, opt.Product.returnErr)
	} else {
		m.EXPECT().Product(nil).MaxTimes(0)
	}
	if opt.AddProduct.expect {
		m.EXPECT().AddProduct(opt.AddProduct.argProduct).Return(opt.AddProduct.returnProduct, opt.AddProduct.returnErr)
	} else {
		m.EXPECT().AddProduct(nil).MaxTimes(0)
	}
	if opt.ReplaceProduct.expect {
		m.EXPECT().ReplaceProduct(opt.ReplaceProduct.argProduct).Return(opt.ReplaceProduct.returnProduct, opt.ReplaceProduct.returnErr)
	} else {
		m.EXPECT().ReplaceProduct(nil).MaxTimes(0)
	}
	if opt.DeleteProduct.expect {
		m.EXPECT().DeleteProduct(opt.DeleteProduct.argID).Return(opt.DeleteProduct.returnErr)
	} else {
		m.EXPECT().DeleteProduct(nil).MaxTimes(0)
	}
}

func TestMarket_Products(t *testing.T) {
	type mocks struct {
		UserService    MockUserService
		ProductService MockProductService
	}
	tests := []struct {
		name    string
		mocks   mocks
		want    []*market.Product
		wantErr bool
	}{
		{
			name: "Returns products",
			mocks: mocks{
				ProductService: MockProductService{
					Products: MockFuncProducts{
						expect: true,
						returnProducts: []*market.Product{
							{ID: 1, Name: "p1", Price: 100, Seller: "1"},
							{ID: 2, Name: "p2", Price: 200, Seller: "2"},
						},
					},
				},
			},
			want: []*market.Product{
				{ID: 1, Name: "p1", Price: 100, Seller: "1"},
				{ID: 2, Name: "p2", Price: 200, Seller: "2"},
			},
		},
		{
			name: "Returns empty products",
			mocks: mocks{
				ProductService: MockProductService{
					Products: MockFuncProducts{
						expect:         true,
						returnProducts: []*market.Product{},
					},
				},
			},
			want: []*market.Product{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			us := mock.NewMockUserService(ctrl)
			tt.mocks.UserService.Setup(us)

			ps := mock.NewMockProductService(ctrl)
			tt.mocks.ProductService.Setup(ps)

			m := &market.Market{
				UserService:    us,
				ProductService: ps,
			}
			got, err := m.Products()
			if (err != nil) != tt.wantErr {
				t.Errorf("Market.Products() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Market.Products() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarket_Product(t *testing.T) {
	type mocks struct {
		UserService    MockUserService
		ProductService MockProductService
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    *market.Product
		wantErr bool
	}{
		{
			name: "Returns a product",
			mocks: mocks{
				ProductService: MockProductService{
					Product: MockFuncProduct{
						expect:        true,
						argID:         1,
						returnProduct: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
					},
				},
			},
			args: args{
				id: 1,
			},
			want: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
		},
		{
			name: "Returns an error for not existing product",
			mocks: mocks{
				ProductService: MockProductService{
					Product: MockFuncProduct{
						expect:    true,
						argID:     1,
						returnErr: market.ErrProductNotFound,
					},
				},
			},
			args: args{
				id: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			us := mock.NewMockUserService(ctrl)
			tt.mocks.UserService.Setup(us)

			ps := mock.NewMockProductService(ctrl)
			tt.mocks.ProductService.Setup(ps)

			m := &market.Market{
				UserService:    us,
				ProductService: ps,
			}
			got, err := m.Product(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Market.Products() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Market.Products() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarket_AddProduct(t *testing.T) {
	type mocks struct {
		UserService    MockUserService
		ProductService MockProductService
	}
	type args struct {
		p      *market.Product
		userID string
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    *market.Product
		wantErr bool
	}{
		{
			name: "Adds product",
			mocks: mocks{
				UserService: MockUserService{
					User: MockFuncUser{
						expect:     true,
						argID:      "1",
						returnUser: &market.User{ID: "1", Name: "u1"},
					},
				},
				ProductService: MockProductService{
					AddProduct: MockFuncAddProduct{
						expect:        true,
						argProduct:    &market.Product{Name: "p1", Price: 100, Seller: "1"},
						returnProduct: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
					},
				},
			},
			args: args{
				p:      &market.Product{Name: "p1", Price: 100, Seller: "1"},
				userID: "1",
			},
			want: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
		},
		{
			name: "Returns an error for not existing user",
			mocks: mocks{
				UserService: MockUserService{
					User: MockFuncUser{
						expect:    true,
						argID:     "1",
						returnErr: &market.ErrUserNotFound{},
					},
				},
			},
			args: args{
				p:      &market.Product{Name: "p1", Price: 100, Seller: "1"},
				userID: "1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			us := mock.NewMockUserService(ctrl)
			tt.mocks.UserService.Setup(us)

			ps := mock.NewMockProductService(ctrl)
			tt.mocks.ProductService.Setup(ps)

			m := &market.Market{
				UserService:    us,
				ProductService: ps,
			}
			got, err := m.AddProduct(tt.args.p, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Market.Products() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Market.Products() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarket_ReplaceProduct(t *testing.T) {
	type mocks struct {
		UserService    MockUserService
		ProductService MockProductService
	}
	type args struct {
		p      *market.Product
		userID string
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    *market.Product
		wantErr bool
	}{
		{
			name: "Should replace a product",
			mocks: mocks{
				UserService: MockUserService{
					User: MockFuncUser{
						expect:     true,
						argID:      "1",
						returnUser: &market.User{ID: "1", Name: "u1"},
					},
				},
				ProductService: MockProductService{
					ReplaceProduct: MockFuncReplaceProduct{
						expect:        true,
						argProduct:    &market.Product{ID: 1, Name: "p2", Price: 200, Seller: "1"},
						returnProduct: &market.Product{ID: 1, Name: "p2", Price: 200, Seller: "1"},
					},
				},
			},
			args: args{
				p:      &market.Product{ID: 1, Name: "p2", Price: 200, Seller: "1"},
				userID: "1",
			},
			want: &market.Product{ID: 1, Name: "p2", Price: 200, Seller: "1"},
		},
		{
			name: "Returns an error for not existing user",
			mocks: mocks{
				UserService: MockUserService{
					User: MockFuncUser{
						expect:    true,
						argID:     "1",
						returnErr: &market.ErrUserNotFound{},
					},
				},
			},
			args: args{
				p:      &market.Product{Name: "p1", Price: 100, Seller: "1"},
				userID: "1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			us := mock.NewMockUserService(ctrl)
			tt.mocks.UserService.Setup(us)

			ps := mock.NewMockProductService(ctrl)
			tt.mocks.ProductService.Setup(ps)

			m := &market.Market{
				UserService:    us,
				ProductService: ps,
			}
			got, err := m.ReplaceProduct(tt.args.p, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Market.Products() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Market.Products() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarket_DeleteProduct(t *testing.T) {
	type mocks struct {
		UserService    MockUserService
		ProductService MockProductService
	}
	type args struct {
		id     int
		userID string
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    *market.Product
		wantErr bool
	}{
		{
			name: "Deletes product",
			mocks: mocks{
				ProductService: MockProductService{
					Product: MockFuncProduct{
						expect:        true,
						argID:         1,
						returnProduct: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
					},
					DeleteProduct: MockFuncDeleteProduct{
						expect: true,
						argID:  1,
					},
				},
			},
			args: args{
				id:     1,
				userID: "1",
			},
			want: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
		},
		{
			name: "Returns an error for not existing product",
			mocks: mocks{
				ProductService: MockProductService{
					Product: MockFuncProduct{
						expect:    true,
						argID:     1,
						returnErr: market.ErrProductNotFound,
					},
				},
			},
			args: args{
				id:     1,
				userID: "1",
			},
			wantErr: true,
		},
		{
			name: "Returns an error for user mismatch",
			mocks: mocks{
				ProductService: MockProductService{
					Product: MockFuncProduct{
						expect:        true,
						argID:         1,
						returnProduct: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
					},
				},
			},
			args: args{
				id:     1,
				userID: "2",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			us := mock.NewMockUserService(ctrl)
			tt.mocks.UserService.Setup(us)

			ps := mock.NewMockProductService(ctrl)
			tt.mocks.ProductService.Setup(ps)

			m := &market.Market{
				UserService:    us,
				ProductService: ps,
			}
			err := m.DeleteProduct(tt.args.id, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Market.Products() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
