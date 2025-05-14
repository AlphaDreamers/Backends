#!/bin/zsh

mkdir -p cmd/{auth,chat,gig,order,wallet}
mkdir -p cmd/middleware
mkdir -p cmd/provider
mkdir -p migration
mkdir -p internal/hanlder/{auth,chat,gig,order,wallet}
mkdir -p internal/service/{auth,chat,gig,order,wallet}
mkdir -p internal/repo/{auth,chat,gig,order,wallet}
touch  docker-compose.yaml
touch  cmd/auth/server.go
touch  cmd/auth/routes.go
touch  cmd/auth/middleware.go
touch cmd/auth/lifeCycle.go

touch  cmd/chat/server.go
touch  cmd/chat/routes.go
touch  cmd/chat/middleware.go
touch cmd/chat/lifeCycle.go


touch  cmd/gig/server.go
touch  cmd/gig/routes.go
touch  cmd/gig/middleware.go
touch cmd/gig/lifeCycle.go

touch  cmd/order/server.go
touch  cmd/order/routes.go
touch  cmd/order/middleware.go
touch cmd/order/lifeCycle.go

touch  cmd/wallet/server.go
touch  cmd/wallet/routes.go
touch  cmd/wallet/middleware.go
touch cmd/wallet/lifeCycle.go


touch cmd/middleware/jwt_middleware.go
touch cmd/middleware/auth_middleware.go
