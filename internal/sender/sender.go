package sender

import (
	"encoding/json"
	"gotel_alpha/data"
	"gotel_alpha/internal/handler"
	"gotel_alpha/util"
	"io"
	"net/http"
	"strconv"
)

type Sender struct {
	token *string
}

type SendingEntity struct {
	ChatId string
	Value  interface{}
}

func NewSender(token *string) *Sender {
	return &Sender{token: token}
}

func (sender *Sender) SendMessage(sendingMessage *SendingEntity) (*data.Message, error) {
	params := make(map[string]string)
	params["chat_id"] = sendingMessage.ChatId
	params["text"] = sendingMessage.Value.(string)
	resp, err := util.GetRequest(*sender.token, "sendMessage", params)
	if err != nil {
		return nil, err
	}
	respMessage, err := getMessage(resp)
	return respMessage, err
}

func getMessage(resp *http.Response) (*data.Message, error) {
	var messageInfo data.MessageInfo
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &messageInfo); err != nil {
		return nil, err
	}
	return &messageInfo.Message, nil
}

func (sender *Sender) ReactionSend(
	update *handler.Update,
	value interface{},
	sendFunction func(*SendingEntity) (*data.Message, error),
	handleMessageFunction func(*data.Message)) error {

	entity := SendingEntity{
		ChatId: strconv.FormatInt(update.Message.Chat.Id, 10),
		Value:  value,
	}
	msg, err := sendFunction(&entity)
	if err != nil {
		return err
	}
	if handleMessageFunction != nil {
		handleMessageFunction(msg)
	}
	return nil
}
