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
