package config

import "database/sql"

func DBConnection() (*sql.DB, error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "crud_go"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"/@"+dbName)
	return db, err
}