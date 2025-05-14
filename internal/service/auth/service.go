package auth

import (
  	"github.com/SwanHtetAungPhyo/backend/internal/repo/auth"
 	"github.com/sirupsen/logrus"
)

type Service struct {
  repo auth.Repository
  log  *logrus.Logger
}

func NewService(
  repo auth.Repository,
  log *logrus.Logger,
) *Service {
  return &Service{
    repo: repo,
    log: log,
  }
}
