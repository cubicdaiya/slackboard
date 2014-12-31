package slackboard

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SlackboardResponse struct {
	Message string `json:"message"`
}

func sendResponse(w http.ResponseWriter, msg string, code int) {
	var (
		respBody       []byte
		respSlackboard SlackboardResponse
	)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", fmt.Sprintf("slackboard %s", Version))

	respSlackboard.Message = msg
	respBody, err := json.Marshal(respSlackboard)
	if err != nil {
		msg := "Response-body could not be created"
		http.Error(w, msg, http.StatusInternalServerError)
		LogError.Error(msg)
		return
	}

	switch code {
	case http.StatusOK:
		fmt.Fprintf(w, string(respBody))
	default:
		http.Error(w, string(respBody), code)
		LogError.Error(msg)
	}
}
