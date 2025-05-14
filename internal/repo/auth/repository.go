package auth

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository struct {
	gorm *gorm.DB
	log  *logrus.Logger
}

func NewRepository(
	db *gorm.DB,
	log *logrus.Logger,
) *Repository {
	return &Repository{
		gorm: db,
		log:  log,
	}
}
