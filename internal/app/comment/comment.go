package comment

type CreateRequest struct {
	PostId   uint   `json:"post_id" extensions:"x-order=1" example:"1" validate:"required" valid:"required~post_id|invalid"`     // ID of the post
	Body     string `json:"body" extensions:"x-order=2" example:"New Post..." validate:"required" valid:"required~body|invalid"` // Body of the post
	ParentID *uint  `json:"parent_id" extensions:"x-order=3" example:"1" validate:"-"`                                           // ID of the parent comment id
}

type UpdateRequest struct {
	Id     uint   `json:"id" extensions:"x-order=1" example:"1" validate:"required" valid:"required~id|invalid"`           // ID of the comment
	PostId uint   `json:"post_id" extensions:"x-order=2" example:"3" validate:"required" valid:"required~post_id|invalid"` // ID of the post
	Body   string `json:"body" extensions:"x-order=3" example:"3" validate:"required" valid:"required~body|invalid"`       // Body of the post
}
