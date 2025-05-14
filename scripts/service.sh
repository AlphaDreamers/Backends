#!/bin/zsh

module_names=("auth" "chat" "gig" "wallet" "order")

read -r -d '' COMMON_CODE <<EOF
package {{pkg}}

import (
  	"github.com/SwanHtetAungPhyo/backend/internal/repo/{{pkg}}"
 	"github.com/sirupsen/logrus"
)

type Service struct {
  repo {{pkg}}.Repository
  log  *logrus.Logger
}

func NewService(
  repo {{pkg}}.Repository,
  log *logrus.Logger,
) *Service {
  return &Service{
    repo: repo,
    log: log,
  }
}
EOF

# Create service.go files for each module
for mod in "${module_names[@]}"; do
	mkdir -p "internal/service/$mod"
	echo "${COMMON_CODE//\{\{pkg\}\}/$mod}" > "internal/service/$mod/service.go"
done
