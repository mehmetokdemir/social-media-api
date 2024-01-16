package post

type CreateRequest struct {
	Body string `json:"body"`
}

type UpdateRequest struct {
	Id    uint   `json:"id"`
	Body  string `json:"body"`
	Image string `json:"image"`
}

type ReadPostResponse struct {
	Id        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	Body      string `json:"body"`
	Image     string `json:"image"`
}
