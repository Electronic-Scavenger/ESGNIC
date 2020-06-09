package orm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
)

var DB *gorm.DB

func Init() {
	db, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		logrus.Fatal("failed to connect database")
	}
	DB = db
	initTables()
	return
}

func initTables() {
	DB.AutoMigrate(User{})
	DB.AutoMigrate(Network{})
}
