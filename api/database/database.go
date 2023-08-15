package database

import (
	"log"
	"os"

	"github.com/fbcharles747/fiber-api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {
	dsn := "root:mysql@tcp(mysqldb:3306)/user1_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Fail to connect to the database!\n", err.Error())
		os.Exit(2)
	}

	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migration")
	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})

	Database = DbInstance{Db: db}

}
