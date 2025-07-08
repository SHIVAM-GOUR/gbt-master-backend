package config

import (
	"log"

	"github.com/SHIVAM-GOUR/gbt-master-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// dsn := os.Getenv("RAILWAY_DATABASE_URL")
	dsn := "postgresql://postgres:fTCNqYhMCcceJvrZJknDNNvxmgRzUXLc@switchyard.proxy.rlwy.net:30216/railway"

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(&models.Class{}, &models.Inquiry{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully")
}
