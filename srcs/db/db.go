package db

import (
	"database/sql"
	"fmt"

	"github.com/Palm78070/basic_web_app/settings"
	_ "github.com/lib/pq" //driver for PostgreSQL database
)

//Function to connect to db
//*sql.DB object that have important information to connect to db
func Connect(dbConfig *settings.DBSettings) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
	)
	// connStr := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

	//db driver, interface to interact with db system
	db, err := sql.Open("postgres", connStr)
	return db, err
}
