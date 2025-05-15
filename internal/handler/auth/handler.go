package auth

import (
	"context"
	"github.com/SwanHtetAungPhyo/backend/internal/model"
	ar "github.com/SwanHtetAungPhyo/backend/internal/repo/auth"
	sa "github.com/SwanHtetAungPhyo/backend/internal/service/auth"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type Handler struct {
	log                *logrus.Logger
	cognitoClient      *cognitoidentityprovider.Client
	rekognitiionClient *rekognition.Client
	textractClient     *textract.Client
	repo               *ar.Repository
	srv                *sa.Service
	ctx                context.Context
	clientId           string
	v                  *viper.Viper
}

func NewHandler(
	cognitoClient *cognitoidentityprovider.Client,
	rekognitiionClient *rekognition.Client,
	textractClient *textract.Client,
	repo *ar.Repository,
	srv *sa.Service,
	log *logrus.Logger,
	v *viper.Viper,
) *Handler {
	return &Handler{
		log:                log,
		cognitoClient:      cognitoClient,
		rekognitiionClient: rekognitiionClient,
		textractClient:     textractClient,
		repo:               repo,
		ctx:                context.Background(),
		clientId:           v.GetString("client_id"),
		srv:                srv,
		v:                  v,
	}
}
func (h Handler) SignUp(c *fiber.Ctx) error {
	var req *model.UserSignUpRequest
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.Response{
			Message: err.Error(),
		})
	}
	err := h.srv.SignUp(req)
	if err != nil {
		return c.JSON(model.Response{
			Message: err.Error(),
		})
	}
	return c.JSON(model.Response{
		Message: "Email verification code is send to the email that u used in the sign up ",
	})
}

func (h Handler) SignIn(c *fiber.Ctx) error {
	var req *model.UserSignInReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.Response{
			Message: err.Error(),
		})
	}
	userData, respFromC, err := h.srv.SignIn(req)
	if err != nil {
		return c.JSON(model.Response{
			Message: err.Error(),
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    *respFromC.AuthenticationResult.RefreshToken,
		Secure:   true,
		HTTPOnly: true,
		MaxAge:   time.Now().Add(time.Hour * 24 * 365 * 10).Minute(),
	})

	//userData.AccessToken = model.AccessToken{
	//	AccessToken: *respFromC.AuthenticationResult.AccessToken,
	//}
	//userData.IdTOKEN = model.IdTOKEN{
	//	IdToken: *respFromC.AuthenticationResult.IdToken,
	//}
	userData.AccessToken = model.AccessToken{
		AccessToken: *respFromC.AuthenticationResult.AccessToken,
	}
	userData.IdTOKEN = model.IdTOKEN{
		IdToken: *respFromC.AuthenticationResult.IdToken,
	}
	return c.JSON(model.Response{
		Message: "OK",
		Data:    userData,
	})
}

func (h Handler) Confirm(c *fiber.Ctx) error {
	var req *model.EmailVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.Response{
			Message: err.Error(),
		})
	}
	err := h.srv.Confirm(req)
	if err != nil {
		return c.JSON(model.Response{
			Message: err.Error(),
		})
	}
	return c.JSON(model.Response{
		Message: "OK",
	})
}

func (h Handler) ResendConfirmation(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.JSON(model.Response{
			Message: "Please provide a valid email address",
		})
	}
	err := h.srv.ResendConfirmation(email)
	if err != nil {
		return c.JSON(model.Response{
			Message: err.Error(),
		})
	}
	return c.JSON(model.Response{
		Message: "Confirm code is resend",
	})
}

func (h Handler) ForgotPassword(c *fiber.Ctx) error {
	var req model.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}
	err := h.srv.ForgotPassword(req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password reset instructions sent to email",
	})
}

// ResetPasswordConfirm handler to confirm the new password
func (h Handler) ResetPasswordConfirm(c *fiber.Ctx) error {
	var req model.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}

	err := h.srv.ResetPasswordConfirm(req.Email, req.Code, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password successfully reset",
	})
}

// Logout handler to log the user out
func (h Handler) Logout(c *fiber.Ctx) error {
	accessToken := c.Get("Authorization")
	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing access token",
		})
	}

	err := h.srv.Logout(accessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (h Handler) KYCVerify(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email address is required in the param",
		})
	}
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files to verify",
		})
	}

	verification, err := h.srv.KYCVerification(files, email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.Response{
		Message: "KYC verification success",
		Data:    verification,
	})
}

func (h Handler) Me(c *fiber.Ctx) error {
	return c.JSON(model.Response{
		Message: "OK",
	})
}

func (h Handler) ForgotPasswordConfirm(ctx *fiber.Ctx) error {
	return ctx.JSON(model.Response{
		Message: "OK",
	})
}
