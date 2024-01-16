package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/constants"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/httpmodel"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/httpresponse"
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"github.com/mehmetokdemir/social-media-api/internal/app/guard"
	"github.com/mehmetokdemir/social-media-api/internal/middleware"
	"go.uber.org/zap"
	"net/http"
)

type HttpHandler struct {
	userService   IUserService
	logger        *zap.SugaredLogger
	jwtPrivateKey string
	guardService  guard.IGuardService
}

func NewHttpHandler(guardService guard.IGuardService, userService IUserService, logger *zap.SugaredLogger, jwtPrivateKey string) *HttpHandler {
	return &HttpHandler{guardService: guardService, userService: userService, logger: logger, jwtPrivateKey: jwtPrivateKey}
}

func (h *HttpHandler) RegisterRoutes(app *fiber.App) {
	authGroup := app.Group("/private").Use(middleware.AuthMiddleware(h.jwtPrivateKey))
	authGroup.Put("/update/photo", h.UpdatePhoto)

	noAuthGroup := app.Group("/public")
	noAuthGroup.Post("/register", h.Register)
}

// Register godoc
// @Summary Authenticate user
// @Description authenticates given user by giving an access jwttoken.
// @Tags User
// @Accept  json
// @Produce  json
// @Param request body RegisterRequest true "body params"
// @Success 200 {object} RegisterResponse
// @Failure 400
// @Failure 500
// @Router /public/register [post]
func (h *HttpHandler) Register(ctx *fiber.Ctx) error {
	//h.logger.Info("New register request arrived")
	var req RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse body", err.Error(), http.StatusBadRequest))
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
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not create user", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, RegisterResponse{Username: user.Username, Email: user.Email}))
}

// UpdatePhoto godoc
// @Summary Update own profile pic
// @Description Update profile pic with authentication
// @Tags User
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param image formData file true "The image file to upload"
// @Success 200 {object} httpmodel.UpdateImageResponse
// @Failure 400
// @Failure 500
// @Router /private/update/photo [put]
func (h *HttpHandler) UpdatePhoto(ctx *fiber.Ctx) error {
	//h.logger.Info("New register request arrived")
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}
	userID, exists := ctx.Locals(constants.UserIdKey).(uint)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context", "can not get user from context", http.StatusBadRequest))
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get file", err.Error(), http.StatusBadRequest))
	}

	uploadedImage, err := h.userService.UpdateProfilePhoto(userID, file)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not update profile photo", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, httpmodel.UpdateImageResponse{UploadedFileName: uploadedImage}))
}
