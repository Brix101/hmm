package services

import "github.com/jmoiron/sqlx"

type FeatureServices struct {
	Conn *sqlx.DB
}

func NewFeatureServices(conn *sqlx.DB) *FeatureServices {
	return &FeatureServices{Conn: conn}
}

type Table struct {
	Name string `db:"name"`
}

func (s *FeatureServices) GetAll() ([]string, error) {
	var tables []Table
	err := s.Conn.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		return nil, err
	}

	tableNames := make([]string, len(tables))
	for i, table := range tables {
		tableNames[i] = table.Name
	}

	return tableNames, nil
}
