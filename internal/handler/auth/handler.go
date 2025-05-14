package auth

import (
	"github.com/SwanHtetAungPhyo/backend/internal/model"
	"github.com/SwanHtetAungPhyo/backend/internal/service/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
)

func NewHandler(
	srv auth.Service,
	log *logrus.Logger,
) *Handler {
	return &Handler{
		srv: srv,
		log: log,
	}
}

type Handler struct {
	srv auth.Service
	log *logrus.Logger
}

func (h Handler) SignIn(ctx *fiber.Ctx) error {
	var req *model.UserSignInReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResp{
			Message: err.Error(),
		})
	}
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   ctx.Hostname(),
		Secure:   true,
		HTTPOnly: true,
		SameSite: string(rune(http.SameSiteLaxMode)),
	})

	return ctx.Status(fiber.StatusOK).JSON(&model.UserSignInResp{})
}

func (h Handler) SignUp(ctx *fiber.Ctx) error {
	var req *model.UserSignUpRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResp{
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(&model.UserSignUpResp{})
}

func (h Handler) Logout(ctx *fiber.Ctx) error {
	var req *model.UserLogoutRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResp{
			Message: err.Error(),
		})
	}

	ctx.ClearCookie("refresh_token")
	return ctx.Status(fiber.StatusOK).JSON(model.Response{})
}

func (h Handler) VerifyEmail(ctx *fiber.Ctx) error {
	var req *model.EmailVerificationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResp{
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(model.Response{})
}

func (h Handler) ForgotPassword(ctx *fiber.Ctx) error {
	var req *model.ForgotPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResp{
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(model.Response{})
}

func (h Handler) ForgotPasswordConfirm(ctx *fiber.Ctx) error {
	var req *model.ForgotPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResp{
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(model.Response{
		Message: "you will get the forgot password code to the email that you registered",
	})
}

func (h Handler) Me(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")
	if userId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResp{
			Message: "userId is required",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(model.Response{
		Message: "you are logged in",
	})
}

func (h Handler) KycVerify(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")
	if userId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResp{
			Message: "userId is required in param",
		})
	}
	files, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResp{
			Message: err.Error(),
		})
	}
	idPhoto := files.File["id_photo"]
	if idPhoto == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResp{
			Message: "id_photo is required",
		})
	}
	selfie := files.File["selfie"]
	if selfie == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResp{
			Message: "selfie is required",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(model.Response{
		Message: userId + "!! you kyc is verified",
	})

}
