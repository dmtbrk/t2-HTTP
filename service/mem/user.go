package mem

import (
	"sync"

	"github.com/ortymid/t2-http/market"
)

type UserService struct {
	mu     sync.RWMutex
	lastID int
	users  []*market.User
}

func NewUserService() *UserService {
	users := []*market.User{
		{ID: "1", Name: "admin"},
		{ID: "2", Name: "Dmytro"},
	}
	return &UserService{users: users, lastID: 2}
}

func (srv *UserService) User(id string) (*market.User, error) {
	srv.mu.RLock()
	defer srv.mu.RUnlock()

	for _, u := range srv.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, &market.ErrUserNotFound{UserID: id}
}

// func (srv *UserService) AddUser(u *market.User) (*market.User, error) {
// 	srv.mu.Lock()
// 	defer srv.mu.Unlock()

// 	srv.lastID++
// 	u.ID = market.UUID(strconv.Itoa(srv.lastID))

// 	srv.users = append(srv.users, u)

// 	return u, nil
// }
