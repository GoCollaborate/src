package restful

type ErrorPayload struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Source string `json:"source"`
}

// create error types
func Error401Unauthorized() ErrorPayload {
	return ErrorPayload{[]Error{
		Error{401, "Unauthorized", "Although the HTTP standard specifies 'unauthorized', semantically this response means 'unauthenticated'.", ""}},
	}
}

func Error403Forbidden() ErrorPayload {
	return ErrorPayload{[]Error{
		Error{403, "Forbidden", "The client does not have access rights to the content.", ""}},
	}
}

func Error404NotFound() ErrorPayload {
	return ErrorPayload{[]Error{
		Error{404, "Not Found", "The server can not find requested resource.", ""}},
	}
}

func Error405MethodNotAllowed() ErrorPayload {
	return ErrorPayload{[]Error{
		Error{404, "Method Not Allowed", "The request method is known by the server but has been disabled and cannot be used.", ""}},
	}
}

func Error408RequestTimeout() ErrorPayload {
	return ErrorPayload{[]Error{
		Error{408, "Method Not Allowed", "This response is sent on an idle connection by some servers, even without any previous request by the client.", ""}},
	}
}

func Error409Conflict() ErrorPayload {
	return ErrorPayload{[]Error{
		Error{409, "Conflict", "This response is sent when a request conflicts with the current state of the server.", ""}},
	}
}

func Error415UnsupportedMediaType() ErrorPayload {
	return ErrorPayload{[]Error{
		Error{415, "Bad Request", "This response means that server could not understand the request due to invalid syntax.", ""}},
	}
}

func Error500InternalServerError() ErrorPayload {
	return ErrorPayload{[]Error{
		Error{500, "Internal Server Error", "The server has encountered a situation it doesn't know how to handle.", ""}},
	}
}

func Error502BadGateway() ErrorPayload {
	return ErrorPayload{[]Error{
		Error{502, "Bad Gateway", "This error response means that the server, while working as a gateway to get a response needed to handle the request, got an invalid response.", ""}},
	}
}

func Error503ServiceUnavailable() ErrorPayload {
	return ErrorPayload{[]Error{
		Error{503, "Service Unavailable", "The server is not ready to handle the request.", ""}},
	}
}
