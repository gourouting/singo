package model

import (
	"log"
	"os"
	"singo/util"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the singleton database connection.
var DB *gorm.DB

// Database initializes the MySQL connection.
func Database(connString string) {
	// Initialize the GORM logger configuration.
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level. Adjust as needed.
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(connString), &gorm.Config{
		Logger: newLogger,
	})
	// Error
	if connString == "" || err != nil {
		util.Log().Error("mysql lost: %v", err)
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		util.Log().Error("mysql lost: %v", err)
		panic(err)
	}

	// Set the connection pool.
	// Idle connections.
	sqlDB.SetMaxIdleConns(10)
	// Open connections.
	sqlDB.SetMaxOpenConns(20)
	DB = db

	migration()
}
