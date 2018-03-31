package restfulHelper

import (
	"encoding/json"
	"github.com/GoCollaborate/src/artifacts/restful"
	"github.com/GoCollaborate/src/constants"
	"github.com/GoCollaborate/src/utils"
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
