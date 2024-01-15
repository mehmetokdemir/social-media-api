package comment

type CreateRequest struct {
	PostId uint   `json:"post_id"`
	Body   string `json:"body"`
	Image  string `json:"image"`
}
