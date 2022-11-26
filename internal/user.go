package internal

import (
	"net/http"
	"time"
)

type User struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserList map[string]User

type UserStore struct {
	Increment int      `json:"increment"`
	List      UserList `json:"list"`
}

type UpdateUserRequest struct {
	Name string `json:"name"`
}

func (c *UpdateUserRequest) Bind(r *http.Request) error { return nil }

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c *CreateUserRequest) Bind(r *http.Request) error { return nil }
