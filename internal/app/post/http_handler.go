package post

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
	"strconv"
)

type HttpHandler struct {
	postService   IPostService
	logger        *zap.SugaredLogger
	jwtPrivateKey string
	guardService  guard.IGuardService
}

func NewHttpHandler(guardService guard.IGuardService, postService IPostService, logger *zap.SugaredLogger, jwtPrivateKey string) *HttpHandler {
	return &HttpHandler{guardService: guardService, postService: postService, logger: logger, jwtPrivateKey: jwtPrivateKey}
}

func (h *HttpHandler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/post").Use(middleware.AuthMiddleware(h.jwtPrivateKey))
	appGroup.Post("/create", h.Create)
	appGroup.Put("/update", h.Update)
	appGroup.Put("/update/:post_id/image", h.UpdateImage)
	appGroup.Delete("/delete/:post_id", h.Delete)
	appGroup.Get("/list", h.List)
	appGroup.Get("/get/:post_id", h.Get)
}

// Create godoc
// @Summary Create post
// @Description Create post from body with authenticated user
// @Tags Post
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param request body CreateRequest true "body params"
// @Success 200 {object} httpmodel.CreateResponse
// @Failure 400
// @Failure 500
// @Router /post/create [post]
func (h *HttpHandler) Create(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}

	var req CreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse body", err.Error(), http.StatusBadRequest))
	}

	userID, exists := ctx.Locals(constants.UserIdKey).(uint)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context", "can not get user from context", http.StatusBadRequest))
	}

	// TODO ADD VALIDATOR

	post, err := h.postService.CreatePost(entity.Post{
		UserID: userID,
		Body:   req.Body,
	})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not create post", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, httpmodel.CreateResponse{Id: post.ID}))
}

// Update godoc
// @Description Update post from body with authenticated user
// @Tags Post
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param request body UpdateRequest true "body params"
// @Success 200 {object} httpresponse.Response "Success"
// @Failure 400
// @Failure 500
// @Router /post/update [put]
func (h *HttpHandler) Update(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}

	var req UpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse body", err.Error(), http.StatusBadRequest))
	}

	userID, exists := ctx.Locals(constants.UserIdKey).(uint)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context", "can not get user from context", http.StatusBadRequest))
	}

	if _, err := h.postService.UpdatePost(userID, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not update post", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, nil))
}

// UpdateImage godoc
// @Description Update post image with authenticated user
// @Tags Post
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param post_id path string true "ID of the post to update the image for"
// @Param image formData file true "The image file to upload"
// @Success 200 {object} httpmodel.UpdateImageResponse
// @Failure 400
// @Failure 500
// @Router /post/update/{post_id}/image [put]
func (h *HttpHandler) UpdateImage(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}

	userID, exists := ctx.Locals(constants.UserIdKey).(uint)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context", "can not get user from context", http.StatusBadRequest))
	}

	postIDStr := ctx.Params("post_id")
	if postIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get post_id on params", "can not get post_id on params", http.StatusBadRequest))
	}

	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse post_id", err.Error(), http.StatusBadRequest))
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get file", err.Error(), http.StatusBadRequest))
	}

	uploadedImage, err := h.postService.UpdatePostImage(uint(postID), userID, file)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not update post image", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, httpmodel.UpdateImageResponse{UploadedFileName: uploadedImage}))
}

// Delete godoc
// @Summary Delete post
// @Description authenticates given user by giving an access jwttoken.
// @Tags Post
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param post_id path integer true "ID of the post to delete"
// @Success 200 {object} httpresponse.Response "Success"
// @Failure 400
// @Failure 500
// @Router /post/delete/{post_id} [delete]
func (h *HttpHandler) Delete(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}

	userID, exists := ctx.Locals(constants.UserIdKey).(uint)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context", "can not get user from context", http.StatusBadRequest))
	}

	postIDStr := ctx.Params("post_id")
	if postIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get post_id on params", "can not get post_id on params", http.StatusBadRequest))
	}

	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse post_id", err.Error(), http.StatusBadRequest))
	}

	if err = h.postService.DeletePostById(userID, uint(postID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not delete post", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, nil))
}

// List godoc
// @Summary List post
// @Description authenticates given user by giving an access jwttoken.
// @Tags Post
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Success 200 {object} []ReadPostResponse "Success"
// @Failure 400
// @Failure 500
// @Router /post/list [get]
func (h *HttpHandler) List(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}

	posts, err := h.postService.ListPosts()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not get posts", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, posts))
}

// Get godoc
// @Summary Get post by id post
// @Description authenticates given user by giving an access jwttoken.
// @Tags Post
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Success 200 {object} ReadPostResponse "Success"
// @Failure 400
// @Failure 500
// @Router /post/get/{post_id} [get]
func (h *HttpHandler) Get(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}

	postIDStr := ctx.Params("post_id")
	if postIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get post_id on params", "can not get post_id on params", http.StatusBadRequest))
	}

	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse post_id", err.Error(), http.StatusBadRequest))
	}

	post, err := h.postService.GetPostById(uint(postID))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(httpresponse.NewError("can not get post", err.Error(), http.StatusNotFound))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, post))
}
