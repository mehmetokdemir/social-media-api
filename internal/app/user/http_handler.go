package user

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type HttpHandler struct {
	userService IUserService
	logger      *zap.SugaredLogger
}

func NewHttpHandler(userService IUserService, logger *zap.SugaredLogger) *HttpHandler {
	return &HttpHandler{userService: userService, logger: logger}
}

func (h *HttpHandler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/user")
	appGroup.Post("/login", h.Login)
	appGroup.Post("/register", h.Register)

	// Add,update profile pic
	// Like post
	// Like comment
	// Add friendship

	// Post create (in post repo)
}

func (h *HttpHandler) Register(ctx *fiber.Ctx) error {
	return nil
}

func (h *HttpHandler) Login(ctx *fiber.Ctx) error {
	return nil
}
