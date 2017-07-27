package restful

type Payload struct {
	Data     []Resource `json:"data"`
	Included []Resource `json:"included,omitempty"`
	Links    Links      `json:"links,omitempty"`
}

type Resource struct {
	Type          string                 `json:"type"`
	ID            string                 `json:"id,omitempty"`
	Attributes    map[string]interface{} `json:"attributes"`
	Relationships []Relationship         `json:"relationships"`
}

type Relationship struct {
	Key  string `json:"key"`
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
	Last  string `json:"last,omitempty"`
}
