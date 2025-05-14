package gig

import (
    "github.com/SwanHtetAungPhyo/backend/internal/service/gig"
 	"github.com/sirupsen/logrus"
)

type Handler struct {
  srv gig.Service
  log  *logrus.Logger
}

func NewHandler(
  srv gig.Service,
  log *logrus.Logger,
) *Handler {
  return &Handler{
    srv: srv,
    log: log,
  }
}
