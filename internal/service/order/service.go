package order

import (
  	"github.com/SwanHtetAungPhyo/backend/internal/repo/order"
 	"github.com/sirupsen/logrus"
)

type Service struct {
  repo order.Repository
  log  *logrus.Logger
}

func NewService(
  repo order.Repository,
  log *logrus.Logger,
) *Service {
  return &Service{
    repo: repo,
    log: log,
  }
}
