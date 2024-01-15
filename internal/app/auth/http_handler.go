package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/httperror"
	"github.com/mehmetokdemir/social-media-api/internal/middleware"
	"go.uber.org/zap"
)

type HttpHandler struct {
	authService   IAuthService
	logger        *zap.SugaredLogger
	jwtPrivateKey string
}

func NewHttpHandler(authService IAuthService, logger *zap.SugaredLogger, jwtPrivateKey string) *HttpHandler {
	return &HttpHandler{authService: authService, logger: logger, jwtPrivateKey: jwtPrivateKey}
}

func (h *HttpHandler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/auth")
	// Use(middleware.AuthMiddleware(h.jwtPrivateKey))
	appGroup.Post("/login", h.Login)
	appGroup.Post("/logout", h.Logout).Use(middleware.AuthMiddleware(h.jwtPrivateKey))
}

// Login godoc
// @Summary Login user
// @Description authenticates given user by giving an access jwttoken.
// @Param request body LoginRequest true "body params"
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 201 {object} LoginResponse
// @Failure 400
// @Failure 500
// @Router /auth/login [post]
func (h *HttpHandler) Login(ctx *fiber.Ctx) error {
	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httperror.NewBadRequestError(err.Error()))
	}

	rsp, err := h.authService.CreateToken(req.Username, req.Username)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperror.NewInternalServerError(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(rsp)
}

// Logout godoc
// @Summary Logout user
// @Description authenticates given user by giving an access jwttoken.
// @Param request body LoginRequest true "body params"
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 201 {object} LoginResponse
// @Failure 400
// @Failure 500
// @Router /auth/logout [post]
func (h *HttpHandler) Logout(ctx *fiber.Ctx) error {
	if err := h.authService.DeleteToken(ctx.Get("X-Auth-Token")); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httperror.NewBadRequestError(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON("")
}
