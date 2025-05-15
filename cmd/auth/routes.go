package auth

import "github.com/gofiber/fiber/v2"

func (s *ServerState) setupAuthRoutes() {
	authGroup := s.fiberApp.Group("/auth")
	authGroup.Post("/sign-in", s.handler.SignIn)
	authGroup.Post("/sign-up", s.handler.SignUp)
	authGroup.Post("/log-out", s.handler.Logout)
	authGroup.Post("/verify-email", s.handler.Confirm)
	authGroup.Post("/forgot-password", s.handler.ForgotPassword)
	authGroup.Post("/forgot-password/confirm", s.handler.ForgotPasswordConfirm)
	authGroup.Post("/reset-password", s.handler.ForgotPassword)
	authGroup.Post("/reset-password/confirm", s.handler.ForgotPasswordConfirm)
	authGroup.Post("/me", s.handler.Me)
	authGroup.Post("/kyc-verify", s.handler.KYCVerify)
	authGroup.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "OK from auth service",
			"data":    "I am on live ",
		})
	})
}

func (s *ServerState) setUpChatRoutes() {
	chatGroup := s.fiberApp.Group("/chat")
	chatGroup.Get("/sse", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "OK from chat service",
		})
	})
}

func (s *ServerState) setUpGigRoutes() {
	gigGroup := s.fiberApp.Group("/gig")
	gigGroup.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "OK from gig service",
		})
	})
}

func (s *ServerState) setUpOrderRoutes() {
	orderGroup := s.fiberApp.Group("/order")
	orderGroup.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "OK from order service",
		})
	})
}

func (s *ServerState) setUpPaymentRoutes() {
	paymentGroup := s.fiberApp.Group("/payment")
	paymentGroup.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "OK from payment service",
		})
	})
}

func (s *ServerState) setUpWalletRoutes() {
	walletGroup := s.fiberApp.Group("/wallet")
	walletGroup.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "OK from wallet service",
		})
	})
}
