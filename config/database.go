package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// DB is a exported connection
	DB *gorm.DB
)

// InitDatabase is initial Setup for DB Connection
func InitDatabase() *gorm.DB {
	var err error
	appEnv := GoDotEnvVariable("APPLICATION_ENV")
	sslMode := "sslmode=require"
	if appEnv == "development" {
		sslMode = "sslmode=disable"
	}
	dsn := "host=" + GoDotEnvVariable("DB_HOST") + " user=" + GoDotEnvVariable("DB_USERNAME") + " password=" + GoDotEnvVariable("DB_PASSWORD") + " dbname=" + GoDotEnvVariable("DB_DATABASE") + " port=" + GoDotEnvVariable("DB_PORT") + " " + sslMode + " TimeZone=" + GoDotEnvVariable("TZ")

	var logLvl logger.LogLevel
	if appEnv == "production" {
		logLvl = logger.Silent
	} else if appEnv == "staging" {
		logLvl = logger.Warn
	} else {
		logLvl = logger.Info
		// logLvl = logger.Warn
	}
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logLvl,      // Log level
			Colorful:      true,        // Enable color
		},
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      dbLogger,
		PrepareStmt: true,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		SkipDefaultTransaction: true,
	})

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	sqlDB, errConnPool := DB.DB()
	if errConnPool != nil {
		fmt.Println(errConnPool.Error())
		panic(errConnPool.Error())
	}

	// Ping
	errSqlPing := sqlDB.Ping()
	if errSqlPing != nil {
		fmt.Println(errSqlPing)
		panic("failed to connect database")
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("⚡️DB Connection opened!")

	return DB
}
