#!/bin/zsh

module_names=("auth" "chat" "gig" "wallet" "order")

read -r -d '' COMMON_CODE <<EOF
package {{pkg}}

import (
  	"github.com/sirupsen/logrus"
  	"gorm.io/gorm"
)
type Repository struct {
  db *gorm.DB
  log *logrus.Logger
}

func NewRepository(
  db *gorm.DB,
  log *logrus.Logger,
) *Repository {
  return  &Repository{
    db: db,
    log: log,
  }
}
EOF

for mod in "${module_names[@]}"; do
	mkdir -p "cmd/$mod"
	echo "${COMMON_CODE//\{\{pkg\}\}/$mod}" > "internal/repo/$mod/repository.go"
done
