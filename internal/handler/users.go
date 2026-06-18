package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Ainyx-backend/internal/models"
	"github.com/Ainyx-backend/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const userDateLayout = "2006-01-02"

type UserHandler struct {
	service  service.UserService
	validate *validator.Validate
	logger   *zap.Logger
}

type createUserRequest struct {
	Name string `json:"name" validate:"required"`
	DOB  string `json:"dob" validate:"required"`
}

type updateUserRequest struct {
	Name string `json:"name" validate:"required"`
	DOB  string `json:"dob" validate:"required"`
}

type listUsersQuery struct {
	Limit  int32 `query:"limit" validate:"gte=1,lte=100"`
	Offset int32 `query:"offset" validate:"gte=0"`
}

type userResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  *int   `json:"age,omitempty"`
}

func NewUserHandler(service service.UserService, validate *validator.Validate, logger *zap.Logger) *UserHandler {
	return &UserHandler{service: service, validate: validate, logger: logger}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req createUserRequest
	if err := c.BodyParser(&req); err != nil {
		return h.badRequest(c, "invalid request body")
	}
	if err := h.validate.Struct(req); err != nil {
		return h.validationError(c, err)
	}

	dob, err := time.Parse(userDateLayout, req.DOB)
	if err != nil {
		return h.badRequest(c, "dob must be in YYYY-MM-DD format")
	}

	user, err := h.service.CreateUser(c.UserContext(), models.CreateUserInput{
		Name: req.Name,
		DOB:  dob,
	})
	if err != nil {
		return h.handleServiceError(c, err)
	}

	h.logger.Info("user created", zap.Int32("id", user.ID), zap.String("request_id", requestID(c)))
	return c.Status(http.StatusCreated).JSON(toUserResponse(user))
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return h.badRequest(c, "invalid user id")
	}

	user, err := h.service.GetByID(c.UserContext(), id)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return c.JSON(toUserResponseWithAge(user))
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	query := listUsersQuery{Limit: 10}
	if err := c.QueryParser(&query); err != nil {
		return h.badRequest(c, "invalid query parameters")
	}
	if err := h.validate.Struct(query); err != nil {
		return h.validationError(c, err)
	}

	users, err := h.service.GetAll(c.UserContext(), query.Limit, query.Offset)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	response := make([]userResponse, 0, len(users))
	for _, user := range users {
		response = append(response, toUserResponseWithAge(user))
	}

	return c.JSON(response)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return h.badRequest(c, "invalid user id")
	}

	var req updateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return h.badRequest(c, "invalid request body")
	}
	if err := h.validate.Struct(req); err != nil {
		return h.validationError(c, err)
	}

	dob, err := time.Parse(userDateLayout, req.DOB)
	if err != nil {
		return h.badRequest(c, "dob must be in YYYY-MM-DD format")
	}

	user, err := h.service.Update(c.UserContext(), models.UpdateUserInput{
		ID:   id,
		Name: req.Name,
		DOB:  dob,
	})
	if err != nil {
		return h.handleServiceError(c, err)
	}

	h.logger.Info("user updated", zap.Int32("id", user.ID), zap.String("request_id", requestID(c)))
	return c.JSON(toUserResponse(user))
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return h.badRequest(c, "invalid user id")
	}


	if _, getUserErr := h.service.GetByID(c.UserContext(),id);getUserErr != nil{
		return h.badRequest(c,"User Does not exist")
	}
	if err := h.service.Delete(c.UserContext(), id); err != nil {
		return h.handleServiceError(c, err)
	}

	h.logger.Info("user deleted", zap.Int32("id", id), zap.String("request_id", requestID(c)))
	return c.SendStatus(http.StatusNoContent)
}

func (h *UserHandler) handleServiceError(c *fiber.Ctx, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	h.logger.Error("request failed", zap.Error(err), zap.String("request_id", requestID(c)))
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
}

func (h *UserHandler) badRequest(c *fiber.Ctx, message string) error {
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": message})
}

func (h *UserHandler) validationError(c *fiber.Ctx, err error) error {
	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
}

func parseID(raw string) (int32, error) {
	parsed, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}

	return int32(parsed), nil
}

func toUserResponse(user models.User) userResponse {
	return userResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.DOB.Format(userDateLayout),
	}
}

func toUserResponseWithAge(user models.User) userResponse {
	age := user.Age
	return userResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.DOB.Format(userDateLayout),
		Age:  &age,
	}
}

func requestID(c *fiber.Ctx) string {
	value, _ := c.Locals("requestId").(string)
	return value
}
