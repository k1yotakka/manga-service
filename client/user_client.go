package client

import (
	"encoding/json"
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
		Get(fmt.Sprintf("http://user-service:8001/api/users/%d", id))

	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(resp.Body(), &user); err != nil {
		return nil, err
	}

	return &user, nil
}
