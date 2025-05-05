package sender

import (
	"bytes"
	"encoding/json"
	"gotel_alpha/data"
	"gotel_alpha/internal/handler"
	"gotel_alpha/util/network"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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
	const op = "sendMessage"
	params := make(map[string]string)
	params["chat_id"] = sendingMessage.ChatId
	params["text"] = sendingMessage.Value.(string)
	resp, err := network.GetRequest(*sender.token, op, params)
	if err != nil {
		return nil, err
	}
	respMessage, err := getMessage(resp)
	return respMessage, err
}

func (sender *Sender) SendPhoto(sendingPhoto *SendingEntity) (*data.Message, error) {
	const op = "sendPhoto"
	var resp *http.Response
	var err error
	params := make(map[string]string)
	params["chat_id"] = sendingPhoto.ChatId
	photo := sendingPhoto.Value.(*data.Photo)
	if photo.FileId == "" {
		body, writer, err := createForm(photo.FilePath)
		if err != nil {
			return nil, err
		}
		writer.Close()
		resp, err = network.PostRequest(*sender.token, op, params, body, writer)
	} else {
		params["photo"] = photo.FileId
		resp, err = network.GetRequest(*sender.token, op, params)
	}
	if err != nil {
		return nil, err
	}
	return getMessage(resp)
}

func createForm(filePath string) (*bytes.Buffer, *multipart.Writer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	part, _ := writer.CreateFormFile("photo", filePath)
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}
	return body, writer, nil
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
