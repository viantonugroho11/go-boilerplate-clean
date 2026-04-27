package model

import (
	userEntity "go-boilerplate-clean/internal/entity/users"
)

type User struct {
	ID    string
	Name  string
	Email string
}

func (u *User) TableName() string {
	return "users"
}

func ToEntity(u *User) userEntity.User {
	return userEntity.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func ToModel(u userEntity.User) User {
	return User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}