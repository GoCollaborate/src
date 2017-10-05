package restful

type DashboardPayload struct {
	Base    Base           `json:"base"`
	Entries []EntriesGroup `json:"entries"`
	Models  []ModelsGroup  `json:"models"`
}

type EntriesGroup struct {
	Name    string  `json:"groupName"`
	Desc    string  `json:"groupDesc,omitempty"`
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Type       string `json:"type"`
	API        string `json:"api"`
	Desc       string `json:"desc,omitempty"`
	Deprecated bool   `json:"deprecated,omitempty"`
}

type ModelsGroup struct {
	Name   string  `json:"groupName"`
	Desc   string  `json:"groupDesc,omitempty"`
	Models []Model `json:"models"`
}

type Model struct {
	Name       string               `json:"name"`
	Attributes map[string]Attribute `json:"attributes"`
}

type Attribute struct {
	Type    string      `json:"type"`
	Default interface{} `json:"default,omitempty"`
}

type Base struct {
	Title    string `json:"title"`
	SubTitle string `json:"subTitle"`
}
