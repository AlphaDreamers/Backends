package gig

import (
	"github.com/SwanHtetAungPhyo/backend/internal/repo/gig"
	"github.com/sirupsen/logrus"
)

type Service struct {
	repo gig.Repository
	log  *logrus.Logger
}

func NewService(
	repo gig.Repository,
	log *logrus.Logger,
) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}
