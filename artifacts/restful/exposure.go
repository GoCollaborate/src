package restful

type ExposurePayload struct {
	Data     []Resource `json:"data"`
	Included []Resource `json:"included,omitempty"`
	Links    Links      `json:"links,omitempty"`
}
