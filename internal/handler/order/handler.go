package order

import (
	"github.com/SwanHtetAungPhyo/backend/internal/service/order"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	srv order.Service
	log *logrus.Logger
}

func NewHandler(
	srv order.Service,
	log *logrus.Logger,
) *Handler {
	return &Handler{
		srv: srv,
		log: log,
	}
}
