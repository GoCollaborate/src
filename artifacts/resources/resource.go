package resources

type Resource struct {
	Id            string                   `json:"id,omitempty"`
	Type          string                   `json:"type"`
	Relationships map[string]*Relationship `json:"relationships"`
}

func (r *Resource) GetType() string {
	return r.Type
}

func (r *Resource) GetId() string {
	return r.Id
}

func (r *Resource) GetRelationships() map[string]*Relationship {
	return r.Relationships
}

type Relationship struct {
	Links Links `json:"links"`
}

type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
	Last  string `json:"last,omitempty"`
}
