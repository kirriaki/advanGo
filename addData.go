package main

import (
	"encoding/csv"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"project/configs"
	"project/models"
	"strconv"
)

func importFromCSV(db *gorm.DB, csvPath string) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("could not open CSV file: %w", err)
	}
	defer file.Close()

	r := csv.NewReader(file)

	header, err := r.Read()
	if err != nil {
		return fmt.Errorf("could not read header from CSV: %w", err)
	}
	fmt.Println("Header columns:", header) // для отладки

	for {
		record, err := r.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			print(err)
		}

		name := record[1]
		city := record[3]
		store := record[4]

		quantity, err := strconv.Atoi(record[5])
		if err != nil {
			return fmt.Errorf("failed to parse quantity: %w", err)
		}

		price, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			return fmt.Errorf("failed to parse price: %w", err)
		}

		product := models.Product{
			Name:  name,
			City:  city,
			Shop:  store,
			Stock: quantity,
			Price: price,
		}

		if err := db.Create(&product).Error; err != nil {
			return fmt.Errorf("failed to insert product: %w", err)
		}
	}

	return nil
}

func main() {
	config := configs.LoadConfig()
	dsn := config.DatabaseURL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	if err := db.AutoMigrate(&models.Product{}); err != nil {
		log.Fatal("failed to migrate: ", err)
	}

	if err := importFromCSV(db, "MECHTA1.csv"); err != nil {
		log.Fatal("failed to import csv: ", err)
	}
}
