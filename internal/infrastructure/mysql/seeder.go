// Package mysql-seeder.go is used to seed the database with some data only for testing purposes/dev env
package mysql

import (
	"log"

	"go.uber.org/dig"
	"gorm.io/gorm"
	"twitter/internal/models"
)

type MysqlSeederParams struct {
	dig.In

	MySQLClientHandler *MySQLClientHandler
}

type MySQLSeederHandler struct {
	db *gorm.DB
}

func NewMySQLSeederHandler(params MysqlSeederParams) *MySQLSeederHandler {
	log.Println("Creating new MySQL seeder handler")
	return &MySQLSeederHandler{
		db: params.MySQLClientHandler.DB,
	}
}

func (m *MySQLSeederHandler) MigrateSchema() {
	// Migrate the schema
	if err := m.db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate User model: %v", err)
	}
	if err := m.db.AutoMigrate(&models.Tweet{}); err != nil {
		log.Fatalf("Failed to migrate Tweet model: %v", err)
	}
	if err := m.db.AutoMigrate(&models.Followers{}); err != nil {
		log.Fatalf("Failed to migrate Followers model: %v", err)
	}
}

func (m *MySQLSeederHandler) SeedDatabase() {
	log.Println("Seeding the database with some data only for dev environment")
	//// Seed the database with some data only for dev environment
	//users := []models.User{
	//	{Name: "Piyush Singh", Email: "piyush@example.com", Password: "aspire"},
	//}
	//
	//for _, user := range users {
	//	if err := m.db.Where(models.User{Email: user.Email}).FirstOrCreate(&user).Error; err != nil {
	//		log.Printf("Error while seeding the database: %e", err)
	//	}
	//}
}
