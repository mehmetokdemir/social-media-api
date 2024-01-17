package post

import "github.com/mehmetokdemir/social-media-api/internal/app/common/httpmodel"

type CreateRequest struct {
	Body string `json:"body"`
}

type UpdateRequest struct {
	Id    uint   `json:"id"`
	Body  string `json:"body"`
	Image string `json:"image"`
}

type ReadPostResponse struct {
	Id        uint                 `json:"id"`
	CreatedAt string               `json:"created_at"`
	User      httpmodel.CommonUser `json:"user"`
	Body      string               `json:"body"`
	Image     string               `json:"image"`
}
