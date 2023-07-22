package users

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var schema = `
	CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
	);
`

func NewUserServices(conn *sqlx.DB) *UserServices {
	// create the table users if not exists
	conn.MustExec(schema)

	return &UserServices{Conn: conn}
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
	user, err := s.GetByID(uuid.MustParse(newID))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServices) GetAll() (UserEntities, error) {
	users := UserEntities{}
	err := s.Conn.Select(&users, "SELECT id, name, email FROM users")
	if err != nil {
		return users, err
	}
	return users, nil
}

func (s *UserServices) GetByID(userId uuid.UUID) (*UserEntity, error) {
	user := UserEntity{}
	err := s.Conn.Get(&user, "SELECT * FROM users WHERE id=?", userId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserServices) GetByEmail(email string) (*UserEntity, error) {
	user := UserEntity{}
	err := s.Conn.Get(&user, "SELECT * FROM users WHERE email=?", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
