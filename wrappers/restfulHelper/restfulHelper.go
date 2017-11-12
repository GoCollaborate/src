package restfulHelper

import (
	"encoding/json"
	"github.com/GoCollaborate/artifacts/restful"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/utils"
	"io"
	"net/http"
)

func SendErrorWith(w http.ResponseWriter, errPayload restful.ErrorPayload, status int) error {
	mal, err := json.Marshal(errPayload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, constants.HeaderContentTypeJSON)
	utils.AdaptHTTPWithStatus(w, status)
	io.WriteString(w, string(mal))
	return nil
}
