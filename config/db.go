package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const SearchLimit = 10

func ConnectMySQL() *gorm.DB {
	databaseConfig := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err := gorm.Open(mysql.Open(databaseConfig), initConfig())

	if err != nil {
		panic("Fail To Connect MySQL Database")
	}
	migrate(db)
	return db
}

func ConnectSQLite() *gorm.DB {
	// Use app.db as the SQLite database file
	dbPath := "app.db"
	var err error
	db, err := gorm.Open(sqlite.Open(dbPath), initConfig())
	if err != nil {
		panic("Fail To Connect SQLite Database")
	}
	migrate(db)
	return db
}

// InitConfig Initialize Config
func initConfig() *gorm.Config {
	return &gorm.Config{
		Logger:         initLog(),
		NamingStrategy: initNamingStrategy(),
	}
}

// InitLog Connection Log Configuration
func initLog() logger.Interface {
	f, _ := os.Create("gorm.log")
	// Log to both file and standard output
	multiWriter := io.MultiWriter(f, os.Stdout)
	newLogger := logger.New(log.New(multiWriter, "\r\n", log.LstdFlags), logger.Config{
		Colorful:      true,
		LogLevel:      logger.Info, // This will show SQL queries
		SlowThreshold: time.Second,
	})
	return newLogger
}

// InitNamingStrategy Init NamingStrategy
func initNamingStrategy() *schema.NamingStrategy {
	return &schema.NamingStrategy{
		SingularTable: false,
		TablePrefix:   "",
	}
}
