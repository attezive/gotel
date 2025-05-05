package deleter

import (
	"encoding/json"
	"gotel_alpha/util/network"
	"io"
	"net/http"
)

type DeleteResponse struct {
	Success     bool   `json:"ok"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

type Deleter struct {
	token *string
}

func NewDeleter(token *string) *Deleter {
	return &Deleter{token: token}
}

func (deleter Deleter) DeleteMessage(chatId string, messageId string) (*DeleteResponse, error) {
	const op = "deleteMessage"
	params := make(map[string]string, 2)
	params["chat_id"] = chatId
	params["message_id"] = messageId
	resp, err := network.GetRequest(*deleter.token, op, params)
	if err != nil {
		return nil, err
	}
	return getDeleteResponse(resp)
}

func getDeleteResponse(resp *http.Response) (*DeleteResponse, error) {
	var deleteResponse DeleteResponse
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &deleteResponse); err != nil {
		return nil, err
	}
	return &deleteResponse, nil
}
