package restful

import (
	"encoding/json"
	. "github.com/GoCollaborate/src/artifacts/resources"
	"net/http"
)

type Request struct {
	Resources interface{}   `json:"data"`
	Included  interface{}   `json:"included"`
	reader    *http.Request `json:"-"`
}

type Response struct {
	Resources interface{}         `json:"data"`
	Links     *Links              `json:"links"`
	writer    http.ResponseWriter `json:"-"`
}

type RestfulResource interface {
	GetType() string
	GetId() string
	GetRelationships() map[string]*Relationship
}

func NewRequest() *Request {
	return new(Request)
}

func NewResponse() *Response {
	return new(Response)
}

func Writer(w http.ResponseWriter) *Response {

	var (
		res *Response = new(Response)
	)

	res.writer = w
	return res
}

func (res *Response) WithResource(resource RestfulResource) *Response {

	res.Resources = resource
	return res
}

func (res *Response) WithResources(resources ...RestfulResource) *Response {

	if len(resources) > 0 {
		res.Resources = resources
	}

	return res
}

func (res *Response) WithResourceArr(resources interface{}) *Response {
	res.Resources = resources
	return res
}

func (res *Response) WithLinks(links *Links) *Response {

	res.Links = links
	return res
}

func (res *Response) WithHeader(key, val string) *Response {

	res.writer.Header().Add(key, val)
	return res
}

func (res *Response) WithStatus(status int) *Response {

	res.writer.WriteHeader(status)
	return res
}

func (res *Response) Send() {

	mal, _ := json.Marshal(res)
	res.writer.Write(mal)
}

func Reader(r *http.Request) *Request {
	var (
		req *Request = new(Request)
	)

	req.reader = r
	return req
}

func (req *Request) WithResource(resource RestfulResource) *Request {

	req.Resources = resource

	return req
}

func (req *Request) WithResources(resources ...RestfulResource) *Request {

	if len(resources) > 0 {
		req.Resources = resources
	}

	return req
}

func (req *Request) WithResourceArr(resources interface{}) *Request {

	req.Resources = resources
	return req
}

func (req *Request) WithIncluded(resources ...RestfulResource) *Request {

	if len(resources) > 0 {
		req.Included = resources
	}

	return req
}

func (req *Request) Receive() *Request {

	json.NewDecoder(
		req.reader.Body,
	).Decode(req)
	return req
}
