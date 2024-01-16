package httpmodel

type UpdateImageResponse struct {
	UploadedFileName string `json:"uploaded_file_name" extensions:"x-order=1" example:"https://res-cdn.com/postId"` // Uploaded file name
}

type CreateResponse struct {
	Id uint `json:"id" extensions:"x-order=1" example:"1"` // ID of the created model`
}
