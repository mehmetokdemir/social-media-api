package httpresponse

func NewError(message, detail string, statusCode int) Response {
	return Response{
		Success:    false,
		StatusCode: statusCode,
		Error: &ResponseError{
			Message: message,
			Detail:  detail,
		},
		Data: nil,
	}
}
