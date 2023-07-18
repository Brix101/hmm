package services

import "github.com/jmoiron/sqlx"

type UserStorage struct {
	Conn *sqlx.DB
}

func NewUserStorage(conn *sqlx.DB) *UserStorage {
	return &UserStorage{Conn: conn}
}

type NewUser struct {
	Name     string
	Email    string
	Password string
}

func (s *UserStorage) CreateNewUser(data NewUser) (int, error) {
	res, err := s.Conn.Exec("insert into users (name, email, password) values (?, ?, ?)", data.Name, data.Email, data.Password)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
