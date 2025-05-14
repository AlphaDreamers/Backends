package chat

import (
  	"github.com/SwanHtetAungPhyo/backend/internal/repo/chat"
 	"github.com/sirupsen/logrus"
)

type Service struct {
  repo chat.Repository
  log  *logrus.Logger
}

func NewService(
  repo chat.Repository,
  log *logrus.Logger,
) *Service {
  return &Service{
    repo: repo,
    log: log,
  }
}
