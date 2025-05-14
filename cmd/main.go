package main

import (
	"github.com/SwanHtetAungPhyo/backend/cmd/auth"
	"github.com/SwanHtetAungPhyo/backend/cmd/chat"
	"github.com/SwanHtetAungPhyo/backend/cmd/gig"
	"github.com/SwanHtetAungPhyo/backend/cmd/order"
	"github.com/SwanHtetAungPhyo/backend/cmd/provider"
	"github.com/SwanHtetAungPhyo/backend/cmd/wallet"
	"go.uber.org/fx"
)

var ServerModule = fx.Module("server", fx.Provide(
	auth.NewServerState,
	chat.NewServerState,
	gig.NewServerState,
	order.NewServerState,
	wallet.NewServerState,
))

var ServerLifeCycleModule = fx.Module("lifecycle", fx.Provide(
	auth.RegisterLifeauthCycle,
	chat.RegisterLifechatCycle,
	gig.RegisterLifegigCycle,
	order.RegisterLifeorderCycle,
	wallet.RegisterLifewalletCycle,
))

func main() {
	app := fx.New(
		provider.ProviderModule,
		ServerModule,
		fx.Invoke(
			ServerLifeCycleModule,
		),
	)
	app.Run()
}
