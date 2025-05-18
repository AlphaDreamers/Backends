package gig

import (
	"github.com/SwanHtetAungPhyo/backend/internal/model"
	"github.com/SwanHtetAungPhyo/backend/internal/service/gig"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type HandlerInterface interface {
	CreateGig(c *fiber.Ctx) error
	EditGig(c *fiber.Ctx) error
	UpdateGigByUserId(c *fiber.Ctx) error
	DeleteGig(c *fiber.Ctx) error
	GetSpecificGigByUserId(c *fiber.Ctx) error
	GetAllGigsByUserId(c *fiber.Ctx) error
}
type Handler struct {
	srv gig.Service
	log *logrus.Logger
}

func (h Handler) CreateGig(c *fiber.Ctx) error {
	var gigCreationReq *model.GigCreationReq
	if err := c.BodyParser(&gigCreationReq); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	return c.Status(201).SendString("New Gig is created successfully")
}

func (h Handler) EditGig(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (h Handler) UpdateGigByUserId(c *fiber.Ctx) error {
	userId := c.Params("userId")
	gigId := c.Params("gigId")
	if userId == "" || gigId == "" {
		return c.Status(400).SendString("UserId or GigId is empty")
	}
	return c.Status(200).SendString("Update Gig by user_id: " + userId)
}
func (h Handler) DeleteGig(c *fiber.Ctx) error {
	gigId := c.Params("gigId")
	userId := c.Params("userId")
	if gigId == "" || userId == "" {
		return c.Status(400).SendString("UserId or GigId is empty")
	}

	return c.Status(200).SendString("Delete Gig by user_id: " + userId)
}

func (h Handler) GetSpecificGigByUserId(c *fiber.Ctx) error {
	gigId := c.Params("gigId")
	userId := c.Params("userId")
	if userId == "" || gigId == "" {
		return c.Status(400).SendString("UserId or GigId is empty")
	}
	return c.Status(200).SendString("Get Specific Gig by user_id: " + userId)
}

func (h Handler) GetAllGigsByUserId(c *fiber.Ctx) error {
	userId := c.Params("userId")
	if userId == "" {
		return c.Status(400).SendString("UserId is empty")
	}

	return c.Status(200).SendString("Get All Gig by user_id: " + userId)
}

func NewHandler(
	srv gig.Service,
	log *logrus.Logger,
) *Handler {
	return &Handler{
		srv: srv,
		log: log,
	}
}

var _ HandlerInterface = (*Handler)(nil)
