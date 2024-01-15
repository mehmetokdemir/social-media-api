package httpresponse

type Response struct {
	Success    bool           `json:"success" extensions:"x-order=1" example:"true"`
	StatusCode int            `json:"status_code" extensions:"x-order=2" example:"200"`
	Error      *ResponseError `json:"error,omitempty" extensions:"x-order=4"`
	Data       interface{}    `json:"data,omitempty" extensions:"x-order=5"`
}

type ResponseError struct {
	Message string `json:"message" extensions:"x-order=1" example:"NOT_FOUND"`
	Detail  string `json:"detail" extensions:"x-order=2" example:"user not found"`
}
