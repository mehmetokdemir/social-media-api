package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/httperror"
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"go.uber.org/zap"
)

type HttpHandler struct {
	userService IUserService
	logger      *zap.SugaredLogger

	// TODO: REMOVE
	jwtPrivateKey string
}

func NewHttpHandler(userService IUserService, logger *zap.SugaredLogger, jwtPrivateKey string) *HttpHandler {
	return &HttpHandler{userService: userService, logger: logger, jwtPrivateKey: jwtPrivateKey}
}

func (h *HttpHandler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/user")
	appGroup.Post("/register", h.Register)
}

// Register godoc
// @Summary Authenticate user
// @Description authenticates given user by giving an access jwttoken.
// @Param request body RegisterRequest true "body params"
// @Tags User
// @Accept  json
// @Produce  json
// @Success 201 {object} RegisterResponse
// @Failure 400
// @Failure 500
// @Router /user/register [post]
func (h *HttpHandler) Register(ctx *fiber.Ctx) error {
	//h.logger.Info("New register request arrived")
	var req RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httperror.NewBadRequestError(err.Error()))
	}

	user, err := h.userService.CreateUser(entity.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Username:    req.Username,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httperror.NewBadRequestError(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(RegisterResponse{Username: user.Username, Email: user.Email})
}
