package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/httpresponse"
	"github.com/mehmetokdemir/social-media-api/internal/app/guard"
	"github.com/mehmetokdemir/social-media-api/internal/middleware"
	"go.uber.org/zap"
	"net/http"
)

type HttpHandler struct {
	authService   IAuthService
	logger        *zap.SugaredLogger
	jwtPrivateKey string
	guardService  guard.IGuardService
}

func NewHttpHandler(guardService guard.IGuardService, authService IAuthService, logger *zap.SugaredLogger, jwtPrivateKey string) *HttpHandler {
	return &HttpHandler{guardService: guardService, authService: authService, logger: logger, jwtPrivateKey: jwtPrivateKey}
}

func (h *HttpHandler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/auth")
	// Use(middleware.AuthMiddleware(h.jwtPrivateKey))
	appGroup.Post("/login", h.Login)
	appGroup.Post("/logout", h.Logout).Use(middleware.AuthMiddleware(h.jwtPrivateKey))
}

// Login godoc
// @Summary Login user
// @Description Create token with given credentials
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body LoginRequest true "body params"
// @Success 200 {object} LoginResponse
// @Failure 400
// @Failure 500
// @Router /auth/login [post]
func (h *HttpHandler) Login(ctx *fiber.Ctx) error {
	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse body", err.Error(), http.StatusBadRequest))
	}

	rsp, err := h.authService.CreateToken(req.Username, req.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not create token", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, rsp))
}

// Logout godoc
// @Summary Logout user
// @Description Delete session with given access token.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Success 200 {object} httpresponse.Response "Success"
// @Failure 400
// @Failure 500
// @Router /auth/logout [post]
func (h *HttpHandler) Logout(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}

	if err := h.authService.DeleteToken(token); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not add tokens to black list", err.Error(), http.StatusBadRequest))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, nil))
}
