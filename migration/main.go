package main

import (
	"github.com/SwanHtetAungPhyo/backend/migration/model"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(v *viper.Viper) (*gorm.DB, error) {
	dsn := v.GetString("aws.rds.local")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&model.Badge{},
		&model.UserBadge{},
		&model.User{},
		&model.Skill{},
		&model.UserSkill{},
		&model.Biometrics{},
		&model.GigTag{},
		&model.Gig{},
		&model.GigImage{},
		&model.GigPackage{},
		&model.Category{},
		&model.Order{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
 