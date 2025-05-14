#!/bin/zsh

module_names=("auth" "chat" "gig" "wallet" "order")

read -r -d '' COMMON_CODE <<EOF
package {{pkg}}

import (
    "github.com/SwanHtetAungPhyo/backend/internal/service/{{pkg}}"
 	"github.com/sirupsen/logrus"
)

type Handler struct {
  srv {{pkg}}.Service
  log  *logrus.Logger
}

func NewHandler(
  srv {{pkg}}.Service,
  log *logrus.Logger,
) *Handler {
  return &Handler{
    srv: srv,
    log: log,
  }
}
EOF

for mod in "${module_names[@]}"; do
	echo "${COMMON_CODE//\{\{pkg\}\}/$mod}" > "internal/handler/$mod/handler.go"
done
