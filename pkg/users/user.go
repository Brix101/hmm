package users

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserServices struct {
	Conn *sqlx.DB
}

type UserEntity struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

type UserEntities []UserEntity

type NewUser struct {
	Name     string
	Email    string
	Password string
}

type UserRequestBody struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (user *UserEntity) ToJSON() []byte {
	userJSON, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	return userJSON
}

func (users UserEntities) ToJSON() []byte {
	usersJson, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}

	return usersJson
}
