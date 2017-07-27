package restful

import (
	"encoding/json"
)

func Decode(js string) *Payload {
	payload := Payload{}
	json.Unmarshal([]byte(js), &payload)
	return &payload
}
