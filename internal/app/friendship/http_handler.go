package friendship

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mehmetokdemir/social-media-api/internal/middleware"
	"go.uber.org/zap"
)

type HttpHandler struct {
	FriendshipService IFriendshipService
	logger            *zap.SugaredLogger
	jwtPrivateKey     string
}

func NewHttpHandler(FriendshipService IFriendshipService, logger *zap.SugaredLogger, jwtPrivateKey string) *HttpHandler {
	return &HttpHandler{FriendshipService: FriendshipService, logger: logger, jwtPrivateKey: jwtPrivateKey}
}

func (h *HttpHandler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/Friendship").Use(middleware.AuthMiddleware(h.jwtPrivateKey))
	appGroup.Post("/add", h.Add)
	appGroup.Post("/delete", h.Delete)
}

// Add godoc
// @Summary Add user as friend
// @Description authenticates given user by giving an access jwttoken.
// @Tags Friendship
// @Accept  json
// @Produce  json
// @Failure 400
// @Failure 500
// @Router /Friendship/add [post]
func (h *HttpHandler) Add(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("")
}

// Delete godoc
// @Summary Delete user as friend
// @Description authenticates given user by giving an access jwttoken.
// @Tags Friendship
// @Accept  json
// @Produce  json
// @Failure 400
// @Failure 500
// @Router /Friendship/delete [post]
func (h *HttpHandler) Delete(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("")
}
