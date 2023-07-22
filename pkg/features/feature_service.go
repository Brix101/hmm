package features

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type FeatureServices struct {
	Conn *sqlx.DB
}

func NewFeatureServices(conn *sqlx.DB) *FeatureServices {
	return &FeatureServices{Conn: conn}
}

type Table struct {
	Name string `db:"name"`
}

type Features []string

func (features Features) ToJSON() []byte {
	jsonData, err := json.Marshal(features)
	if err != nil {
		panic(err)
	}

	return jsonData
}

func (s *FeatureServices) GetAll() (Features, error) {
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
