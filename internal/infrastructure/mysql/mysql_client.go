package mysql

import (
	"fmt"
	"go.uber.org/dig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"twitter/internal/closer"
)

type mySqlClientParams struct {
	dig.In

	Environment string `name:"environment"`
}

var connectionString = map[string]string{
	"testing": "root:aspire@tcp(127.0.0.1:3306)",
}

const mainDB = "twitter"

// MySQLClientHandler is a struct that holds the database connection.
type MySQLClientHandler struct {
	DB *gorm.DB
}

// newMySQLClientHandler creates a new instance of the MySQLHandler.
func newMySQLClientHandler(params mySqlClientParams) (*MySQLClientHandler, closer.CloserResult, error) {
	log.Println("Creating new MySQL client handler")
	// Open up our database connection.
	mysqlDSN := fmt.Sprintf("%s/%s?parseTime=true", connectionString[params.Environment], mainDB)
	db, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Println("Error while connecting to the database", err)
		return nil, closer.CloserResult{}, err
	}
	//
	sqlDB, err := db.DB() // Get the underlying sql.DB instance to get the close method
	if err != nil {
		return nil, closer.CloserResult{}, err
	}

	// this will ensure that the mysql connection is closed when the application is closed
	cl := func() error {
		return sqlDB.Close()
	}

	return &MySQLClientHandler{
		DB: db,
	}, closer.CloserResult{Close: cl}, nil
}
