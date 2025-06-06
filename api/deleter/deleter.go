package deleter

import (
	"encoding/json"
	"github.com/attezive/gotel/data"
	"github.com/attezive/gotel/util/network"
	"io"
	"net/http"
)

type Deleter struct {
	token *string
}

func CreateDeleter(token *string) *Deleter {
	return &Deleter{token: token}
}

func (deleter Deleter) DeleteMessage(chatId string, messageId string,
	rspCh chan<- *data.SuccessResponse, errCh chan<- error) {
	const op = "deleteMessage"
	params := make(map[string]string, 2)
	params["chat_id"] = chatId
	params["message_id"] = messageId
	resp, err := network.GetRequest(*deleter.token, op, params)
	if err != nil {
		errCh <- err
		return
	}
	success, err := getDeleteResponse(resp)
	rspCh <- success
	errCh <- err
}

func getDeleteResponse(resp *http.Response) (*data.SuccessResponse, error) {
	var deleteResponse data.SuccessResponse
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &deleteResponse); err != nil {
		return nil, err
	}
	return &deleteResponse, nil
}
