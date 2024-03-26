package controller

func generateResponseError(status int, message string) Response {
	return Response{
		Status: status,
		Message: message,
	}
}

func generateResponseSuccess(status int, message string, data any) Response {
	return Response{
		Status: status,
		Data:    data,
		Message: message,
	}
}