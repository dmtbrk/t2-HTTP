package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ortymid/t2-http/market"
)

type UserService struct {
	URL string
}

func NewUserService(url string) *UserService {
	return &UserService{URL: url}
}

func (srv *UserService) User(id string) (*market.User, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", srv.URL, id))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, &market.ErrUserNotFound{UserID: id}
		}
		return nil, errors.New("something went wrong")
	}

	data := struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Balance  int    `json:"balance"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &market.User{
		ID:   strconv.Itoa(data.ID),
		Name: data.Username,
	}, nil
}
