package wallet

import (
	"github.com/SwanHtetAungPhyo/backend/internal/service/wallet"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	srv wallet.Service
	log *logrus.Logger
}

func NewHandler(
	srv wallet.Service,
	log *logrus.Logger,
) *Handler {
	return &Handler{
		srv: srv,
		log: log,
	}
}
