package main

import (
	"database/sql"
	"fmt"
	q "magic/db_ops"
)

func connectDB(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbName) // if path doesn't exist, it will fail
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setupDB(dbName string) (*sql.DB, error) {

	db, err := connectDB(dbName)
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	schemas := []string{q.UserSchema, q.CardSchema, q.ImageUrisSchema, q.LegalitiesSchema, q.PricesSchema, q.RelatedUrisSchema, q.PurchaseUrisSchema}
	for _, schema := range schemas {
		_, err := db.Exec(schema)
		if err != nil {
			fmt.Println("Error while building schemas in DB: ", err)
			return nil, err
		}
	}

	fmt.Println("DB setup completed")
	return db, nil
}
