package chat

import (
	"github.com/SwanHtetAungPhyo/backend/internal/service/chat"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	srv chat.Service
	log *logrus.Logger
}

func NewHandler(
	srv chat.Service,
	log *logrus.Logger,
) *Handler {
	return &Handler{
		srv: srv,
		log: log,
	}
}
