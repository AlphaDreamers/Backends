package provider

import (
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

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(v.GetString("database.url")), &gorm.Config{
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

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		return db
	}

	log.Fatal("Could not connect to the database after 10 attempts.")
	return nil
}
