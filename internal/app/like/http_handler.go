package like

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/constants"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/httpresponse"
	"github.com/mehmetokdemir/social-media-api/internal/app/guard"
	"github.com/mehmetokdemir/social-media-api/internal/middleware"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type HttpHandler struct {
	likeService   ILikeService
	logger        *zap.SugaredLogger
	jwtPrivateKey string
	guardService  guard.IGuardService
}

func NewHttpHandler(guardService guard.IGuardService, likeService ILikeService, logger *zap.SugaredLogger, jwtPrivateKey string) *HttpHandler {
	return &HttpHandler{guardService: guardService, likeService: likeService, logger: logger, jwtPrivateKey: jwtPrivateKey}
}

func (h *HttpHandler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/like").Use(middleware.AuthMiddleware(h.jwtPrivateKey))
	appGroup.Post("/posts/:post_id", h.LikePost)
	appGroup.Post("/comments/:comment_id", h.LikeComment)
}

// LikePost godoc
// @Summary Like Post
// @Description Like post by post id. This id must be taken from path and need authorization
// @Tags Like
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param post_id path integer true "ID of the post to be liked"
// @Success 200 {object} httpresponse.Response "Success"
// @Failure 400
// @Failure 500
// @Router /like/posts/{post_id} [post]
func (h *HttpHandler) LikePost(ctx *fiber.Ctx) error {
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

	userID, exists := ctx.Locals(constants.UserIdKey).(uint)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context", "can not get user from context", http.StatusBadRequest))
	}

	if err = h.likeService.LikePost(userID, uint(postID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not like post", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, nil))
}

// LikeComment godoc
// @Summary Like Comment
// @Description Like comment by comment id. This id must be taken from path and need authorization
// @Tags Like
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param comment_id path integer true "ID of the comment to be liked"
// @Success 200 {object} httpresponse.Response "Success"
// @Failure 400
// @Failure 500
// @Router /like/comments/{comment_id} [post]
func (h *HttpHandler) LikeComment(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}

	commentIDStr := ctx.Params("comment_id")
	if commentIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get comment_id on params", "can not get comment_id on params", http.StatusBadRequest))
	}

	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse comment_id", err.Error(), http.StatusBadRequest))
	}

	userID, exists := ctx.Locals(constants.UserIdKey).(uint)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context", "can not get user from context", http.StatusBadRequest))
	}

	if err = h.likeService.LikeComment(userID, uint(commentID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not like comment", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, nil))
}
