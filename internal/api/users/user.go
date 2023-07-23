package users

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserServices struct {
	Conn *sqlx.DB
}

type User struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

type Users []User

func (user *User) ToJSON() []byte {
	userJSON, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	return userJSON
}

func (user *User) CheckPassword(password string) bool {
	return user.Password == password
}

func (users Users) ToJSON() []byte {
	usersJson, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}

	return usersJson
}
