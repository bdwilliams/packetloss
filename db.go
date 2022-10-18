package main

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() {
	log.Println("Connecting to the database...")
	dsn := os.Getenv("MYSQL_USERNAME") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ":3306)/" + os.Getenv("MYSQL_DATABASE") + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}

func DBMigrate() {
	log.Println("Running migrations...")
	DB.AutoMigrate(&PingResults{}, &Pings{})
}

func DBInsertPingResults(pingResults PingResults) error {
	err := DB.Create(&pingResults)
	return err.Error
}

func DBVerify() bool {
	MYSQL_ENABLED := os.Getenv("MYSQL_ENABLED")
	if MYSQL_ENABLED == "true" {
		log.Println("MySQL is enabled")
		return true
	} else {
		log.Println("MySQL is disabled")
	}

	return false
}
