package restful

import ()

type ErrorPayload struct {
	JSONAPI string  `json:"jsonapi, omitempty"`
	Errors  []Error `json:"errors"`
}

type Error struct {
	Code   string      `json:"code"`
	Source ErrorSource `json:"source"`
	Title  string      `json:"title"`
	Detail string      `json:"detail"`
}

type ErrorSource struct {
	Pointer string `json:"pointer"`
}

// create specific Error Types Below
func Error409Conflict() ErrorPayload {
	return ErrorPayload{"add_your_api_route_here", []Error{Error{"409", ErrorSource{"add_your_error_source_here"}, "add_your_error_title_here", "add_your_error_detail_here"}}}
}

func Error404NotFound() ErrorPayload {
	return ErrorPayload{"add_your_api_route_here", []Error{Error{"404", ErrorSource{"add_your_error_source_here"}, "add_your_error_title_here", "add_your_error_detail_here"}}}
}
