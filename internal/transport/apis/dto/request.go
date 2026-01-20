package dto

import "go-boilerplate-clean/internal/entity/users"

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// to entity
func (r *CreateUserRequest) ToEntity() users.User {
	return users.User{
		Name:  r.Name,
		Email: r.Email,
	}
}			

