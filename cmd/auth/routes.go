package auth

import "github.com/gofiber/fiber/v2"

func (s *ServerState) setupRoutes() {
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
