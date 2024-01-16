package comment

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/constants"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/httpmodel"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/httpresponse"
	"github.com/mehmetokdemir/social-media-api/internal/app/guard"
	"github.com/mehmetokdemir/social-media-api/internal/middleware"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type HttpHandler struct {
	commentService ICommentService
	logger         *zap.SugaredLogger
	jwtPrivateKey  string
	guardService   guard.IGuardService
}

func NewHttpHandler(guardService guard.IGuardService, commentService ICommentService, logger *zap.SugaredLogger, jwtPrivateKey string) *HttpHandler {
	return &HttpHandler{guardService: guardService, commentService: commentService, logger: logger, jwtPrivateKey: jwtPrivateKey}
}

func (h *HttpHandler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/comment").Use(middleware.AuthMiddleware(h.jwtPrivateKey))
	appGroup.Post("/create", h.Create)
	appGroup.Put("/update", h.Update)
	appGroup.Put("/update/:comment_id/image", h.UpdateImage)
	appGroup.Delete("/delete/:comment_id", h.Delete)
	appGroup.Get("/list/:post_id", h.List)

	appGroup.Get("/get/:comment_id", h.Get)
}

// Create godoc
// @Summary Create comment
// @Description Create comment from payload with POST method; need Authorization
// @Tags Comment
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param request body CreateRequest true "body params"
// @Success 200 {object} httpmodel.CreateResponse
// @Failure 400
// @Failure 500
// @Router /comment/create [post]
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

	comment, err := h.commentService.CreateComment(userID, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not create comment", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, httpmodel.CreateResponse{Id: comment.ID}))
}

// Update godoc
// @Summary Update Comment By ID
// @Description Update comment from payload with PUT method; need Authorization
// @Tags Comment
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param request body UpdateRequest true "body params"
// @Success 200 {object} httpresponse.Response "Success"
// @Failure 400
// @Failure 500
// @Router /comment/update [put]
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
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context ", "can not get user from context", http.StatusBadRequest))
	}

	if _, err := h.commentService.UpdateComment(userID, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not update comment", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, nil))
}

// UpdateImage godoc
// @Description Update comment image with authenticated user
// @Tags Comment
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param comment_id path string true "ID of the comment to update the image for"
// @Param image formData file true "The image file to upload"
// @Success 200 {object} httpmodel.UpdateImageResponse
// @Failure 400
// @Failure 500
// @Router /comment/update/{comment_id}/image [put]
func (h *HttpHandler) UpdateImage(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}
	userID, exists := ctx.Locals(constants.UserIdKey).(uint)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context", "can not get user from context", http.StatusBadRequest))
	}

	commentIDStr := ctx.Params("comment_id")
	if commentIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get comment_id on params", "can not get comment_id on params", http.StatusBadRequest))
	}

	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse comment_id", err.Error(), http.StatusBadRequest))
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get file", err.Error(), http.StatusBadRequest))
	}

	uploadedImage, err := h.commentService.UpdateCommentImage(uint(commentID), userID, file)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not update comment image", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, httpmodel.UpdateImageResponse{UploadedFileName: uploadedImage}))
}

// Delete godoc
// @Summary Delete Comment
// @Description Delete comment by comment id
// @Tags Comment
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param comment_id path integer true "ID of the comment to delete"
// @Success 200 {object} httpresponse.Response "Success"
// @Failure 400
// @Failure 500
// @Router /comment/delete/{comment_id} [delete]
func (h *HttpHandler) Delete(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}
	commentIDStr := ctx.Params("comment_id")
	if commentIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not find comment id on params", "can not find comment id on params", http.StatusBadRequest))
	}

	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse comment id ", err.Error(), http.StatusBadRequest))
	}

	userID, exists := ctx.Locals(constants.UserIdKey).(uint)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context ", "can not get user from context", http.StatusBadRequest))
	}

	if err = h.commentService.DeleteCommentById(userID, uint(commentID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not delete comment by id ", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, nil))
}

// List godoc
// @Summary List comments by post id
// @Description List comments with giving post_id which is requested from params. Authenticates given user by giving an access jwttoken.
// @Tags Comment
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param post_id path integer true "ID of the post to list comments"
// @Success 200 {object} []entity.Comment "Success"
// @Failure 400
// @Failure 500
// @Router /comment/list/{post_id} [get]
func (h *HttpHandler) List(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}
	postIDStr := ctx.Params("post_id")
	if postIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not find post id on params", "can not find post id on params", http.StatusBadRequest))
	}

	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse post id", err.Error(), http.StatusBadRequest))
	}

	comments, err := h.commentService.ListCommentsByPostID(uint(postID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not get comments", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, comments))
}

// Get godoc
// @Summary Get Comment BY ID
// @Description Get comment by id. This id must be taken from path and need authorization
// @Tags Comment
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param comment_id path integer true "ID of the comment to get comment"
// @Success 200 {object} entity.Comment "Success"
// @Failure 400
// @Failure 500
// @Router /comment/get/{comment_id} [get]
func (h *HttpHandler) Get(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}

	commentIDStr := ctx.Params("comment_id")
	if commentIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not find comment id on params", "", http.StatusBadRequest))
	}

	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not convert comment id from params", err.Error(), http.StatusBadRequest))
	}

	comment, err := h.commentService.GetCommentById(uint(commentID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not get comment from database", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, comment))

}
