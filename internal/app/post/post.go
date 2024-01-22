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
	Id         uint                      `json:"id"`
	CreatedAt  string                    `json:"created_at"`
	User       httpmodel.CommonUser      `json:"user"`
	Body       string                    `json:"body"`
	Image      string                    `json:"image"`
	LikedCount int64                     `json:"liked_count"`
	Comments   []ReadPostResponseComment `json:"comments,omitempty"`
}

type ReadPostResponseComment struct {
	Id          uint                      `json:"id"`
	Body        string                    `json:"body"`
	Image       string                    `json:"image"`
	User        httpmodel.CommonUser      `json:"user"`
	LikedCount  int64                     `json:"liked_count"`
	SubComments []ReadPostResponseComment `json:"sub_comments,omitempty"`
}
