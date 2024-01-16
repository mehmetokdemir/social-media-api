package friendship

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/constants"
	"github.com/mehmetokdemir/social-media-api/internal/app/common/httpresponse"
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
	"github.com/mehmetokdemir/social-media-api/internal/app/guard"
	"github.com/mehmetokdemir/social-media-api/internal/middleware"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type HttpHandler struct {
	friendshipService IFriendshipService
	logger            *zap.SugaredLogger
	jwtPrivateKey     string
	guardService      guard.IGuardService
}

func NewHttpHandler(guardService guard.IGuardService, friendshipService IFriendshipService, logger *zap.SugaredLogger, jwtPrivateKey string) *HttpHandler {
	return &HttpHandler{guardService: guardService, friendshipService: friendshipService, logger: logger, jwtPrivateKey: jwtPrivateKey}
}

func (h *HttpHandler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/friendship").Use(middleware.AuthMiddleware(h.jwtPrivateKey))
	appGroup.Post("/add/:user_id", h.Add)
	appGroup.Post("/remove/:user_id", h.Remove)
	appGroup.Post("/reject/:user_id", h.Reject)
	appGroup.Post("/accept/:user_id", h.Accept)
	appGroup.Get("/list/:status", h.List)
}

// Add godoc
// @Summary Add user as friend
// @Description Add user as a friend, this endpoint needs authentication
// @Tags Friendship
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param user_id path integer true "ID of the user to add as friendship"
// @Success 200 {object} httpresponse.Response "Success"
// @Failure 400
// @Failure 500
// @Router /friendship/add/{user_id} [post]
func (h *HttpHandler) Add(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}
	friendIDStr := ctx.Params("user_id")
	if friendIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user_id on params", "can not get user_id on params", http.StatusBadRequest))
	}

	friendID, err := strconv.ParseUint(friendIDStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse user_id", err.Error(), http.StatusBadRequest))
	}

	userID, exists := ctx.Locals(constants.UserIdKey).(uint)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context", "can not get user from context", http.StatusBadRequest))
	}

	if err = h.friendshipService.AddFriend(userID, uint(friendID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not request friendship", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, nil))
}

// Reject godoc
// @Summary Reject user from pending friendship request
// @Description Reject user from pending friendship request, this endpoint needs authentication
// @Tags Friendship
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param user_id path integer true "ID of the user to add as friendship"
// @Success 200 {object} httpresponse.Response "Success"
// @Failure 400
// @Failure 500
// @Router /friendship/reject/{request_id} [post]
func (h *HttpHandler) Reject(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}

	requestIDStr := ctx.Params("request_id")
	if requestIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get request_id on params", "can not get request_id on params", http.StatusBadRequest))
	}

	requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse request_id", err.Error(), http.StatusBadRequest))
	}

	userID, exists := ctx.Locals(constants.UserIdKey).(uint)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context", "can not get user from context", http.StatusBadRequest))
	}

	if err = h.friendshipService.RejectFriend(userID, uint(requestID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not reject user", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, nil))
}

// Accept godoc
// @Summary Accept user from pending friendship request
// @Description Accept user from pending friendship request, this endpoint needs authentication
// @Tags Friendship
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param user_id path integer true "ID of the user to add as friendship"
// @Success 200 {object} httpresponse.Response "Success"
// @Failure 400
// @Failure 500
// @Router /friendship/accept/{request_id} [post]
func (h *HttpHandler) Accept(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}
	requestIDStr := ctx.Params("request_id")
	if requestIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get request_id on params", "can not get request_id on params", http.StatusBadRequest))
	}

	requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse request_id", err.Error(), http.StatusBadRequest))
	}

	userID, exists := ctx.Locals(constants.UserIdKey).(uint)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context", "can not get user from context", http.StatusBadRequest))
	}

	if err = h.friendshipService.AcceptFriend(userID, uint(requestID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not accept friend request", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, nil))
}

// Remove godoc
// @Summary Remove user from friend
// @Description Remove user from friend, this endpoint needs authentication
// @Tags Friendship
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param request_id path integer true "ID of the request to be removed which is in list in friendship endpoint"
// @Success 200 {object} httpresponse.Response "Success"
// @Failure 400
// @Failure 500
// @Router /friendship/remove/{request_id} [post]
func (h *HttpHandler) Remove(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}
	requestIDStr := ctx.Params("request_id")
	if requestIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get request_id on params", "can not get request_id on params", http.StatusBadRequest))
	}

	requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not parse request_id", err.Error(), http.StatusBadRequest))
	}

	userID, exists := ctx.Locals(constants.UserIdKey).(uint)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context", "can not get user from context", http.StatusBadRequest))
	}

	if err = h.friendshipService.RemoveFriend(userID, uint(requestID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not remove friend", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, nil))
}

// List godoc
// @Summary List friendship
// @Description List friendship by status, if status empty it will return all data. Status is enum v
// @Description Status is enum values and get; pending and accepted values
// @Tags Friendship
// @Accept  json
// @Produce  json
// @Param X-Auth-Token header string true "Auth token of logged-in user."
// @Param status path string false "Filter with status it takes enum values which are; pending, accepted and also empty string. if status is empty string all documents will return"
// @Success 200 {object} []ReadFriendship "Success"
// @Failure 400
// @Failure 500
// @Router /friendship/list/{status} [get]
func (h *HttpHandler) List(ctx *fiber.Ctx) error {
	token := ctx.Get("X-Auth-Token")
	if ok := h.guardService.CheckTokenInBlacklist(token); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(httpresponse.NewError("invalid token", "token is not valid", http.StatusForbidden))
	}

	var filterStatus *entity.FriendshipStatusEnum
	status := ctx.Params("status")
	if status == "" {
		filterStatus = nil
	} else {
		switch status {
		case string(entity.FriendshipStatusPending):
			p := entity.FriendshipStatusPending
			filterStatus = &p
		case string(entity.FriendshipStatusAccepted):
			p := entity.FriendshipStatusAccepted
			filterStatus = &p
		default:
			filterStatus = nil
		}
	}

	userID, exists := ctx.Locals(constants.UserIdKey).(uint)
	if !exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(httpresponse.NewError("can not get user from context", "can not get user from context", http.StatusBadRequest))
	}

	friendships, err := h.friendshipService.ListFriends(userID, filterStatus)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(httpresponse.NewError("can not get friend list", err.Error(), http.StatusInternalServerError))
	}

	return ctx.Status(fiber.StatusOK).JSON(httpresponse.NewSuccess(ctx, http.StatusOK, friendships))
}
