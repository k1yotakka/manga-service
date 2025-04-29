package client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

var Client = resty.New()

type User struct {
	ID       uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func GetUserByID(id uint) (*User, error) {
	resp, err := Client.R().
		SetResult(&User{}).
		Get(fmt.Sprintf("http://localhost:8001/api/users/%d", id))

	if err != nil {
		return nil, err
	}

	return resp.Result().(*User), nil
}
