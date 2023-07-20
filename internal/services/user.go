package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserServices struct {
	Conn *sqlx.DB
}

func NewUserServices(conn *sqlx.DB) *UserServices {
	return &UserServices{Conn: conn}
}

type UserEntity struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

type NewUser struct {
	Name     string
	Email    string
	Password string
}

func (s *UserServices) CreateUser(data NewUser) (*UserEntity, error) {
	// Generate a new UUID
	newID := uuid.New().String()

	res, err := s.Conn.Exec("insert into users (id, name, email, password) values (?, ?, ?, ?)", newID, data.Name, data.Email, data.Password)
	fmt.Println(res)
	if err != nil {
		return nil, err
	}
	// Fetch the created user from the database using the UUID
	user, err := s.GetUserByID(uuid.MustParse(newID))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServices) ListUser() ([]UserEntity, error) {
	users := []UserEntity{}
	err := s.Conn.Select(&users, "SELECT id, name, email FROM users")
	if err != nil {
		return users, err
	}
	return users, nil
}

func (s *UserServices) GetUserByID(userId uuid.UUID) (*UserEntity, error) {
	user := UserEntity{}
	err := s.Conn.Get(&user, "SELECT * FROM users WHERE id=?", userId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserServices) GetUserByEmail(email string) (*UserEntity, error) {
	user := UserEntity{}
	err := s.Conn.Get(&user, "SELECT * FROM users WHERE email=?", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
