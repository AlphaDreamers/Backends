package wallet

import (
  	"github.com/SwanHtetAungPhyo/backend/internal/repo/wallet"
 	"github.com/sirupsen/logrus"
)

type Service struct {
  repo wallet.Repository
  log  *logrus.Logger
}

func NewService(
  repo wallet.Repository,
  log *logrus.Logger,
) *Service {
  return &Service{
    repo: repo,
    log: log,
  }
}
