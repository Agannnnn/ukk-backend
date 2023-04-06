package utils

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, username, password, dbname, port)
	var err error
	// db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	db.AutoMigrate(
		&Petugas{},
		&Masyarakat{},
		&Barang{},
		&Lelang{},
		&Penawaran{},
		&FotoBarang{},
		&Kategori{},
		&DetailKategori{},
		&Langganan{},
	)

	if err != nil {
		panic(err)
	}
}

func DB() *gorm.DB {
	return db
}
