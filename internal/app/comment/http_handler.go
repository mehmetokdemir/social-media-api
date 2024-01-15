package comment

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mehmetokdemir/social-media-api/internal/middleware"
	"go.uber.org/zap"
)

type HttpHandler struct {
	commentService ICommentService
	logger         *zap.SugaredLogger
	jwtPrivateKey  string
}

func NewHttpHandler(commentService ICommentService, logger *zap.SugaredLogger, jwtPrivateKey string) *HttpHandler {
	return &HttpHandler{commentService: commentService, logger: logger, jwtPrivateKey: jwtPrivateKey}
}

func (h *HttpHandler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/comment").Use(middleware.AuthMiddleware(h.jwtPrivateKey))
	appGroup.Post("/create", h.Create)
	appGroup.Post("/update", h.Update)
	appGroup.Post("/delete", h.Delete)
	appGroup.Post("/list", h.List)
	appGroup.Post("/get", h.Get)
}

// Create godoc
// @Summary Create comment
// @Description authenticates given user by giving an access jwttoken.
// @Tags Comment
// @Accept  json
// @Produce  json
// @Failure 400
// @Failure 500
// @Router /comment/create [post]
func (h *HttpHandler) Create(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("")
}

// Update godoc
// @Summary Create post
// @Description authenticates given user by giving an access jwttoken.
// @Tags Comment
// @Accept  json
// @Produce  json
// @Failure 400
// @Failure 500
// @Router /comment/update [post]
func (h *HttpHandler) Update(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("")
}

// Delete godoc
// @Summary Create post
// @Description authenticates given user by giving an access jwttoken.
// @Tags Comment
// @Accept  json
// @Produce  json
// @Failure 400
// @Failure 500
// @Router /comment/delete [post]
func (h *HttpHandler) Delete(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("")
}

// List godoc
// @Summary Create post
// @Description authenticates given user by giving an access jwttoken.
// @Tags Comment
// @Accept  json
// @Produce  json
// @Failure 400
// @Failure 500
// @Router /comment/list [post]
func (h *HttpHandler) List(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("")
}

// Get godoc
// @Summary Create post
// @Description authenticates given user by giving an access jwttoken.
// @Tags Comment
// @Accept  json
// @Produce  json
// @Failure 400
// @Failure 500
// @Router /comment/get [post]
func (h *HttpHandler) Get(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON("")
}
