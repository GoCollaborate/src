package restful

type Resource struct {
	ID            string         `json:"id,omitempty"`
	Type          string         `json:"type"`
	Relationships []Relationship `json:"relationships,omitempty"`
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
