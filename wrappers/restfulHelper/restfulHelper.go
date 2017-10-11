package restfulHelper

import (
	"encoding/json"
	"github.com/GoCollaborate/artifacts/restful"
	"github.com/GoCollaborate/constants"
	"github.com/GoCollaborate/utils"
	"io"
	"net/http"
)

func SendErrorWith(w http.ResponseWriter, errPayload restful.ErrorPayload, header constants.Header) error {
	mal, err := json.Marshal(errPayload)
	if err != nil {
		return err
	}
	utils.AdaptHTTPWithHeader(w, header)
	io.WriteString(w, string(mal))
	return nil
}
