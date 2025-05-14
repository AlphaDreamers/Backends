package provider

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GormPostgres(v *viper.Viper, log *logrus.Logger) *gorm.DB {
	var db *gorm.DB
	var err error

	// Get database configuration from viper
	host := v.GetString("database.host")
	port := v.GetString("database.port")
	user := v.GetString("database.user")
	password := v.GetString("database.password")
	dbname := v.GetString("database.name")
	sslmode := v.GetString("database.sslmode")

	// If any of these are empty, use defaults
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if password == "" {
		password = "postgres"
	}
	if dbname == "" {
		dbname = "appDatabase"
	}
	if sslmode == "" {
		sslmode = "disable"
	}

	// Build connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	log.Infof("Attempting to connect to database with DSN: host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		host, port, user, dbname, sslmode, password)

	maxAttempts := v.GetInt("database.max_attempts")
	if maxAttempts <= 0 {
		maxAttempts = 10
	}

	for i := 0; i < maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			log.Warnf("Failed to connect to DB (attempt %d): %v", i+1, err)
			time.Sleep(3 * time.Second)
			continue
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Warnf("Failed to get raw DB (attempt %d): %v", i+1, err)
			time.Sleep(3 * time.Second)
			continue
		}

		maxIdleConns := v.GetInt("database.max_idle_conns")
		if maxIdleConns <= 0 {
			maxIdleConns = 10
		}
		sqlDB.SetMaxIdleConns(maxIdleConns)

		maxOpenConns := v.GetInt("database.max_open_conns")
		if maxOpenConns <= 0 {
			maxOpenConns = 100
		}
		sqlDB.SetMaxOpenConns(maxOpenConns)

		connMaxLifetime := v.GetDuration("database.conn_max_lifetime")
		if connMaxLifetime <= 0 {
			connMaxLifetime = time.Hour
		}
		sqlDB.SetConnMaxLifetime(connMaxLifetime)

		log.Info("Successfully connected to the database")
		return db
	}

	log.Fatal("Could not connect to the database after multiple attempts. Check your database credentials and ensure PostgreSQL is running.")
	return nil
}
